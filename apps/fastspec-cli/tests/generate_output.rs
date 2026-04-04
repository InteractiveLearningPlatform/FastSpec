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
        "apiVersion: fastspec.dev/v0alpha1\nkind: ProjectSpec\nmetadata:\n  id: demo\n  title: Demo\n  summary: Demo project\nspec:\n  modules:\n    - id: api\n      purpose: API module\n",
    )
    .expect("project fixture should write");
}

#[test]
fn generate_writes_scaffold_to_output_dir() {
    let output_dir = unique_temp_dir("generate-text");
    let output = cli_command()
        .args(["generate", "--out", &output_dir.display().to_string(), &workspace_path("examples/archlint-reproduction/specs")])
        .output()
        .expect("generate command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));

    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.contains("output_dir:"));
    assert!(stdout.contains("fastspec-manifest.json"));

    let project_readme = fs::read_to_string(output_dir.join("README.md")).expect("project readme should exist");
    assert!(project_readme.contains("Archlint Reproduction"));
    assert!(output_dir.join("modules/api/README.md").exists());
    assert!(output_dir.join("workflows/generate/README.md").exists());
    assert!(output_dir.join("fastspec-manifest.json").exists());

    fs::remove_dir_all(output_dir).expect("generated output should be removed");
}

#[test]
fn generate_json_reports_artifacts() {
    let output_dir = unique_temp_dir("generate-json");
    let output = cli_command()
        .args(["generate", "--json", "--out", &output_dir.display().to_string(), &workspace_path("examples/archlint-reproduction/specs")])
        .output()
        .expect("generate command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));

    let value: serde_json::Value = serde_json::from_slice(&output.stdout).expect("stdout should be json");
    assert_eq!(value["output_dir"], output_dir.display().to_string());
    let artifacts = value["artifacts"].as_array().expect("artifacts should be an array");
    assert!(artifacts.iter().any(|artifact| artifact["path"] == output_dir.join("README.md").display().to_string()));
    assert!(artifacts.iter().any(|artifact| artifact["path"] == output_dir.join("fastspec-manifest.json").display().to_string()));

    fs::remove_dir_all(output_dir).expect("generated output should be removed");
}

#[test]
fn generate_rejects_invalid_tree() {
    let root = unique_temp_dir("generate-invalid-fixture");
    let output_dir = unique_temp_dir("generate-invalid-out");
    write_invalid_fixture_tree(&root);

    let output = cli_command()
        .args(["generate", "--json", "--out", &output_dir.display().to_string(), &root.display().to_string()])
        .output()
        .expect("generate command should run");

    assert!(!output.status.success(), "invalid tree should fail");
    let value: serde_json::Value = serde_json::from_slice(&output.stderr).expect("stderr should be json");
    assert!(value["error"].as_str().expect("error string").contains("validation-clean tree"));
    assert!(!output_dir.exists(), "output directory should not be created on failure");

    fs::remove_dir_all(root).expect("fixture dir should be removed");
}
