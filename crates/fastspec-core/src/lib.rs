use std::collections::{HashMap, HashSet};
use std::fs;
use std::io;
use std::path::{Path, PathBuf};

use fastspec_model::{FastSpecDocument, SpecKind, parse_document};
use serde::Serialize;

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct SpecSummary {
    pub path: PathBuf,
    pub kind: SpecKind,
    pub id: String,
    pub title: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct SummaryOutput {
    pub documents: Vec<SpecSummary>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct InspectDocument {
    pub path: PathBuf,
    pub metadata: fastspec_model::Metadata,
    #[serde(flatten)]
    pub document: FastSpecDocument,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct InspectOutput {
    pub documents: Vec<InspectDocument>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "snake_case")]
pub enum ValidationSeverity {
    Error,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct ValidationFinding {
    pub code: String,
    pub severity: ValidationSeverity,
    pub message: String,
    pub path: PathBuf,
    pub document_id: Option<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct ValidationOutput {
    pub valid: bool,
    pub findings: Vec<ValidationFinding>,
}

pub fn collect_spec_files(root: &Path) -> io::Result<Vec<PathBuf>> {
    let mut files = Vec::new();
    visit(root, &mut files, true)?;
    files.sort();
    Ok(files)
}

pub fn summarize_specs(root: &Path) -> io::Result<Vec<SpecSummary>> {
    parse_spec_tree(root).map(|documents| documents.into_iter().map(SpecDocument::into_summary).collect())
}

pub fn validate_spec_tree(root: &Path) -> io::Result<Vec<SpecSummary>> {
    let documents = parse_spec_tree(root)?;
    if documents.is_empty() {
        return Err(io::Error::new(io::ErrorKind::NotFound, format!("no .yaml files found under {}", root.display())));
    }

    Ok(documents.into_iter().map(SpecDocument::into_summary).collect())
}

pub fn validate_findings(path: &Path) -> io::Result<ValidationOutput> {
    let documents = parse_spec_path(path)?;
    if documents.is_empty() {
        return Err(io::Error::new(io::ErrorKind::NotFound, format!("no .yaml files found under {}", path.display())));
    }

    let mut findings = Vec::new();

    let mut id_to_documents: HashMap<String, Vec<&SpecDocument>> = HashMap::new();
    for document in &documents {
        id_to_documents.entry(document.document.metadata().id.clone()).or_default().push(document);
    }

    for (id, matches) in id_to_documents {
        if matches.len() > 1 {
            findings.push(ValidationFinding {
                code: "duplicate_identifier".to_string(),
                severity: ValidationSeverity::Error,
                message: format!("document identifier `{id}` is used by multiple documents"),
                path: matches[0].path.clone(),
                document_id: Some(id),
            });
        }
    }

    let module_documents: Vec<&SpecDocument> =
        documents.iter().filter(|document| matches!(document.document, FastSpecDocument::Module(_))).collect();
    let actual_module_ids: HashSet<String> = module_documents.iter().map(|document| document.document.metadata().id.clone()).collect();

    let project_documents: Vec<&SpecDocument> =
        documents.iter().filter(|document| matches!(document.document, FastSpecDocument::Project(_))).collect();

    for project_document in project_documents {
        let FastSpecDocument::Project(project) = &project_document.document else {
            continue;
        };

        let declared_module_ids: HashSet<String> = project.spec.modules.iter().map(|module| module.id.clone()).collect();

        for module_id in &declared_module_ids {
            if !actual_module_ids.contains(module_id) {
                findings.push(ValidationFinding {
                    code: "missing_module_document".to_string(),
                    severity: ValidationSeverity::Error,
                    message: format!("project declares module `{module_id}` but no matching module document exists"),
                    path: project_document.path.clone(),
                    document_id: Some(project.metadata.id.clone()),
                });
            }
        }

        for module_document in &module_documents {
            let FastSpecDocument::Module(module) = &module_document.document else {
                continue;
            };

            if !declared_module_ids.contains(&module.metadata.id) {
                findings.push(ValidationFinding {
                    code: "undeclared_module_document".to_string(),
                    severity: ValidationSeverity::Error,
                    message: format!(
                        "module document `{}` exists but is not declared in project `{}`",
                        module.metadata.id, project.metadata.id
                    ),
                    path: module_document.path.clone(),
                    document_id: Some(module.metadata.id.clone()),
                });
            }

            for dependency in &module.spec.dependencies {
                let references_known_module_doc = actual_module_ids.contains(&dependency.id);
                let references_declared_module = declared_module_ids.contains(&dependency.id);

                if references_declared_module && !references_known_module_doc {
                    findings.push(ValidationFinding {
                        code: "invalid_module_dependency".to_string(),
                        severity: ValidationSeverity::Error,
                        message: format!(
                            "module `{}` depends on declared project module `{}` but no matching module document exists",
                            module.metadata.id, dependency.id
                        ),
                        path: module_document.path.clone(),
                        document_id: Some(module.metadata.id.clone()),
                    });
                } else if references_known_module_doc && !references_declared_module {
                    findings.push(ValidationFinding {
                        code: "invalid_module_dependency".to_string(),
                        severity: ValidationSeverity::Error,
                        message: format!(
                            "module `{}` depends on module document `{}` that is not declared in project `{}`",
                            module.metadata.id, dependency.id, project.metadata.id
                        ),
                        path: module_document.path.clone(),
                        document_id: Some(module.metadata.id.clone()),
                    });
                }
            }
        }
    }

    findings.sort_by(|left, right| left.code.cmp(&right.code).then(left.path.cmp(&right.path)));

    Ok(ValidationOutput { valid: findings.is_empty(), findings })
}

#[derive(Debug, Clone)]
pub struct SpecDocument {
    pub path: PathBuf,
    pub document: FastSpecDocument,
}

impl SpecDocument {
    pub fn into_summary(self) -> SpecSummary {
        SpecSummary {
            path: self.path,
            kind: self.document.kind(),
            id: self.document.metadata().id.clone(),
            title: self.document.metadata().title.clone(),
        }
    }

    pub fn into_inspect(self) -> InspectDocument {
        InspectDocument { path: self.path, metadata: self.document.metadata().clone(), document: self.document }
    }
}

pub fn parse_spec_path(path: &Path) -> io::Result<Vec<SpecDocument>> {
    if path.is_dir() { parse_spec_tree(path) } else { parse_spec_file(path).map(|document| vec![document]) }
}

pub fn parse_spec_tree(root: &Path) -> io::Result<Vec<SpecDocument>> {
    collect_spec_files(root)?.into_iter().map(|path| parse_spec_file(&path)).collect()
}

pub fn parse_spec_file(path: &Path) -> io::Result<SpecDocument> {
    let contents = fs::read_to_string(path)?;
    let document = parse_document(&contents)
        .map_err(|error| io::Error::new(io::ErrorKind::InvalidData, format!("invalid FastSpec document in {}: {error}", path.display())))?;

    Ok(SpecDocument { path: path.to_path_buf(), document })
}

fn visit(root: &Path, files: &mut Vec<PathBuf>, require_yaml: bool) -> io::Result<()> {
    if root.is_file() {
        if !require_yaml || matches!(root.extension().and_then(|ext| ext.to_str()), Some("yaml" | "yml")) {
            files.push(root.to_path_buf());
        }
        return Ok(());
    }

    for entry in fs::read_dir(root)? {
        let entry = entry?;
        let path = entry.path();
        if path.is_dir() {
            visit(&path, files, require_yaml)?;
        } else if matches!(path.extension().and_then(|ext| ext.to_str()), Some("yaml" | "yml")) {
            files.push(path);
        }
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use std::fs;
    use std::path::PathBuf;
    use std::time::{SystemTime, UNIX_EPOCH};

    use fastspec_model::SpecKind;

    use super::{parse_spec_file, validate_findings, validate_spec_tree};

    #[test]
    fn validates_archlint_example_tree() {
        let root = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../examples/archlint-reproduction/specs");
        let summaries = validate_spec_tree(&root).expect("example specs should validate");

        assert_eq!(summaries.len(), 3);
        assert!(summaries.iter().any(|summary| summary.kind == SpecKind::Project));
        assert_eq!(summaries.iter().filter(|summary| summary.kind == SpecKind::Module).count(), 2);
        assert!(summaries.iter().any(|summary| summary.id == "archlint-reproduction"));
    }

    #[test]
    fn rejects_invalid_document() {
        let path = unique_temp_file("invalid.fastspec.yaml");
        fs::write(
            &path,
            "apiVersion: fastspec.dev/v0alpha1\nkind: ModuleSpec\nmetadata:\n  id: broken\n  title: Broken module\n  summary: Missing purpose field\nspec:\n  inputs: []\n",
        )
        .expect("temp file should write");

        let error = parse_spec_file(&path).expect_err("invalid spec should fail");
        assert!(error.to_string().contains("invalid FastSpec document"));
        assert!(error.to_string().contains("purpose"));

        fs::remove_file(path).expect("temp file should be cleaned");
    }

    #[test]
    fn reports_clean_validation_for_example_tree() {
        let root = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../examples/archlint-reproduction/specs");
        let output = validate_findings(&root).expect("example tree should validate");
        assert!(output.valid);
        assert!(output.findings.is_empty());
    }

    #[test]
    fn reports_cross_document_findings() {
        let root = unique_temp_dir("validation-fixtures");
        fs::create_dir_all(root.join("modules")).expect("fixture directories should be created");

        fs::write(
            root.join("project.fastspec.yaml"),
            "apiVersion: fastspec.dev/v0alpha1\nkind: ProjectSpec\nmetadata:\n  id: demo\n  title: Demo\n  summary: Demo project\nspec:\n  modules:\n    - id: api\n      purpose: API module\n    - id: web\n      purpose: Web module\n",
        )
        .expect("project fixture should write");
        fs::write(
            root.join("modules/api.fastspec.yaml"),
            "apiVersion: fastspec.dev/v0alpha1\nkind: ModuleSpec\nmetadata:\n  id: api\n  title: API\n  summary: API module\nspec:\n  purpose: Serve data\n  dependencies:\n    - id: ghost\n      reason: Internal ghost dependency\n",
        )
        .expect("api fixture should write");
        fs::write(
            root.join("modules/ghost.fastspec.yaml"),
            "apiVersion: fastspec.dev/v0alpha1\nkind: ModuleSpec\nmetadata:\n  id: ghost\n  title: Ghost\n  summary: Ghost module\nspec:\n  purpose: Hidden module\n",
        )
        .expect("ghost fixture should write");
        fs::write(
            root.join("modules/api-duplicate.fastspec.yaml"),
            "apiVersion: fastspec.dev/v0alpha1\nkind: ModuleSpec\nmetadata:\n  id: api\n  title: API Duplicate\n  summary: Duplicate API module\nspec:\n  purpose: Duplicate module\n",
        )
        .expect("duplicate fixture should write");

        let output = validate_findings(&root).expect("validation should succeed with findings");
        assert!(!output.valid);
        assert!(output.findings.iter().any(|finding| finding.code == "duplicate_identifier"));
        assert!(output.findings.iter().any(|finding| finding.code == "missing_module_document"));
        assert!(output.findings.iter().any(|finding| finding.code == "undeclared_module_document"));
        assert!(output.findings.iter().any(|finding| finding.code == "invalid_module_dependency"));

        fs::remove_dir_all(root).expect("fixture dir should be removed");
    }

    fn unique_temp_file(name: &str) -> PathBuf {
        let unique = SystemTime::now().duration_since(UNIX_EPOCH).expect("time should move forward").as_nanos();
        std::env::temp_dir().join(format!("fastspec-{unique}-{name}"))
    }

    fn unique_temp_dir(name: &str) -> PathBuf {
        let unique = SystemTime::now().duration_since(UNIX_EPOCH).expect("time should move forward").as_nanos();
        std::env::temp_dir().join(format!("fastspec-{unique}-{name}"))
    }
}
