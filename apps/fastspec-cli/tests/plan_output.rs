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
fn plan_json_exports_ordered_steps() {
    let output = cli_command()
        .args(["plan", "--json", &workspace_path("examples/archlint-reproduction/specs")])
        .output()
        .expect("plan command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));
    let value: serde_json::Value = serde_json::from_slice(&output.stdout).expect("stdout should be json");
    let steps = value["steps"].as_array().expect("steps should be an array");
    assert!(steps.iter().any(|step| step["id"] == "project:archlint-reproduction"));
    assert!(steps.iter().any(|step| step["id"] == "module:web" && step["depends_on"].as_array().is_some()));
    assert!(steps.iter().any(|step| step["id"] == "workflow:plan"));
}

#[test]
fn plan_rejects_invalid_tree() {
    let root = unique_temp_dir("plan-invalid-fixture");
    fs::create_dir_all(root.join("modules")).expect("fixture directories should be created");
    fs::write(
        root.join("project.fastspec.yaml"),
        "apiVersion: fastspec.dev/v0alpha1\nkind: ProjectSpec\nmetadata:\n  id: demo\n  title: Demo\n  summary: Demo project\nspec:\n  modules:\n    - id: api\n      purpose: API module\n",
    )
    .expect("project fixture should write");

    let output = cli_command().args(["plan", "--json", &root.display().to_string()]).output().expect("plan command should run");

    assert!(!output.status.success(), "invalid tree should fail");
    let value: serde_json::Value = serde_json::from_slice(&output.stderr).expect("stderr should be json");
    assert!(value["error"].as_str().expect("error string").contains("validation-clean tree"));

    fs::remove_dir_all(root).expect("fixture dir should be removed");
}
