use std::fs;
use std::path::PathBuf;
use std::process::Command;
use std::time::{SystemTime, UNIX_EPOCH};

fn cli_command() -> Command {
    Command::new(env!("CARGO_BIN_EXE_fastspec"))
}

fn workspace_path(relative: &str) -> String {
    PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../..").join(relative).display().to_string()
}

fn unique_temp_dir(name: &str) -> PathBuf {
    let unique = SystemTime::now().duration_since(UNIX_EPOCH).expect("time should move forward").as_nanos();
    std::env::temp_dir().join(format!("fastspec-{unique}-{name}"))
}

#[test]
fn graph_json_exports_normalized_graph() {
    let output = cli_command()
        .args(["graph", "--json", &workspace_path("examples/archlint-reproduction/specs")])
        .output()
        .expect("graph command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));
    let value: serde_json::Value = serde_json::from_slice(&output.stdout).expect("stdout should be json");
    let nodes = value["nodes"].as_array().expect("nodes should be an array");
    let edges = value["edges"].as_array().expect("edges should be an array");
    assert!(nodes.iter().any(|node| node["id"] == "archlint-reproduction"));
    assert!(nodes.iter().any(|node| node["id"] == "workflow:plan"));
    assert!(edges.iter().any(|edge| edge["from"] == "web" && edge["to"] == "api" && edge["kind"] == "depends_on"));
}

#[test]
fn graph_rejects_invalid_tree() {
    let root = unique_temp_dir("graph-invalid-fixture");
    fs::create_dir_all(root.join("modules")).expect("fixture directories should be created");
    fs::write(
        root.join("project.fastspec.yaml"),
        "apiVersion: fastspec.dev/v0alpha1\nkind: ProjectSpec\nmetadata:\n  id: demo\n  title: Demo\n  summary: Demo project\nspec:\n  modules:\n    - id: api\n      purpose: API module\n",
    )
    .expect("project fixture should write");

    let output = cli_command().args(["graph", "--json", &root.display().to_string()]).output().expect("graph command should run");

    assert!(!output.status.success(), "invalid tree should fail");
    let value: serde_json::Value = serde_json::from_slice(&output.stderr).expect("stderr should be json");
    assert!(value["error"].as_str().expect("error string").contains("validation-clean tree"));

    fs::remove_dir_all(root).expect("fixture dir should be removed");
}

#[test]
fn graph_mermaid_produces_flowchart() {
    let output = cli_command()
        .args(["graph", "--format", "mermaid", &workspace_path("examples/archlint-reproduction/specs")])
        .output()
        .expect("graph command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));
    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.starts_with("flowchart LR"), "expected Mermaid header, got: {stdout}");
    // Colon-containing IDs (workflow:plan) are sanitized to double-underscore.
    assert!(stdout.contains("workflow__plan"), "expected sanitized workflow id");
    assert!(stdout.contains("archlint-reproduction"), "expected project id");
    // Edges should use arrows.
    assert!(stdout.contains("-->"), "expected edge arrows");
}

#[test]
fn graph_dot_produces_digraph() {
    let output = cli_command()
        .args(["graph", "--format", "dot", &workspace_path("examples/archlint-reproduction/specs")])
        .output()
        .expect("graph command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));
    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.contains("digraph fastspec {"), "expected DOT digraph header");
    assert!(stdout.contains("archlint-reproduction"), "expected project id node");
    assert!(stdout.contains("->"), "expected DOT edge arrow");
    assert!(stdout.trim_end().ends_with("}"), "expected closing brace");
}

#[test]
fn graph_format_flag_rejects_unknown_format() {
    let output = cli_command()
        .args(["graph", "--format", "xml", &workspace_path("examples/archlint-reproduction/specs")])
        .output()
        .expect("graph command should run");

    assert!(!output.status.success(), "unknown format should fail");
    let stderr = String::from_utf8(output.stderr).expect("stderr should be utf-8");
    assert!(stderr.contains("unknown graph format"), "expected format error message");
}