use std::collections::{HashMap, HashSet};
use std::fs;
use std::io;
use std::path::{Path, PathBuf};

use fastspec_model::{FastSpecDocument, ProjectSpecDocument, SpecKind, parse_document};
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

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "snake_case")]
pub enum GraphNodeKind {
    Project,
    Module,
    Workflow,
    AgentCapability,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "snake_case")]
pub enum GraphEdgeKind {
    Contains,
    DefinesWorkflow,
    DefinesCapability,
    DependsOn,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct GraphNode {
    pub id: String,
    pub kind: GraphNodeKind,
    pub title: String,
    pub path: PathBuf,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct GraphEdge {
    pub from: String,
    pub to: String,
    pub kind: GraphEdgeKind,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct GraphOutput {
    pub nodes: Vec<GraphNode>,
    pub edges: Vec<GraphEdge>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "snake_case")]
pub enum PlanPhase {
    Bootstrap,
    Module,
    AgentCapability,
    Workflow,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct PlanStep {
    pub id: String,
    pub phase: PlanPhase,
    pub title: String,
    pub depends_on: Vec<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct PlanOutput {
    pub steps: Vec<PlanStep>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "snake_case")]
pub enum GeneratedArtifactKind {
    Directory,
    File,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct GeneratedArtifact {
    pub path: PathBuf,
    pub kind: GeneratedArtifactKind,
    pub description: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct ScaffoldOutput {
    pub output_dir: PathBuf,
    pub artifacts: Vec<GeneratedArtifact>,
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

    for project_document in &project_documents {
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

    let agent_capability_documents: Vec<&SpecDocument> =
        documents.iter().filter(|document| matches!(document.document, FastSpecDocument::AgentCapability(_))).collect();
    let actual_capability_ids: HashSet<String> =
        agent_capability_documents.iter().map(|document| document.document.metadata().id.clone()).collect();

    for project_document in project_documents {
        let FastSpecDocument::Project(project) = &project_document.document else {
            continue;
        };

        let declared_capability_ids: HashSet<String> =
            project.spec.agent_capabilities.iter().map(|capability| capability.id.clone()).collect();

        for capability_id in &declared_capability_ids {
            if !actual_capability_ids.contains(capability_id) {
                findings.push(ValidationFinding {
                    code: "missing_agent_capability_document".to_string(),
                    severity: ValidationSeverity::Error,
                    message: format!(
                        "project declares agent capability `{capability_id}` but no matching agent capability document exists"
                    ),
                    path: project_document.path.clone(),
                    document_id: Some(project.metadata.id.clone()),
                });
            }
        }

        for capability_document in &agent_capability_documents {
            if !declared_capability_ids.contains(&capability_document.document.metadata().id) {
                findings.push(ValidationFinding {
                    code: "undeclared_agent_capability_document".to_string(),
                    severity: ValidationSeverity::Error,
                    message: format!(
                        "agent capability document `{}` exists but is not declared in project `{}`",
                        capability_document.document.metadata().id,
                        project.metadata.id
                    ),
                    path: capability_document.path.clone(),
                    document_id: Some(capability_document.document.metadata().id.clone()),
                });
            }
        }
    }

    let workflow_documents: Vec<&SpecDocument> =
        documents.iter().filter(|document| matches!(document.document, FastSpecDocument::Workflow(_))).collect();
    let actual_workflow_ids: HashSet<String> =
        workflow_documents.iter().map(|document| document.document.metadata().id.clone()).collect();

    for project_document in &project_documents {
        let FastSpecDocument::Project(project) = &project_document.document else {
            continue;
        };

        let declared_workflow_ids: HashSet<String> =
            project.spec.workflows.iter().map(|workflow| workflow.id.clone()).collect();

        for workflow_id in &declared_workflow_ids {
            if !actual_workflow_ids.contains(workflow_id) {
                findings.push(ValidationFinding {
                    code: "missing_workflow_document".to_string(),
                    severity: ValidationSeverity::Error,
                    message: format!(
                        "project declares workflow `{workflow_id}` but no matching workflow document exists"
                    ),
                    path: project_document.path.clone(),
                    document_id: Some(project.metadata.id.clone()),
                });
            }
        }

        for workflow_document in &workflow_documents {
            if !declared_workflow_ids.contains(&workflow_document.document.metadata().id) {
                findings.push(ValidationFinding {
                    code: "undeclared_workflow_document".to_string(),
                    severity: ValidationSeverity::Error,
                    message: format!(
                        "workflow document `{}` exists but is not declared in project `{}`",
                        workflow_document.document.metadata().id,
                        project.metadata.id
                    ),
                    path: workflow_document.path.clone(),
                    document_id: Some(workflow_document.document.metadata().id.clone()),
                });
            }
        }
    }

    // Build module dependency graph and detect cycles via DFS.
    // A cycle is a structural error that prevents meaningful planning.
    let module_dep_graph: HashMap<&str, Vec<&str>> = module_documents
        .iter()
        .filter_map(|document| {
            if let FastSpecDocument::Module(module) = &document.document {
                let deps: Vec<&str> = module.spec.dependencies.iter().map(|d| d.id.as_str()).collect();
                Some((module.metadata.id.as_str(), deps))
            } else {
                None
            }
        })
        .collect();

    // Two-color DFS: `color` tracks completed nodes (done), `in_stack` tracks the active path (gray).
    let all_module_ids: Vec<&str> = module_dep_graph.keys().copied().collect();
    let mut color: HashMap<&str, bool> = HashMap::new(); // true = done, absent = unvisited
    let mut in_stack: HashSet<&str> = HashSet::new();
    let mut reported_cycles: HashSet<String> = HashSet::new();

    for start in &all_module_ids {
        if color.contains_key(start) {
            continue;
        }
        let mut stack: Vec<(&str, usize)> = vec![(start, 0)];
        in_stack.insert(start);
        while let Some((node, idx)) = stack.last().copied() {
            let deps = module_dep_graph.get(node).map(|v| v.as_slice()).unwrap_or(&[]);
            if idx >= deps.len() {
                stack.pop();
                in_stack.remove(node);
                color.insert(node, true);
            } else {
                let last = stack.last_mut().unwrap();
                last.1 += 1;
                let next = deps[idx];
                if !module_dep_graph.contains_key(next) {
                    // External or unknown dependency — already caught by other rules.
                    continue;
                }
                if in_stack.contains(next) {
                    // Back-edge detected: cycle from `next` back to itself via `node`.
                    let cycle_key = {
                        let mut pair = [node, next];
                        pair.sort_unstable();
                        format!("{}-{}", pair[0], pair[1])
                    };
                    if reported_cycles.insert(cycle_key) {
                        let path =
                            module_documents.iter().find(|d| d.document.metadata().id == node).map(|d| d.path.clone()).unwrap_or_default();
                        findings.push(ValidationFinding {
                            code: "module_dependency_cycle".to_string(),
                            severity: ValidationSeverity::Error,
                            message: format!("module `{node}` participates in a dependency cycle via `{next}`"),
                            path,
                            document_id: Some(node.to_string()),
                        });
                    }
                } else if !color.contains_key(next) {
                    stack.push((next, 0));
                    in_stack.insert(next);
                }
            }
        }
    }

    findings.sort_by(|left, right| left.code.cmp(&right.code).then(left.path.cmp(&right.path)));

    Ok(ValidationOutput { valid: findings.is_empty(), findings })
}

pub fn export_graph(path: &Path) -> io::Result<GraphOutput> {
    let validation = validate_findings(path)?;
    if !validation.valid {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            format!("graph export requires a validation-clean tree; found {} validation finding(s)", validation.findings.len()),
        ));
    }

    let documents = parse_spec_path(path)?;
    let mut nodes = Vec::new();
    let mut edges = Vec::new();
    let mut module_paths: HashMap<String, PathBuf> = HashMap::new();
    let mut module_titles: HashMap<String, String> = HashMap::new();
    let mut capability_paths: HashMap<String, PathBuf> = HashMap::new();
    let mut capability_titles: HashMap<String, String> = HashMap::new();
    let mut workflow_paths: HashMap<String, PathBuf> = HashMap::new();
    let mut workflow_titles: HashMap<String, String> = HashMap::new();

    for document in &documents {
        match &document.document {
            FastSpecDocument::Module(module) => {
                module_paths.insert(module.metadata.id.clone(), document.path.clone());
                module_titles.insert(module.metadata.id.clone(), module.metadata.title.clone());
            }
            FastSpecDocument::AgentCapability(capability) => {
                capability_paths.insert(capability.metadata.id.clone(), document.path.clone());
                capability_titles.insert(capability.metadata.id.clone(), capability.metadata.title.clone());
            }
            FastSpecDocument::Workflow(workflow) => {
                workflow_paths.insert(workflow.metadata.id.clone(), document.path.clone());
                workflow_titles.insert(workflow.metadata.id.clone(), workflow.metadata.title.clone());
            }
            FastSpecDocument::Project(_) => {}
        }
    }
    }

    for document in &documents {
        match &document.document {
            FastSpecDocument::Project(project) => {
                let project_id = project.metadata.id.clone();
                nodes.push(GraphNode {
                    id: project_id.clone(),
                    kind: GraphNodeKind::Project,
                    title: project.metadata.title.clone(),
                    path: document.path.clone(),
                });

                for declared_module in &project.spec.modules {
                    if let (Some(path), Some(title)) = (module_paths.get(&declared_module.id), module_titles.get(&declared_module.id)) {
                        nodes.push(GraphNode {
                            id: declared_module.id.clone(),
                            kind: GraphNodeKind::Module,
                            title: title.clone(),
                            path: path.clone(),
                        });
                        edges.push(GraphEdge { from: project_id.clone(), to: declared_module.id.clone(), kind: GraphEdgeKind::Contains });
                    }
                }

                for workflow in &project.spec.workflows {
                    if let (Some(path), Some(title)) =
                        (workflow_paths.get(&workflow.id), workflow_titles.get(&workflow.id))
                    {
                        let workflow_node_id = format!("workflow:{}", workflow.id);
                        nodes.push(GraphNode {
                            id: workflow_node_id.clone(),
                            kind: GraphNodeKind::Workflow,
                            title: title.clone(),
                            path: path.clone(),
                        });
                        edges.push(GraphEdge {
                            from: project_id.clone(),
                            to: workflow_node_id,
                            kind: GraphEdgeKind::DefinesWorkflow,
                        });
                    }
                }

                for declared_capability in &project.spec.agent_capabilities {
                    if let (Some(path), Some(title)) =
                        (capability_paths.get(&declared_capability.id), capability_titles.get(&declared_capability.id))
                    {
                        nodes.push(GraphNode {
                            id: declared_capability.id.clone(),
                            kind: GraphNodeKind::AgentCapability,
                            title: title.clone(),
                            path: path.clone(),
                        });
                        edges.push(GraphEdge {
                            from: project_id.clone(),
                            to: declared_capability.id.clone(),
                            kind: GraphEdgeKind::DefinesCapability,
                        });
                    }
                }
            }
            FastSpecDocument::Module(module) => {
                for dependency in &module.spec.dependencies {
                    if module_paths.contains_key(&dependency.id) {
                        edges.push(GraphEdge {
                            from: module.metadata.id.clone(),
                            to: dependency.id.clone(),
                            kind: GraphEdgeKind::DependsOn,
                        });
                    }
                }
            }
            FastSpecDocument::AgentCapability(_) | FastSpecDocument::Workflow(_) => {}
        }
    }

    nodes.sort_by(|left, right| left.id.cmp(&right.id).then(left.path.cmp(&right.path)));
    nodes.dedup_by(|left, right| left.id == right.id && left.kind == right.kind);
    edges.sort_by(|left, right| {
        left.from.cmp(&right.from).then(left.to.cmp(&right.to)).then(format!("{:?}", left.kind).cmp(&format!("{:?}", right.kind)))
    });

    Ok(GraphOutput { nodes, edges })
}

pub fn export_plan(path: &Path) -> io::Result<PlanOutput> {
    let graph = export_graph(path)?;
    let mut steps = Vec::new();

    let project_nodes: Vec<&GraphNode> = graph.nodes.iter().filter(|node| node.kind == GraphNodeKind::Project).collect();
    let module_nodes: Vec<&GraphNode> = graph.nodes.iter().filter(|node| node.kind == GraphNodeKind::Module).collect();
    let capability_nodes: Vec<&GraphNode> = graph.nodes.iter().filter(|node| node.kind == GraphNodeKind::AgentCapability).collect();
    let workflow_nodes: Vec<&GraphNode> = graph.nodes.iter().filter(|node| node.kind == GraphNodeKind::Workflow).collect();

    let mut project_step_ids = Vec::new();
    for project in project_nodes {
        let step_id = format!("project:{}", project.id);
        project_step_ids.push(step_id.clone());
        steps.push(PlanStep {
            id: step_id,
            phase: PlanPhase::Bootstrap,
            title: format!("Bootstrap project {}", project.title),
            depends_on: Vec::new(),
        });
    }

    let module_dependency_map: HashMap<String, Vec<String>> =
        graph.edges.iter().filter(|edge| edge.kind == GraphEdgeKind::DependsOn).fold(HashMap::new(), |mut map, edge| {
            map.entry(edge.from.clone()).or_default().push(edge.to.clone());
            map
        });

    let mut remaining_modules: Vec<String> = module_nodes.iter().map(|node| node.id.clone()).collect();
    remaining_modules.sort();
    let mut emitted_modules = HashSet::new();
    let mut module_step_ids = Vec::new();

    while !remaining_modules.is_empty() {
        let mut ready = Vec::new();
        for module_id in &remaining_modules {
            let dependencies = module_dependency_map.get(module_id).cloned().unwrap_or_default();
            if dependencies.iter().all(|dependency| emitted_modules.contains(dependency)) {
                ready.push(module_id.clone());
            }
        }

        if ready.is_empty() {
            ready.push(remaining_modules[0].clone());
        }

        for module_id in ready {
            remaining_modules.retain(|candidate| candidate != &module_id);
            emitted_modules.insert(module_id.clone());

            let node = module_nodes.iter().find(|node| node.id == module_id).expect("module node should exist");
            let mut depends_on = project_step_ids.clone();
            if let Some(module_dependencies) = module_dependency_map.get(&module_id) {
                for dependency in module_dependencies {
                    depends_on.push(format!("module:{dependency}"));
                }
            }
            let step_id = format!("module:{module_id}");
            module_step_ids.push(step_id.clone());
            steps.push(PlanStep { id: step_id, phase: PlanPhase::Module, title: format!("Implement module {}", node.title), depends_on });
        }
    }

    let mut capability_step_ids = Vec::new();
    for capability in capability_nodes {
        let step_id = format!("capability:{}", capability.id);
        capability_step_ids.push(step_id.clone());
        steps.push(PlanStep {
            id: step_id,
            phase: PlanPhase::AgentCapability,
            title: format!("Define capability {}", capability.title),
            depends_on: module_step_ids.clone(),
        });
    }

    for workflow in workflow_nodes {
        steps.push(PlanStep {
            id: format!("workflow:{}", workflow.id.trim_start_matches("workflow:")),
            phase: PlanPhase::Workflow,
            title: format!("Plan workflow {}", workflow.title),
            depends_on: module_step_ids.iter().chain(capability_step_ids.iter()).cloned().collect(),
        });
    }

    Ok(PlanOutput { steps })
}

pub fn generate_scaffold(path: &Path, output_dir: &Path) -> io::Result<ScaffoldOutput> {
    let validation = validate_findings(path)?;
    if !validation.valid {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            format!("scaffold generation requires a validation-clean tree; found {} validation finding(s)", validation.findings.len()),
        ));
    }

    ensure_output_dir_is_empty(output_dir)?;

    let plan = export_plan(path)?;
    let documents = parse_spec_path(path)?;
    let Some(project) = documents.iter().find_map(|document| match &document.document {
        FastSpecDocument::Project(project) => Some(project),
        _ => None,
    }) else {
        return Err(io::Error::new(io::ErrorKind::InvalidData, "scaffold generation requires a project document"));
    };

    let module_documents: HashMap<String, &fastspec_model::ModuleSpecDocument> = documents
        .iter()
        .filter_map(|document| match &document.document {
            FastSpecDocument::Module(module) => Some((module.metadata.id.clone(), module)),
            _ => None,
        })
        .collect();

    let capability_documents: HashMap<String, &fastspec_model::AgentCapabilitySpecDocument> = documents
        .iter()
        .filter_map(|document| match &document.document {
            FastSpecDocument::AgentCapability(cap) => Some((cap.metadata.id.clone(), cap)),
            _ => None,
        })
        .collect();

    let workflow_document_map: HashMap<String, &fastspec_model::WorkflowSpecDocument> = documents
        .iter()
        .filter_map(|document| match &document.document {
            FastSpecDocument::Workflow(workflow) => Some((workflow.metadata.id.clone(), workflow)),
            _ => None,
        })
        .collect();
    fs::create_dir_all(output_dir)?;

    let mut artifacts = Vec::new();
    record_directory(&mut artifacts, output_dir.to_path_buf(), "scaffold root".to_string());

    let modules_dir = output_dir.join("modules");
    fs::create_dir_all(&modules_dir)?;
    record_directory(&mut artifacts, modules_dir.clone(), "module scaffold directory".to_string());

    let capabilities_dir = output_dir.join("capabilities");
    if !project.spec.agent_capabilities.is_empty() {
        fs::create_dir_all(&capabilities_dir)?;
        record_directory(&mut artifacts, capabilities_dir.clone(), "capability scaffold directory".to_string());
    }

    let workflows_dir = output_dir.join("workflows");
    fs::create_dir_all(&workflows_dir)?;
    record_directory(&mut artifacts, workflows_dir.clone(), "workflow scaffold directory".to_string());

    write_artifact_file(
        &mut artifacts,
        output_dir.join("README.md"),
        render_project_readme(project, &plan),
        "project scaffold overview".to_string(),
    )?;

    for module in &project.spec.modules {
        let module_dir = modules_dir.join(&module.id);
        fs::create_dir_all(&module_dir)?;
        record_directory(&mut artifacts, module_dir.clone(), format!("module directory for {}", module.id));

        let module_document = module_documents.get(&module.id).ok_or_else(|| {
            io::Error::new(io::ErrorKind::InvalidData, format!("validated module `{}` should have a matching document", module.id))
        })?;

        write_artifact_file(
            &mut artifacts,
            module_dir.join("README.md"),
            render_module_readme(module_document),
            format!("module stub for {}", module.id),
        )?;
    }

    for workflow in &project.spec.workflows {
        let workflow_document = workflow_document_map.get(&workflow.id).ok_or_else(|| {
            io::Error::new(
                io::ErrorKind::InvalidData,
                format!("validated workflow `{}` should have a matching document", workflow.id),
            )
        })?;
        let workflow_dir = workflows_dir.join(&workflow.id);
        fs::create_dir_all(&workflow_dir)?;
        record_directory(&mut artifacts, workflow_dir.clone(), format!("workflow directory for {}", workflow.id));
        write_artifact_file(
            &mut artifacts,
            workflow_dir.join("README.md"),
            render_workflow_readme(workflow_document),
            format!("workflow stub for {}", workflow.id),
        )?;
    }
    for capability in &project.spec.agent_capabilities {
        let cap_document = capability_documents.get(&capability.id).ok_or_else(|| {
            io::Error::new(io::ErrorKind::InvalidData, format!("validated capability `{}` should have a matching document", capability.id))
        })?;
        write_artifact_file(
            &mut artifacts,
            capabilities_dir.join(format!("{}.md", capability.id)),
            render_capability_stub(cap_document),
            format!("capability stub for {}", capability.id),
        )?;
    }

    let manifest_path = output_dir.join("fastspec-manifest.json");
    let manifest_artifact = GeneratedArtifact {
        path: manifest_path.clone(),
        kind: GeneratedArtifactKind::File,
        description: "machine-readable scaffold manifest".to_string(),
    };
    let mut manifest_artifacts = artifacts.clone();
    manifest_artifacts.push(manifest_artifact.clone());
    let manifest_json =
        serde_json::to_string_pretty(&ScaffoldOutput { output_dir: output_dir.to_path_buf(), artifacts: manifest_artifacts })
            .map_err(|error| io::Error::other(format!("failed to serialize scaffold manifest: {error}")))?;
    fs::write(&manifest_path, manifest_json)?;
    artifacts.push(manifest_artifact);

    Ok(ScaffoldOutput { output_dir: output_dir.to_path_buf(), artifacts })
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct InitOptions {
    pub id: String,
    pub title: String,
    pub modules: Vec<String>,
    pub capabilities: Vec<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
pub struct InitOutput {
    pub dir: PathBuf,
    pub artifacts: Vec<GeneratedArtifact>,
}

/// Scaffold a new FastSpec project tree at `dir` using the provided options.
/// Fails if `dir` already exists and is non-empty.
pub fn init_spec_tree(dir: &Path, opts: InitOptions) -> io::Result<InitOutput> {
    ensure_output_dir_is_empty(dir)?;
    fs::create_dir_all(dir)?;

    let mut artifacts = Vec::new();
    record_directory(&mut artifacts, dir.to_path_buf(), "spec root directory".to_string());

    // Write project spec.
    let project_yaml = render_init_project(&opts);
    write_artifact_file(&mut artifacts, dir.join("project.fastspec.yaml"), project_yaml, "project spec".to_string())?;

    // Write module specs.
    if !opts.modules.is_empty() {
        let modules_dir = dir.join("modules");
        fs::create_dir_all(&modules_dir)?;
        record_directory(&mut artifacts, modules_dir.clone(), "modules directory".to_string());
        for module_id in &opts.modules {
            write_artifact_file(
                &mut artifacts,
                modules_dir.join(format!("{module_id}.fastspec.yaml")),
                render_init_module(module_id),
                format!("module spec for {module_id}"),
            )?;
        }
    }

    // Write capability specs.
    if !opts.capabilities.is_empty() {
        let capabilities_dir = dir.join("capabilities");
        fs::create_dir_all(&capabilities_dir)?;
        record_directory(&mut artifacts, capabilities_dir.clone(), "capabilities directory".to_string());
        for cap_id in &opts.capabilities {
            write_artifact_file(
                &mut artifacts,
                capabilities_dir.join(format!("{cap_id}.fastspec.yaml")),
                render_init_capability(cap_id),
                format!("capability spec for {cap_id}"),
            )?;
        }
    }

    Ok(InitOutput { dir: dir.to_path_buf(), artifacts })
}

fn render_init_project(opts: &InitOptions) -> String {
    let title = if opts.title.is_empty() { opts.id.clone() } else { opts.title.clone() };
    let module_lines: String = if opts.modules.is_empty() {
        String::new()
    } else {
        let items = opts.modules.iter().map(|id| format!("    - id: {id}\n      purpose: TODO")).collect::<Vec<_>>().join("\n");
        format!("  modules:\n{items}\n")
    };
    let capability_lines: String = if opts.capabilities.is_empty() {
        String::new()
    } else {
        let items = opts.capabilities.iter().map(|id| format!("    - id: {id}\n      purpose: TODO")).collect::<Vec<_>>().join("\n");
        format!("  agentCapabilities:\n{items}\n")
    };
    format!(
        "apiVersion: fastspec.dev/v0alpha1\nkind: ProjectSpec\nmetadata:\n  id: {}\n  title: {}\n  summary: TODO\nspec:\n  goals:\n    - TODO\n{module_lines}{capability_lines}",
        opts.id, title
    )
}

fn render_init_module(id: &str) -> String {
    format!(
        "apiVersion: fastspec.dev/v0alpha1\nkind: ModuleSpec\nmetadata:\n  id: {id}\n  title: {id}\n  summary: TODO\nspec:\n  purpose: TODO\n"
    )
}

fn render_init_capability(id: &str) -> String {
    format!(
        "apiVersion: fastspec.dev/v0alpha1\nkind: AgentCapabilitySpec\nmetadata:\n  id: {id}\n  title: {id}\n  summary: TODO\nspec:\n  goal: TODO\n"
    )
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

fn ensure_output_dir_is_empty(output_dir: &Path) -> io::Result<()> {
    if output_dir.exists() {
        let mut entries = fs::read_dir(output_dir)?;
        if entries.next().transpose()?.is_some() {
            return Err(io::Error::new(
                io::ErrorKind::AlreadyExists,
                format!("output directory {} must not already contain files", output_dir.display()),
            ));
        }
    }

    Ok(())
}

fn record_directory(artifacts: &mut Vec<GeneratedArtifact>, path: PathBuf, description: String) {
    artifacts.push(GeneratedArtifact { path, kind: GeneratedArtifactKind::Directory, description });
}

fn write_artifact_file(artifacts: &mut Vec<GeneratedArtifact>, path: PathBuf, contents: String, description: String) -> io::Result<()> {
    fs::write(&path, contents)?;
    artifacts.push(GeneratedArtifact { path, kind: GeneratedArtifactKind::File, description });
    Ok(())
}

fn render_project_readme(project: &ProjectSpecDocument, plan: &PlanOutput) -> String {
    let module_lines = if project.spec.modules.is_empty() {
        "- none".to_string()
    } else {
        project.spec.modules.iter().map(|module| format!("- `{}`: {}", module.id, module.purpose)).collect::<Vec<_>>().join("\n")
    };

    let capability_lines = if project.spec.agent_capabilities.is_empty() {
        "- none".to_string()
    } else {
        project.spec.agent_capabilities.iter().map(|cap| format!("- `{}`: {}", cap.id, cap.purpose)).collect::<Vec<_>>().join("\n")
    };

    let workflow_lines = if project.spec.workflows.is_empty() {
        "- none".to_string()
    } else {
        project.spec.workflows.iter().map(|workflow| format!("- `{}`: {}", workflow.id, workflow.purpose)).collect::<Vec<_>>().join("\n")
    };

    let plan_lines = if plan.steps.is_empty() {
        "- none".to_string()
    } else {
        plan.steps.iter().map(|step| format!("- `{}`: {}", step.id, step.title)).collect::<Vec<_>>().join("\n")
    };

    format!(
        "# {}\n\n{}\n\n## Modules\n{}\n\n## Agent Capabilities\n{}\n\n## Workflows\n{}\n\n## Plan Steps\n{}\n",
        project.metadata.title, project.metadata.summary, module_lines, capability_lines, workflow_lines, plan_lines
    )
}

fn render_module_readme(module: &fastspec_model::ModuleSpecDocument) -> String {
    let input_lines = if module.spec.inputs.is_empty() {
        "- none".to_string()
    } else {
        module.spec.inputs.iter().map(|input| format!("- `{}`: {}", input.name, input.description)).collect::<Vec<_>>().join("\n")
    };
    let output_lines = if module.spec.outputs.is_empty() {
        "- none".to_string()
    } else {
        module.spec.outputs.iter().map(|output| format!("- `{}`: {}", output.name, output.description)).collect::<Vec<_>>().join("\n")
    };
    let dependency_lines = if module.spec.dependencies.is_empty() {
        "- none".to_string()
    } else {
        module
            .spec
            .dependencies
            .iter()
            .map(|dependency| format!("- `{}`: {}", dependency.id, dependency.reason))
            .collect::<Vec<_>>()
            .join("\n")
    };

    format!(
        "# {}\n\n{}\n\n## Purpose\n{}\n\n## Inputs\n{}\n\n## Outputs\n{}\n\n## Dependencies\n{}\n",
        module.metadata.title, module.metadata.summary, module.spec.purpose, input_lines, output_lines, dependency_lines
    )
}

fn render_workflow_readme(workflow: &fastspec_model::WorkflowSpecDocument) -> String {
    let step_lines = if workflow.spec.steps.is_empty() {
        "- none".to_string()
    } else {
        workflow.spec.steps.iter().map(|step| format!("- `{}`: {}", step.name, step.description)).collect::<Vec<_>>().join("\n")
    };
    let input_lines = if workflow.spec.inputs.is_empty() {
        "- none".to_string()
    } else {
        workflow.spec.inputs.iter().map(|input| format!("- `{}`: {}", input.name, input.description)).collect::<Vec<_>>().join("\n")
    };
    let output_lines = if workflow.spec.outputs.is_empty() {
        "- none".to_string()
    } else {
        workflow.spec.outputs.iter().map(|output| format!("- `{}`: {}", output.name, output.description)).collect::<Vec<_>>().join("\n")
    };
    let trigger_lines = if workflow.spec.triggers.is_empty() {
        "- none".to_string()
    } else {
        workflow.spec.triggers.iter().map(|trigger| format!("- {trigger}")).collect::<Vec<_>>().join("\n")
    };
    format!(
        "# Workflow: {}\n\n{}\n\n## Purpose\n{}\n\n## Steps\n{}\n\n## Inputs\n{}\n\n## Outputs\n{}\n\n## Triggers\n{}\n",
        workflow.metadata.title,
        workflow.metadata.summary,
        workflow.spec.purpose,
        step_lines,
        input_lines,
        output_lines,
        trigger_lines
    )
}

fn render_capability_stub(cap: &fastspec_model::AgentCapabilitySpecDocument) -> String {
    let context_lines = if cap.spec.required_context.is_empty() {
        "- none".to_string()
    } else {
        cap.spec.required_context.iter().map(|ctx| format!("- {ctx}")).collect::<Vec<_>>().join("\n")
    };
    let tool_lines = if cap.spec.allowed_tools.is_empty() {
        "- none".to_string()
    } else {
        cap.spec.allowed_tools.iter().map(|tool| format!("- {tool}")).collect::<Vec<_>>().join("\n")
    };
    let disallowed_lines = if cap.spec.disallowed_actions.is_empty() {
        "- none".to_string()
    } else {
        cap.spec.disallowed_actions.iter().map(|action| format!("- {action}")).collect::<Vec<_>>().join("\n")
    };
    let signal_lines = if cap.spec.success_signals.is_empty() {
        "- none".to_string()
    } else {
        cap.spec.success_signals.iter().map(|signal| format!("- {signal}")).collect::<Vec<_>>().join("\n")
    };
    format!(
        "# Capability: {}\n\n{}\n\n## Goal\n{}\n\n## Required Context\n{}\n\n## Allowed Tools\n{}\n\n## Disallowed Actions\n{}\n\n## Success Signals\n{}\n",
        cap.metadata.title, cap.metadata.summary, cap.spec.goal, context_lines, tool_lines, disallowed_lines, signal_lines
    )
}

#[cfg(test)]
mod tests {
    use std::fs;
    use std::path::PathBuf;
    use std::time::{SystemTime, UNIX_EPOCH};

    use fastspec_model::SpecKind;

    use super::{
        GeneratedArtifactKind, export_graph, export_plan, generate_scaffold, parse_spec_file, validate_findings, validate_spec_tree,
    };

    #[test]
    fn validates_archlint_example_tree() {
        let root = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../examples/archlint-reproduction/specs");
        let summaries = validate_spec_tree(&root).expect("example specs should validate");

        assert_eq!(summaries.len(), 6);
        assert!(summaries.iter().any(|summary| summary.kind == SpecKind::Project));
        assert_eq!(summaries.iter().filter(|summary| summary.kind == SpecKind::Module).count(), 2);
        assert_eq!(summaries.iter().filter(|summary| summary.kind == SpecKind::AgentCapability).count(), 1);
        assert_eq!(summaries.iter().filter(|summary| summary.kind == SpecKind::Workflow).count(), 2);
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

    #[test]
    fn exports_graph_for_example_tree() {
        let root = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../examples/archlint-reproduction/specs");
        let graph = export_graph(&root).expect("example graph should export");
        assert!(graph.nodes.iter().any(|node| node.id == "archlint-reproduction"));
        assert!(graph.nodes.iter().any(|node| node.id == "api"));
        assert!(graph.nodes.iter().any(|node| node.id == "web"));
        assert!(graph.nodes.iter().any(|node| node.id == "workflow:plan"));
        assert!(graph.nodes.iter().any(|node| node.id == "lint-agent"));
        assert!(graph.edges.iter().any(|edge| edge.from == "archlint-reproduction" && edge.to == "api"));
        assert!(graph.edges.iter().any(|edge| edge.from == "web" && edge.to == "api"));
        assert!(graph.edges.iter().any(|edge| edge.from == "archlint-reproduction" && edge.to == "lint-agent"));
    }

    #[test]
    fn rejects_graph_export_for_invalid_tree() {
        let root = unique_temp_dir("invalid-graph-fixtures");
        fs::create_dir_all(root.join("modules")).expect("fixture directories should be created");
        fs::write(
            root.join("project.fastspec.yaml"),
            "apiVersion: fastspec.dev/v0alpha1\nkind: ProjectSpec\nmetadata:\n  id: demo\n  title: Demo\n  summary: Demo project\nspec:\n  modules:\n    - id: api\n      purpose: API module\n",
        )
        .expect("project fixture should write");

        let error = export_graph(&root).expect_err("invalid tree should block graph export");
        assert!(error.to_string().contains("validation-clean tree"));

        fs::remove_dir_all(root).expect("fixture dir should be removed");
    }

    #[test]
    fn exports_plan_for_example_tree() {
        let root = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../examples/archlint-reproduction/specs");
        let plan = export_plan(&root).expect("example plan should export");
        assert!(plan.steps.iter().any(|step| step.id == "project:archlint-reproduction"));
        assert!(plan.steps.iter().any(|step| step.id == "module:api"));
        assert!(
            plan.steps.iter().any(|step| step.id == "module:web" && step.depends_on.iter().any(|dependency| dependency == "module:api"))
        );
        assert!(plan.steps.iter().any(|step| step.id == "capability:lint-agent"));
        assert!(plan.steps.iter().any(|step| step.id == "workflow:plan"));
    }

    #[test]
    fn generates_scaffold_for_example_tree() {
        let root = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../examples/archlint-reproduction/specs");
        let output_dir = unique_temp_dir("generated-scaffold");

        let output = generate_scaffold(&root, &output_dir).expect("scaffold should generate");

        assert_eq!(output.output_dir, output_dir);
        assert!(output.artifacts.iter().any(|artifact| artifact.path.ends_with("README.md")));
        assert!(output.artifacts.iter().any(|artifact| artifact.path.ends_with("modules/api/README.md")));
        assert!(output.artifacts.iter().any(|artifact| artifact.path.ends_with("workflows/plan/README.md")));
        assert!(output.artifacts.iter().any(|artifact| artifact.path.ends_with("capabilities/lint-agent.md")));
        assert!(
            output
                .artifacts
                .iter()
                .any(|artifact| { artifact.path.ends_with("fastspec-manifest.json") && artifact.kind == GeneratedArtifactKind::File })
        );

        let project_readme = fs::read_to_string(output_dir.join("README.md")).expect("project readme should exist");
        assert!(project_readme.contains("Archlint Reproduction"));
        assert!(project_readme.contains("`module:api`"));
        assert!(project_readme.contains("`capability:lint-agent`"));
        fs::remove_dir_all(output_dir).expect("output dir should be removed");
    }

    #[test]
    fn rejects_generation_for_non_empty_output_dir() {
        let root = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../../examples/archlint-reproduction/specs");
        let output_dir = unique_temp_dir("non-empty-scaffold");
        fs::create_dir_all(&output_dir).expect("output dir should be created");
        fs::write(output_dir.join("existing.txt"), "occupied").expect("existing file should write");

        let error = generate_scaffold(&root, &output_dir).expect_err("non-empty output dir should fail");
        assert!(error.to_string().contains("must not already contain files"));

        fs::remove_dir_all(output_dir).expect("output dir should be removed");
    }

    fn unique_temp_file(name: &str) -> PathBuf {
        let unique = SystemTime::now().duration_since(UNIX_EPOCH).expect("time should move forward").as_nanos();
        std::env::temp_dir().join(format!("fastspec-{unique}-{name}"))
    }

    #[test]
    fn detects_module_dependency_cycle() {
        let root = unique_temp_dir("cycle-fixtures");
        fs::create_dir_all(root.join("modules")).expect("fixture directories should be created");

        fs::write(
            root.join("project.fastspec.yaml"),
            "apiVersion: fastspec.dev/v0alpha1\nkind: ProjectSpec\nmetadata:\n  id: demo\n  title: Demo\n  summary: Demo project\nspec:\n  modules:\n    - id: a\n      purpose: Module A\n    - id: b\n      purpose: Module B\n",
        )
        .expect("project fixture should write");
        // a depends on b, b depends on a — direct cycle
        fs::write(
            root.join("modules/a.fastspec.yaml"),
            "apiVersion: fastspec.dev/v0alpha1\nkind: ModuleSpec\nmetadata:\n  id: a\n  title: A\n  summary: Module A\nspec:\n  purpose: Does A\n  dependencies:\n    - id: b\n      reason: Needs B\n",
        )
        .expect("module A fixture should write");
        fs::write(
            root.join("modules/b.fastspec.yaml"),
            "apiVersion: fastspec.dev/v0alpha1\nkind: ModuleSpec\nmetadata:\n  id: b\n  title: B\n  summary: Module B\nspec:\n  purpose: Does B\n  dependencies:\n    - id: a\n      reason: Needs A\n",
        )
        .expect("module B fixture should write");

        let output = validate_findings(&root).expect("validation should run");
        assert!(!output.valid, "cycle should produce a finding");
        assert!(
            output.findings.iter().any(|finding| finding.code == "module_dependency_cycle"),
            "expected module_dependency_cycle finding, got: {:?}",
            output.findings
        );

        fs::remove_dir_all(root).expect("fixture dir should be removed");
    }

    #[test]
    fn reports_workflow_document_validation_findings() {
        let root = unique_temp_dir("workflow-validation-fixtures");
        fs::create_dir_all(root.join("workflows")).expect("fixture directories should be created");

        // Project declares 'my-workflow' but no workflow doc exists → missing_workflow_document.
        fs::write(
            root.join("project.fastspec.yaml"),
            "apiVersion: fastspec.dev/v0alpha1\nkind: ProjectSpec\nmetadata:\n  id: demo\n  title: Demo\n  summary: Demo project\nspec:\n  workflows:\n    - id: my-workflow\n      purpose: TODO\n",
        )
        .expect("project fixture should write");
        // An undeclared workflow doc → undeclared_workflow_document.
        fs::write(
            root.join("workflows/ghost.fastspec.yaml"),
            "apiVersion: fastspec.dev/v0alpha1\nkind: WorkflowSpec\nmetadata:\n  id: ghost\n  title: Ghost Workflow\n  summary: Not declared\nspec:\n  purpose: Hidden\n",
        )
        .expect("ghost workflow fixture should write");

        let output = validate_findings(&root).expect("validation should run");
        assert!(!output.valid);
        assert!(output.findings.iter().any(|f| f.code == "missing_workflow_document"), "got: {:?}", output.findings);
        assert!(output.findings.iter().any(|f| f.code == "undeclared_workflow_document"), "got: {:?}", output.findings);

        fs::remove_dir_all(root).expect("fixture dir should be removed");
    }

    fn unique_temp_dir(name: &str) -> PathBuf {
        let unique = SystemTime::now().duration_since(UNIX_EPOCH).expect("time should move forward").as_nanos();
        std::env::temp_dir().join(format!("fastspec-{unique}-{name}"))
    }
}
