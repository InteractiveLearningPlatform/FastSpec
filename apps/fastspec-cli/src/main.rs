use std::env;
use std::path::Path;
use std::process::ExitCode;

use fastspec_core::{InspectOutput, SummaryOutput, ValidationOutput, parse_spec_path, validate_findings, validate_spec_tree};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum CommandKind {
    Summary,
    Inspect,
    Validate,
}

#[derive(Debug, Clone, PartialEq, Eq)]
struct CliCommand {
    kind: CommandKind,
    path: String,
    json: bool,
}

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
        _ => return Err("usage: fastspec <summary|inspect|validate> [--json] <path>".to_string()),
    };

    let mut json = false;
    let mut path = None;
    for arg in args {
        if arg == "--json" {
            json = true;
        } else if path.is_none() {
            path = Some(arg);
        } else {
            return Err("usage: fastspec <summary|inspect> [--json] <path>".to_string());
        }
    }

    let Some(path) = path else {
        return Err("usage: fastspec <summary|inspect|validate> [--json] <path>".to_string());
    };

    Ok(CliCommand { kind, path, json })
}

fn run_command(command: CliCommand) -> ExitCode {
    match command.kind {
        CommandKind::Summary => print_summary(Path::new(&command.path), command.json),
        CommandKind::Inspect => inspect_path(Path::new(&command.path), command.json),
        CommandKind::Validate => validate_path(Path::new(&command.path), command.json),
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

#[cfg(test)]
mod tests {
    use super::{CliCommand, CommandKind, parse_args};

    #[test]
    fn parses_json_flag_before_path() {
        let command = parse_args(["summary".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(command, CliCommand { kind: CommandKind::Summary, path: "specs".to_string(), json: true });
    }

    #[test]
    fn parses_json_flag_after_path() {
        let command = parse_args(["inspect".to_string(), "specs".to_string(), "--json".to_string()]).expect("args should parse");
        assert_eq!(command, CliCommand { kind: CommandKind::Inspect, path: "specs".to_string(), json: true });
    }

    #[test]
    fn rejects_extra_positional_arguments() {
        let error = parse_args(["summary".to_string(), "specs".to_string(), "extra".to_string()]).expect_err("extra args should fail");
        assert!(error.contains("usage: fastspec"));
    }

    #[test]
    fn parses_validate_command() {
        let command = parse_args(["validate".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(command, CliCommand { kind: CommandKind::Validate, path: "specs".to_string(), json: true });
    }
}
