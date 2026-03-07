use std::env;
use std::path::Path;
use std::process::ExitCode;

use fastspec_core::{
    GraphOutput, InspectOutput, PlanOutput, ScaffoldOutput, SummaryOutput, ValidationOutput, export_graph, export_plan, generate_scaffold,
    parse_spec_path, validate_findings, validate_spec_tree,
};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum CommandKind {
    Summary,
    Inspect,
    Validate,
    Graph,
    Plan,
    Generate,
}

#[derive(Debug, Clone, PartialEq, Eq)]
struct CliCommand {
    kind: CommandKind,
    path: String,
    output_dir: Option<String>,
    json: bool,
}

const USAGE: &str = "usage: fastspec <summary|inspect|validate|graph|plan|generate> [--json] [--out <dir>] <path>";

fn main() -> ExitCode {
    match parse_args(env::args().skip(1)) {
        Ok(command) => run_command(command),
        Err(message) => {
            eprintln!("{message}");
            ExitCode::from(2)
        }
    }
}

fn parse_args(args: impl IntoIterator<Item = String>) -> Result<CliCommand, String> {
    let mut args = args.into_iter();
    let kind = match args.next().as_deref() {
        Some("summary") => CommandKind::Summary,
        Some("inspect") => CommandKind::Inspect,
        Some("validate") => CommandKind::Validate,
        Some("graph") => CommandKind::Graph,
        Some("plan") => CommandKind::Plan,
        Some("generate") => CommandKind::Generate,
        _ => return Err(USAGE.to_string()),
    };

    let mut json = false;
    let mut path = None;
    let mut output_dir = None;
    let mut args = args.peekable();
    while let Some(arg) = args.next() {
        if arg == "--json" {
            json = true;
        } else if arg == "--out" {
            let Some(dir) = args.next() else {
                return Err(USAGE.to_string());
            };
            output_dir = Some(dir);
        } else if let Some(dir) = arg.strip_prefix("--out=") {
            if dir.is_empty() {
                return Err(USAGE.to_string());
            }
            output_dir = Some(dir.to_string());
        } else if path.is_none() {
            path = Some(arg);
        } else {
            return Err(USAGE.to_string());
        }
    }

    let Some(path) = path else {
        return Err(USAGE.to_string());
    };

    if kind == CommandKind::Generate && output_dir.is_none() {
        return Err(USAGE.to_string());
    }

    if kind != CommandKind::Generate && output_dir.is_some() {
        return Err(USAGE.to_string());
    }

    Ok(CliCommand { kind, path, output_dir, json })
}

fn run_command(command: CliCommand) -> ExitCode {
    match command.kind {
        CommandKind::Summary => print_summary(Path::new(&command.path), command.json),
        CommandKind::Inspect => inspect_path(Path::new(&command.path), command.json),
        CommandKind::Validate => validate_path(Path::new(&command.path), command.json),
        CommandKind::Graph => graph_path(Path::new(&command.path), command.json),
        CommandKind::Plan => plan_path(Path::new(&command.path), command.json),
        CommandKind::Generate => {
            generate_path(Path::new(&command.path), Path::new(command.output_dir.as_deref().unwrap_or_default()), command.json)
        }
    }
}

fn print_summary(path: &Path, json: bool) -> ExitCode {
    match validate_spec_tree(path) {
        Ok(summaries) => {
            if json {
                return print_json(&SummaryOutput { documents: summaries });
            }

            for summary in summaries {
                println!("{}\t{}\t{}\t{}", summary.kind.as_str(), summary.id, summary.title, summary.path.display());
            }
            ExitCode::SUCCESS
        }
        Err(error) => print_error(&error.to_string(), json),
    }
}

fn inspect_path(path: &Path, json: bool) -> ExitCode {
    match parse_spec_path(path) {
        Ok(documents) => {
            if json {
                let documents = documents.into_iter().map(|document| document.into_inspect()).collect();
                return print_json(&InspectOutput { documents });
            }

            for document in documents {
                println!("path: {}", document.path.display());
                println!("kind: {}", document.document.kind().as_str());
                println!("id: {}", document.document.metadata().id);
                println!("title: {}", document.document.metadata().title);
                println!("summary: {}", document.document.metadata().summary);
                if !document.document.metadata().tags.is_empty() {
                    println!("tags: {}", document.document.metadata().tags.join(", "));
                }
                for detail in document.document.spec_detail_lines() {
                    println!("{detail}");
                }
                println!();
            }
            ExitCode::SUCCESS
        }
        Err(error) => print_error(&error.to_string(), json),
    }
}

fn validate_path(path: &Path, json: bool) -> ExitCode {
    match validate_findings(path) {
        Ok(output) => {
            if json {
                let exit = if output.valid { ExitCode::SUCCESS } else { ExitCode::from(1) };
                return print_json_with_status(&output, exit);
            }

            print_validation_text(&output);
            if output.valid { ExitCode::SUCCESS } else { ExitCode::from(1) }
        }
        Err(error) => print_error(&error.to_string(), json),
    }
}

fn graph_path(path: &Path, json: bool) -> ExitCode {
    match export_graph(path) {
        Ok(output) => {
            if json {
                return print_json(&output);
            }

            print_graph_text(&output);
            ExitCode::SUCCESS
        }
        Err(error) => print_error(&error.to_string(), json),
    }
}

fn plan_path(path: &Path, json: bool) -> ExitCode {
    match export_plan(path) {
        Ok(output) => {
            if json {
                return print_json(&output);
            }

            print_plan_text(&output);
            ExitCode::SUCCESS
        }
        Err(error) => print_error(&error.to_string(), json),
    }
}

fn generate_path(path: &Path, output_dir: &Path, json: bool) -> ExitCode {
    match generate_scaffold(path, output_dir) {
        Ok(output) => {
            if json {
                return print_json(&output);
            }

            print_generation_text(&output);
            ExitCode::SUCCESS
        }
        Err(error) => print_error(&error.to_string(), json),
    }
}

fn print_json<T: serde::Serialize>(value: &T) -> ExitCode {
    print_json_with_status(value, ExitCode::SUCCESS)
}

fn print_json_with_status<T: serde::Serialize>(value: &T, exit: ExitCode) -> ExitCode {
    match serde_json::to_string_pretty(value) {
        Ok(json) => {
            println!("{json}");
            exit
        }
        Err(error) => {
            eprintln!("failed to serialize JSON output: {error}");
            ExitCode::from(1)
        }
    }
}

fn print_error(message: &str, json: bool) -> ExitCode {
    if json {
        match serde_json::to_string_pretty(&serde_json::json!({ "error": message })) {
            Ok(json) => eprintln!("{json}"),
            Err(error) => eprintln!("{{\"error\":\"{message}\",\"serialization_error\":\"{error}\"}}"),
        }
    } else {
        eprintln!("{message}");
    }

    ExitCode::from(1)
}

fn print_validation_text(output: &ValidationOutput) {
    if output.valid {
        println!("valid: true");
        println!("findings: none");
        return;
    }

    println!("valid: false");
    for finding in &output.findings {
        println!("{}\t{}\t{}\t{}", format!("{:?}", finding.severity).to_lowercase(), finding.code, finding.path.display(), finding.message);
    }
}

fn print_graph_text(output: &GraphOutput) {
    println!("nodes:");
    for node in &output.nodes {
        println!("{}\t{:?}\t{}\t{}", node.id, node.kind, node.title, node.path.display());
    }
    println!("edges:");
    for edge in &output.edges {
        println!("{}\t{:?}\t{}", edge.from, edge.kind, edge.to);
    }
}

fn print_plan_text(output: &PlanOutput) {
    for step in &output.steps {
        if step.depends_on.is_empty() {
            println!("{}\t{:?}\t{}", step.id, step.phase, step.title);
        } else {
            println!("{}\t{:?}\t{}\tdepends_on={}", step.id, step.phase, step.title, step.depends_on.join(","));
        }
    }
}

fn print_generation_text(output: &ScaffoldOutput) {
    println!("output_dir:\t{}", output.output_dir.display());
    for artifact in &output.artifacts {
        println!("{:?}\t{}\t{}", artifact.kind, artifact.path.display(), artifact.description);
    }
}

#[cfg(test)]
mod tests {
    use super::{CliCommand, CommandKind, parse_args};

    #[test]
    fn parses_json_flag_before_path() {
        let command = parse_args(["summary".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(command, CliCommand { kind: CommandKind::Summary, path: "specs".to_string(), output_dir: None, json: true });
    }

    #[test]
    fn parses_json_flag_after_path() {
        let command = parse_args(["inspect".to_string(), "specs".to_string(), "--json".to_string()]).expect("args should parse");
        assert_eq!(command, CliCommand { kind: CommandKind::Inspect, path: "specs".to_string(), output_dir: None, json: true });
    }

    #[test]
    fn rejects_extra_positional_arguments() {
        let error = parse_args(["summary".to_string(), "specs".to_string(), "extra".to_string()]).expect_err("extra args should fail");
        assert!(error.contains("usage: fastspec"));
    }

    #[test]
    fn parses_validate_command() {
        let command = parse_args(["validate".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(command, CliCommand { kind: CommandKind::Validate, path: "specs".to_string(), output_dir: None, json: true });
    }

    #[test]
    fn parses_graph_command() {
        let command = parse_args(["graph".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(command, CliCommand { kind: CommandKind::Graph, path: "specs".to_string(), output_dir: None, json: true });
    }

    #[test]
    fn parses_plan_command() {
        let command = parse_args(["plan".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(command, CliCommand { kind: CommandKind::Plan, path: "specs".to_string(), output_dir: None, json: true });
    }

    #[test]
    fn parses_generate_command_with_out_dir() {
        let command =
            parse_args(["generate".to_string(), "--json".to_string(), "--out".to_string(), "out".to_string(), "specs".to_string()])
                .expect("args should parse");
        assert_eq!(
            command,
            CliCommand { kind: CommandKind::Generate, path: "specs".to_string(), output_dir: Some("out".to_string()), json: true }
        );
    }

    #[test]
    fn rejects_generate_without_out_dir() {
        let error = parse_args(["generate".to_string(), "specs".to_string()]).expect_err("generate should require out dir");
        assert!(error.contains("usage: fastspec"));
    }

    #[test]
    fn rejects_out_dir_for_non_generate_command() {
        let error = parse_args(["plan".to_string(), "--out".to_string(), "out".to_string(), "specs".to_string()])
            .expect_err("out dir should not be accepted");
        assert!(error.contains("usage: fastspec"));
    }
}
