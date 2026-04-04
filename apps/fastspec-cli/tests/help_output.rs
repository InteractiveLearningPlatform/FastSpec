use std::process::Command;

fn cli_command() -> Command {
    Command::new(env!("CARGO_BIN_EXE_fastspec"))
}

#[test]
fn version_flag_prints_version() {
    let output = cli_command().arg("--version").output().expect("command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));
    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.starts_with("fastspec "), "expected 'fastspec <version>', got: {stdout}");
    // Version should have at least one dot (e.g. "0.1.0").
    assert!(stdout.contains('.'), "expected semver with dot: {stdout}");
}

#[test]
fn short_version_flag_prints_version() {
    let output = cli_command().arg("-V").output().expect("command should run");

    assert!(output.status.success());
    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.starts_with("fastspec "));
}

#[test]
fn help_flag_prints_help() {
    let output = cli_command().arg("--help").output().expect("command should run");

    assert!(output.status.success(), "stderr: {}", String::from_utf8_lossy(&output.stderr));
    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.contains("fastspec"), "expected fastspec in help");
    assert!(stdout.contains("Commands:"), "expected commands section");
    assert!(stdout.contains("validate"), "expected validate command listed");
    assert!(stdout.contains("graph"), "expected graph command listed");
    assert!(stdout.contains("init"), "expected init command listed");
}

#[test]
fn short_help_flag_prints_help() {
    let output = cli_command().arg("-h").output().expect("command should run");

    assert!(output.status.success());
    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.contains("Commands:"));
}

#[test]
fn no_args_prints_help() {
    let output = cli_command().output().expect("command should run");

    assert!(output.status.success(), "no-arg invocation should succeed with help output");
    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.contains("Commands:"), "expected help output when no args given");
}

#[test]
fn per_command_help_graph() {
    let output = cli_command().args(["graph", "--help"]).output().expect("command should run");

    assert!(output.status.success());
    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.contains("graph"), "expected graph command help");
    assert!(stdout.contains("--format"), "expected --format flag in graph help");
    assert!(stdout.contains("mermaid"), "expected mermaid format mentioned");
}

#[test]
fn per_command_help_init() {
    let output = cli_command().args(["init", "--help"]).output().expect("command should run");

    assert!(output.status.success());
    let stdout = String::from_utf8(output.stdout).expect("stdout should be utf-8");
    assert!(stdout.contains("init"), "expected init command help");
    assert!(stdout.contains("--modules"), "expected --modules flag in init help");
    assert!(stdout.contains("--capabilities"), "expected --capabilities flag in init help");
}
