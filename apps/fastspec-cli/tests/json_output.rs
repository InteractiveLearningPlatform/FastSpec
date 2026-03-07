use std::path::PathBuf;
use std::process::Command;

fn cli_command() -> Command {
    Command::new(env!("CARGO_BIN_EXE_fastspec"))
}

fn workspace_path(relative: &str) -> String {
    PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../..").join(relative).display().to_string()
}

#[test]
fn summary_json_outputs_machine_readable_documents() {
    let output = cli_command()
        .args(["summary", "--json", &workspace_path("examples/archlint-reproduction/specs")])
        .output()
        .expect("summary command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));

    let value: serde_json::Value = serde_json::from_slice(&output.stdout).expect("output should be json");
    let documents = value["documents"].as_array().expect("documents should be an array");
    assert_eq!(documents.len(), 3);
    assert!(documents.iter().any(|document| document["id"] == "archlint-reproduction"));
}

#[test]
fn inspect_json_outputs_document_details() {
    let output = cli_command()
        .args(["inspect", "--json", &workspace_path("examples/archlint-reproduction/specs/project.fastspec.yaml")])
        .output()
        .expect("inspect command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));

    let value: serde_json::Value = serde_json::from_slice(&output.stdout).expect("output should be json");
    let documents = value["documents"].as_array().expect("documents should be an array");
    assert_eq!(documents.len(), 1);
    assert_eq!(documents[0]["metadata"]["id"], "archlint-reproduction");
    assert_eq!(documents[0]["kind"], "Project");
}

#[test]
fn summary_json_returns_json_error_for_invalid_input() {
    let output = cli_command().args(["summary", "--json", &workspace_path("missing-path")]).output().expect("summary command should run");

    assert!(!output.status.success(), "command should fail");

    let value: serde_json::Value = serde_json::from_slice(&output.stderr).expect("stderr should be json");
    assert!(value["error"].as_str().expect("error string").contains("No such file or directory"));
}
