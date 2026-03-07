use std::fs;
use std::path::{Path, PathBuf};
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

fn write_invalid_fixture_tree(root: &Path) {
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
}

#[test]
fn validate_reports_clean_example_tree() {
    let output = cli_command()
        .args(["validate", &workspace_path("examples/archlint-reproduction/specs")])
        .output()
        .expect("validate command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));
    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.contains("valid: true"));
}

#[test]
fn validate_json_reports_findings_for_invalid_tree() {
    let root = unique_temp_dir("validate-json-fixture");
    write_invalid_fixture_tree(&root);

    let output = cli_command().args(["validate", "--json", &root.display().to_string()]).output().expect("validate command should run");

    assert!(!output.status.success(), "command should fail on invalid tree");
    let value: serde_json::Value = serde_json::from_slice(&output.stdout).expect("stdout should be json");
    assert_eq!(value["valid"], false);
    let findings = value["findings"].as_array().expect("findings should be an array");
    assert!(findings.iter().any(|finding| finding["code"] == "missing_module_document"));
    assert!(findings.iter().any(|finding| finding["code"] == "undeclared_module_document"));
    assert!(findings.iter().any(|finding| finding["code"] == "invalid_module_dependency"));

    fs::remove_dir_all(root).expect("fixture dir should be removed");
}
