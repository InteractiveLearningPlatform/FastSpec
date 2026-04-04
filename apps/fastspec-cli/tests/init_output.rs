use std::fs;
use std::path::PathBuf;
use std::process::Command;
use std::time::{SystemTime, UNIX_EPOCH};

fn cli_command() -> Command {
    Command::new(env!("CARGO_BIN_EXE_fastspec"))
}

fn unique_temp_dir(name: &str) -> PathBuf {
    let unique = SystemTime::now().duration_since(UNIX_EPOCH).expect("time should move forward").as_nanos();
    std::env::temp_dir().join(format!("fastspec-{unique}-{name}"))
}

#[test]
fn init_creates_project_spec() {
    let dir = unique_temp_dir("init-basic");

    let output = cli_command()
        .args(["init", "--id", "demo", "--title", "Demo Project", &dir.display().to_string()])
        .output()
        .expect("init command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));

    let project_yaml = fs::read_to_string(dir.join("project.fastspec.yaml")).expect("project spec should exist");
    assert!(project_yaml.contains("id: demo"), "project id should be set");
    assert!(project_yaml.contains("title: Demo Project"), "project title should be set");
    assert!(project_yaml.contains("kind: ProjectSpec"), "should be a ProjectSpec");

    fs::remove_dir_all(&dir).expect("output dir should be removed");
}

#[test]
fn init_creates_module_specs() {
    let dir = unique_temp_dir("init-modules");

    let output = cli_command()
        .args(["init", "--id", "myapp", "--modules", "api,web", &dir.display().to_string()])
        .output()
        .expect("init command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));

    assert!(dir.join("modules/api.fastspec.yaml").exists(), "api module should be created");
    assert!(dir.join("modules/web.fastspec.yaml").exists(), "web module should be created");

    let api_yaml = fs::read_to_string(dir.join("modules/api.fastspec.yaml")).expect("api module should exist");
    assert!(api_yaml.contains("id: api"), "api module id");
    assert!(api_yaml.contains("kind: ModuleSpec"), "api module kind");

    // Project spec should declare both modules.
    let project_yaml = fs::read_to_string(dir.join("project.fastspec.yaml")).expect("project should exist");
    assert!(project_yaml.contains("id: api"), "project should declare api module");
    assert!(project_yaml.contains("id: web"), "project should declare web module");

    fs::remove_dir_all(&dir).expect("output dir should be removed");
}

#[test]
fn init_creates_capability_specs() {
    let dir = unique_temp_dir("init-capabilities");

    let output = cli_command()
        .args(["init", "--id", "myapp", "--capabilities", "lint-agent", &dir.display().to_string()])
        .output()
        .expect("init command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));

    assert!(dir.join("capabilities/lint-agent.fastspec.yaml").exists(), "capability should be created");

    let cap_yaml = fs::read_to_string(dir.join("capabilities/lint-agent.fastspec.yaml")).expect("capability should exist");
    assert!(cap_yaml.contains("id: lint-agent"), "capability id");
    assert!(cap_yaml.contains("kind: AgentCapabilitySpec"), "capability kind");

    fs::remove_dir_all(&dir).expect("output dir should be removed");
}

#[test]
fn init_json_reports_artifacts() {
    let dir = unique_temp_dir("init-json");

    let output = cli_command()
        .args(["init", "--json", "--id", "demo", "--modules", "core", &dir.display().to_string()])
        .output()
        .expect("init command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));

    let value: serde_json::Value = serde_json::from_slice(&output.stdout).expect("stdout should be json");
    assert_eq!(value["dir"], dir.display().to_string());
    let artifacts = value["artifacts"].as_array().expect("artifacts should be an array");
    assert!(artifacts.iter().any(|a| a["path"].as_str().map_or(false, |p| p.contains("project.fastspec.yaml"))));
    assert!(artifacts.iter().any(|a| a["path"].as_str().map_or(false, |p| p.contains("core.fastspec.yaml"))));

    fs::remove_dir_all(&dir).expect("output dir should be removed");
}

#[test]
fn init_rejects_non_empty_dir() {
    let dir = unique_temp_dir("init-occupied");
    fs::create_dir_all(&dir).expect("dir should be created");
    fs::write(dir.join("existing.txt"), "occupied").expect("existing file should write");

    let output = cli_command().args(["init", "--id", "demo", &dir.display().to_string()]).output().expect("init command should run");

    assert!(!output.status.success(), "init on non-empty dir should fail");
    let stderr = String::from_utf8(output.stderr).expect("stderr should be utf-8");
    assert!(stderr.contains("must not already contain files"), "expected non-empty error");

    fs::remove_dir_all(&dir).expect("fixture dir should be removed");
}

#[test]
fn init_uses_dir_name_as_default_id() {
    // Use a fixed base name so we can predict the derived id.
    let base = unique_temp_dir("");
    let dir = base.join("my-spec-project");

    let output = cli_command().args(["init", &dir.display().to_string()]).output().expect("init command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));

    let project_yaml = fs::read_to_string(dir.join("project.fastspec.yaml")).expect("project spec should exist");
    assert!(project_yaml.contains("id: my-spec-project"), "id should default to directory name: {project_yaml}");

    fs::remove_dir_all(&base).expect("output dir should be removed");
}
