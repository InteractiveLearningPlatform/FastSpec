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

    use super::{parse_spec_file, validate_spec_tree};

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

    fn unique_temp_file(name: &str) -> PathBuf {
        let unique = SystemTime::now().duration_since(UNIX_EPOCH).expect("time should move forward").as_nanos();
        std::env::temp_dir().join(format!("fastspec-{unique}-{name}"))
    }
}
