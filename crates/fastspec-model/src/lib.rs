use std::fmt;

use serde::{Deserialize, Serialize};
use serde_yaml::Value;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize)]
pub enum SpecKind {
    Project,
    Module,
    AgentCapability,
    Unknown,
}

impl SpecKind {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::Project => "ProjectSpec",
            Self::Module => "ModuleSpec",
            Self::AgentCapability => "AgentCapabilitySpec",
            Self::Unknown => "Unknown",
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct ParseError {
    message: String,
}

impl ParseError {
    fn new(message: impl Into<String>) -> Self {
        Self { message: message.into() }
    }
}

impl fmt::Display for ParseError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(&self.message)
    }
}

impl std::error::Error for ParseError {}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize, Serialize)]
pub struct Metadata {
    pub id: String,
    pub title: String,
    pub summary: String,
    #[serde(default)]
    pub owners: Vec<String>,
    #[serde(default)]
    pub tags: Vec<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize, Serialize)]
pub struct NamedItem {
    pub name: String,
    pub description: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize, Serialize)]
pub struct IdPurpose {
    pub id: String,
    pub purpose: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize, Serialize)]
pub struct Dependency {
    pub id: String,
    pub reason: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize, Serialize)]
pub struct ProjectSpecBody {
    #[serde(default)]
    pub goals: Vec<String>,
    #[serde(rename = "nonGoals", default)]
    pub non_goals: Vec<String>,
    #[serde(default)]
    pub constraints: Vec<String>,
    #[serde(default)]
    pub modules: Vec<IdPurpose>,
    #[serde(rename = "agentCapabilities", default)]
    pub agent_capabilities: Vec<IdPurpose>,
    #[serde(default)]
    pub workflows: Vec<IdPurpose>,
}
#[derive(Debug, Clone, PartialEq, Eq, Deserialize, Serialize)]
pub struct ModuleSpecBody {
    pub purpose: String,
    #[serde(default)]
    pub inputs: Vec<NamedItem>,
    #[serde(default)]
    pub outputs: Vec<NamedItem>,
    #[serde(default)]
    pub dependencies: Vec<Dependency>,
    #[serde(default)]
    pub invariants: Vec<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize, Serialize)]
pub struct AgentCapabilitySpecBody {
    pub goal: String,
    #[serde(rename = "requiredContext", default)]
    pub required_context: Vec<String>,
    #[serde(rename = "allowedTools", default)]
    pub allowed_tools: Vec<String>,
    #[serde(rename = "disallowedActions", default)]
    pub disallowed_actions: Vec<String>,
    #[serde(rename = "successSignals", default)]
    pub success_signals: Vec<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize, Serialize)]
pub struct ProjectSpecDocument {
    #[serde(rename = "apiVersion")]
    pub api_version: String,
    pub kind: String,
    pub metadata: Metadata,
    pub spec: ProjectSpecBody,
}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize, Serialize)]
pub struct ModuleSpecDocument {
    #[serde(rename = "apiVersion")]
    pub api_version: String,
    pub kind: String,
    pub metadata: Metadata,
    pub spec: ModuleSpecBody,
}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize, Serialize)]
pub struct AgentCapabilitySpecDocument {
    #[serde(rename = "apiVersion")]
    pub api_version: String,
    pub kind: String,
    pub metadata: Metadata,
    pub spec: AgentCapabilitySpecBody,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(tag = "kind", content = "document")]
pub enum FastSpecDocument {
    Project(ProjectSpecDocument),
    Module(ModuleSpecDocument),
    AgentCapability(AgentCapabilitySpecDocument),
}

impl FastSpecDocument {
    pub fn kind(&self) -> SpecKind {
        match self {
            Self::Project(_) => SpecKind::Project,
            Self::Module(_) => SpecKind::Module,
            Self::AgentCapability(_) => SpecKind::AgentCapability,
        }
    }

    pub fn metadata(&self) -> &Metadata {
        match self {
            Self::Project(document) => &document.metadata,
            Self::Module(document) => &document.metadata,
            Self::AgentCapability(document) => &document.metadata,
        }
    }

    pub fn spec_detail_lines(&self) -> Vec<String> {
        match self {
            Self::Project(document) => vec![
                format!("goals: {}", document.spec.goals.len()),
                format!("modules: {}", document.spec.modules.len()),
                format!("agent_capabilities: {}", document.spec.agent_capabilities.len()),
                format!("workflows: {}", document.spec.workflows.len()),
            ],
            Self::Module(document) => vec![
                format!("purpose: {}", document.spec.purpose),
                format!("inputs: {}", document.spec.inputs.len()),
                format!("outputs: {}", document.spec.outputs.len()),
            ],
            Self::AgentCapability(document) => vec![
                format!("goal: {}", document.spec.goal),
                format!("required_context: {}", document.spec.required_context.len()),
                format!("allowed_tools: {}", document.spec.allowed_tools.len()),
            ],
        }
    }
}

pub fn detect_kind(contents: &str) -> SpecKind {
    match parse_document(contents) {
        Ok(document) => document.kind(),
        Err(_) => SpecKind::Unknown,
    }
}

pub fn parse_document(contents: &str) -> Result<FastSpecDocument, ParseError> {
    let value: Value = serde_yaml::from_str(contents).map_err(|error| ParseError::new(error.to_string()))?;
    let kind = value.get("kind").and_then(Value::as_str).ok_or_else(|| ParseError::new("missing required field `kind`"))?;

    match kind {
        "ProjectSpec" => serde_yaml::from_value::<ProjectSpecDocument>(value)
            .map(FastSpecDocument::Project)
            .map_err(|error| ParseError::new(error.to_string())),
        "ModuleSpec" => serde_yaml::from_value::<ModuleSpecDocument>(value)
            .map(FastSpecDocument::Module)
            .map_err(|error| ParseError::new(error.to_string())),
        "AgentCapabilitySpec" => serde_yaml::from_value::<AgentCapabilitySpecDocument>(value)
            .map(FastSpecDocument::AgentCapability)
            .map_err(|error| ParseError::new(error.to_string())),
        other => Err(ParseError::new(format!("unsupported kind `{other}`"))),
    }
}

#[cfg(test)]
mod tests {
    use super::{FastSpecDocument, SpecKind, detect_kind, parse_document};

    #[test]
    fn detects_project_spec() {
        let contents = "apiVersion: fastspec.dev/v0alpha1\nkind: ProjectSpec\nmetadata:\n  id: demo\n  title: Demo\n  summary: Demo project\nspec:\n  goals:\n    - Ship it\n";
        assert_eq!(detect_kind(contents), SpecKind::Project);
    }

    #[test]
    fn parses_module_spec() {
        let contents = "apiVersion: fastspec.dev/v0alpha1\nkind: ModuleSpec\nmetadata:\n  id: demo-module\n  title: Demo Module\n  summary: Demo summary\nspec:\n  purpose: Handle demo requests\n";
        let document = parse_document(contents).expect("module should parse");
        match document {
            FastSpecDocument::Module(document) => {
                assert_eq!(document.metadata.id, "demo-module");
                assert_eq!(document.spec.purpose, "Handle demo requests");
            }
            other => panic!("expected module, got {other:?}"),
        }
    }

    #[test]
    fn returns_unknown_for_missing_kind() {
        let contents = "apiVersion: fastspec.dev/v0alpha1\nmetadata:\n  id: demo\n";
        assert_eq!(detect_kind(contents), SpecKind::Unknown);
    }

    #[test]
    fn rejects_missing_required_metadata_fields() {
        let contents = "apiVersion: fastspec.dev/v0alpha1\nkind: ProjectSpec\nmetadata:\n  title: Demo\n  summary: Demo project\nspec:\n  goals:\n    - Ship it\n";
        let error = parse_document(contents).expect_err("missing id should fail");
        assert!(error.to_string().contains("id"));
    }
}
