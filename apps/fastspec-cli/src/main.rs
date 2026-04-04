use std::env;
use std::path::Path;
use std::process::ExitCode;

use fastspec_core::{
    GraphOutput, InitOptions, InitOutput, InspectOutput, PlanOutput, ScaffoldOutput, SummaryOutput, ValidationOutput, export_graph,
    export_plan, generate_scaffold, init_spec_tree, parse_spec_path, validate_findings, validate_spec_tree,
};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum CommandKind {
    Summary,
    Inspect,
    Validate,
    Graph,
    Plan,
    Generate,
    Init,
}

/// Output format for the `graph` command.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
enum GraphFormat {
    #[default]
    Text,
    Json,
    Mermaid,
    Dot,
}

#[derive(Debug, Clone, PartialEq, Eq)]
struct CliCommand {
    kind: CommandKind,
    path: String,
    output_dir: Option<String>,
    json: bool,
    /// Only used when kind == Graph.
    graph_format: GraphFormat,
    /// Only used when kind == Init.
    init_id: Option<String>,
    /// Only used when kind == Init.
    init_title: Option<String>,
    /// Only used when kind == Init.
    init_modules: Vec<String>,
    /// Only used when kind == Init.
    init_capabilities: Vec<String>,
}

const USAGE: &str = "usage: fastspec <summary|inspect|validate|graph|plan|generate|init> [--json] [--out <dir>] [--id <id>] [--title <title>] [--modules <m1,m2>] [--capabilities <c1,c2>] <path>";

const VERSION: &str = env!("CARGO_PKG_VERSION");

const HELP_TEXT: &str = "fastspec — FastSpec document pipeline\n\nUsage:\n  fastspec <command> [options] <path>\n\nCommands:\n  summary   List all spec documents in a directory tree\n  inspect   Show detailed metadata and spec body for each document\n  validate  Run cross-document validation rules and report findings\n  graph     Export dependency graph (--format text|json|mermaid|dot)\n  plan      Export ordered implementation plan\n  generate  Scaffold implementation stubs (--out <dir>)\n  init      Scaffold a new FastSpec project tree\n\nGlobal flags:\n  --json      Output machine-readable JSON instead of text\n  --help, -h  Show this help or per-command help\n  --version   Show version\n\nRun `fastspec <command> --help` for command-specific details.";

const HELP_SUMMARY: &str = "fastspec summary [--json] <path>\n\nList all spec documents found under <path>, showing kind, id, title, and file path.\n\nFlags:\n  --json  Output as JSON array of documents";

const HELP_INSPECT: &str = "fastspec inspect [--json] <path>\n\nShow detailed document content for each spec found under <path>.\n\nFlags:\n  --json  Output as JSON array of documents with full spec bodies";

const HELP_VALIDATE: &str = "fastspec validate [--json] <path>\n\nRun cross-document validation: check for missing/undeclared modules and capabilities,\nduplicate IDs, invalid dependencies, and dependency cycles.\n\nFlags:\n  --json  Output findings as JSON; exits non-zero when findings exist";

const HELP_GRAPH: &str = "fastspec graph [--json | --format <fmt>] <path>\n\nExport the dependency graph for all specs under <path>.\n\nFlags:\n  --format text     Tab-separated text (default)\n  --format json     JSON nodes and edges\n  --format mermaid  Mermaid flowchart LR\n  --format dot      Graphviz DOT digraph\n  --json            Shorthand for --format json";

const HELP_PLAN: &str = "fastspec plan [--json] <path>\n\nExport an ordered implementation plan derived from the dependency graph.\n\nFlags:\n  --json  Output as JSON array of plan steps";

const HELP_GENERATE: &str = "fastspec generate [--json] --out <dir> <path>\n\nScaffold README stubs for modules, workflows, and capabilities.\n\nFlags:\n  --out <dir>  Required. Output directory (must be empty or non-existent).\n  --json       Output artifact list as JSON";

const HELP_INIT: &str = "fastspec init [--json] [--id <id>] [--title <title>] [--modules <ids>] [--capabilities <ids>] <dir>\n\nScaffold a new FastSpec project tree at <dir>.\n\nFlags:\n  --id <id>              Project ID (default: directory name)\n  --title <title>        Project title (default: same as id)\n  --modules <ids>        Comma-separated module IDs to scaffold\n  --capabilities <ids>   Comma-separated agent capability IDs to scaffold\n  --json                 Output artifact list as JSON\n\nExample:\n  fastspec init ./my-project --id my-project --title \"My Project\" --modules api,web";

fn main() -> ExitCode {
    let args: Vec<String> = env::args().skip(1).collect();

    // Top-level version flag.
    if args.first().map(String::as_str) == Some("--version") || args.first().map(String::as_str) == Some("-V") {
        println!("fastspec {VERSION}");
        return ExitCode::SUCCESS;
    }

    // Top-level or per-command help flag.
    let help_requested =
        args.is_empty() || args.first().map(String::as_str) == Some("--help") || args.first().map(String::as_str) == Some("-h");
    let per_command_help = args.iter().skip(1).any(|a| a == "--help" || a == "-h");

    if help_requested {
        println!("{HELP_TEXT}");
        return ExitCode::SUCCESS;
    }

    if per_command_help {
        let help = match args[0].as_str() {
            "summary" => HELP_SUMMARY,
            "inspect" => HELP_INSPECT,
            "validate" => HELP_VALIDATE,
            "graph" => HELP_GRAPH,
            "plan" => HELP_PLAN,
            "generate" => HELP_GENERATE,
            "init" => HELP_INIT,
            _ => HELP_TEXT,
        };
        println!("{help}");
        return ExitCode::SUCCESS;
    }

    match parse_args(args.into_iter()) {
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
        Some("init") => CommandKind::Init,
        _ => return Err(USAGE.to_string()),
    };

    let mut json = false;
    let mut path = None;
    let mut output_dir = None;
    let mut graph_format: Option<GraphFormat> = None;
    let mut init_id: Option<String> = None;
    let mut init_title: Option<String> = None;
    let mut init_modules: Vec<String> = Vec::new();
    let mut init_capabilities: Vec<String> = Vec::new();
    let mut args = args.peekable();
    while let Some(arg) = args.next() {
        if arg == "--json" {
            json = true;
        } else if arg == "--format" {
            let Some(fmt) = args.next() else {
                return Err(USAGE.to_string());
            };
            graph_format = Some(match fmt.as_str() {
                "text" => GraphFormat::Text,
                "json" => GraphFormat::Json,
                "mermaid" => GraphFormat::Mermaid,
                "dot" => GraphFormat::Dot,
                _ => return Err(format!("unknown graph format `{fmt}`; expected text, json, mermaid, or dot")),
            });
        } else if let Some(fmt) = arg.strip_prefix("--format=") {
            graph_format = Some(match fmt {
                "text" => GraphFormat::Text,
                "json" => GraphFormat::Json,
                "mermaid" => GraphFormat::Mermaid,
                "dot" => GraphFormat::Dot,
                _ => return Err(format!("unknown graph format `{fmt}`; expected text, json, mermaid, or dot")),
            });
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
        } else if arg == "--id" {
            let Some(id) = args.next() else {
                return Err(USAGE.to_string());
            };
            init_id = Some(id);
        } else if let Some(id) = arg.strip_prefix("--id=") {
            init_id = Some(id.to_string());
        } else if arg == "--title" {
            let Some(title) = args.next() else {
                return Err(USAGE.to_string());
            };
            init_title = Some(title);
        } else if let Some(title) = arg.strip_prefix("--title=") {
            init_title = Some(title.to_string());
        } else if arg == "--modules" {
            let Some(mods) = args.next() else {
                return Err(USAGE.to_string());
            };
            init_modules = mods.split(',').map(str::trim).filter(|s| !s.is_empty()).map(str::to_string).collect();
        } else if let Some(mods) = arg.strip_prefix("--modules=") {
            init_modules = mods.split(',').map(str::trim).filter(|s| !s.is_empty()).map(str::to_string).collect();
        } else if arg == "--capabilities" {
            let Some(caps) = args.next() else {
                return Err(USAGE.to_string());
            };
            init_capabilities = caps.split(',').map(str::trim).filter(|s| !s.is_empty()).map(str::to_string).collect();
        } else if let Some(caps) = arg.strip_prefix("--capabilities=") {
            init_capabilities = caps.split(',').map(str::trim).filter(|s| !s.is_empty()).map(str::to_string).collect();
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

    if graph_format.is_some() && kind != CommandKind::Graph {
        return Err("--format is only valid for the graph command".to_string());
    }

    let init_only_flags = init_id.is_some() || init_title.is_some() || !init_modules.is_empty() || !init_capabilities.is_empty();
    if init_only_flags && kind != CommandKind::Init {
        return Err("--id, --title, --modules, and --capabilities are only valid for the init command".to_string());
    }

    // Resolve effective graph format: --format overrides --json for the graph command.
    let effective_graph_format = match graph_format {
        Some(fmt) => fmt,
        None if json => GraphFormat::Json,
        None => GraphFormat::Text,
    };

    Ok(CliCommand {
        kind,
        path,
        output_dir,
        json,
        graph_format: effective_graph_format,
        init_id,
        init_title,
        init_modules,
        init_capabilities,
    })
}

fn run_command(command: CliCommand) -> ExitCode {
    match command.kind {
        CommandKind::Summary => print_summary(Path::new(&command.path), command.json),
        CommandKind::Inspect => inspect_path(Path::new(&command.path), command.json),
        CommandKind::Validate => validate_path(Path::new(&command.path), command.json),
        CommandKind::Graph => graph_path(Path::new(&command.path), command.graph_format, command.json),
        CommandKind::Plan => plan_path(Path::new(&command.path), command.json),
        CommandKind::Generate => {
            generate_path(Path::new(&command.path), Path::new(command.output_dir.as_deref().unwrap_or_default()), command.json)
        }
        CommandKind::Init => {
            let opts = InitOptions {
                id: command.init_id.unwrap_or_else(|| {
                    // Default id from directory name, falling back to "project".
                    Path::new(&command.path).file_name().and_then(|n| n.to_str()).unwrap_or("project").to_string()
                }),
                title: command.init_title.unwrap_or_default(),
                modules: command.init_modules,
                capabilities: command.init_capabilities,
            };
            init_path(Path::new(&command.path), opts, command.json)
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

fn graph_path(path: &Path, format: GraphFormat, json: bool) -> ExitCode {
    match export_graph(path) {
        Ok(output) => match format {
            GraphFormat::Json => print_json(&output),
            GraphFormat::Mermaid => {
                println!("{}", render_graph_mermaid(&output));
                ExitCode::SUCCESS
            }
            GraphFormat::Dot => {
                println!("{}", render_graph_dot(&output));
                ExitCode::SUCCESS
            }
            GraphFormat::Text => {
                print_graph_text(&output);
                ExitCode::SUCCESS
            }
        },
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

fn init_path(dir: &Path, opts: InitOptions, json: bool) -> ExitCode {
    match init_spec_tree(dir, opts) {
        Ok(output) => {
            if json {
                return print_json(&output);
            }
            print_init_text(&output);
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

fn render_graph_mermaid(output: &GraphOutput) -> String {
    use fastspec_core::GraphNodeKind;
    let mut lines = vec!["flowchart LR".to_string()];
    for node in &output.nodes {
        // Sanitize IDs: Mermaid node IDs must not contain colons.
        let safe_id = node.id.replace(':', "__");
        match node.kind {
            GraphNodeKind::Project => lines.push(format!("    {safe_id}[\"{}\"]", node.title)),
            GraphNodeKind::Module => lines.push(format!("    {safe_id}(\"{}\")", node.title)),
            GraphNodeKind::Workflow => lines.push(format!("    {safe_id}{{\"{}\"}}", node.title)),
            GraphNodeKind::AgentCapability => lines.push(format!("    {safe_id}[/\"{}\"/]", node.title)),
        }
    }
    for edge in &output.edges {
        // Mermaid uses different arrow styles: --> for hierarchy, -.-> for workflow/capability.
        use fastspec_core::GraphEdgeKind;
        let from = edge.from.replace(':', "__");
        let to = edge.to.replace(':', "__");
        let arrow = match edge.kind {
            GraphEdgeKind::DefinesWorkflow | GraphEdgeKind::DefinesCapability => "-.->",
            GraphEdgeKind::Contains | GraphEdgeKind::DependsOn => "-->",
        };
        lines.push(format!("    {from} {arrow} {to}"));
    }
    lines.join("\n")
}

fn render_graph_dot(output: &GraphOutput) -> String {
    use fastspec_core::GraphNodeKind;
    let mut lines = vec!["digraph fastspec {".to_string()];
    for node in &output.nodes {
        let shape = match node.kind {
            GraphNodeKind::Project => " shape=box",
            GraphNodeKind::Workflow => " shape=diamond",
            GraphNodeKind::AgentCapability => " shape=parallelogram",
            GraphNodeKind::Module => "",
        };
        let label = node.title.replace('"', "\\\"");
        lines.push(format!("    \"{}\" [label=\"{label}\"{shape}];", node.id));
    }
    for edge in &output.edges {
        lines.push(format!("    \"{}\" -> \"{}\";", edge.from, edge.to));
    }
    lines.push("}".to_string());
    lines.join("\n")
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

fn print_init_text(output: &InitOutput) {
    println!("init_dir:\t{}", output.dir.display());
    for artifact in &output.artifacts {
        println!("{:?}\t{}\t{}", artifact.kind, artifact.path.display(), artifact.description);
    }
}

#[cfg(test)]
mod tests {
    use super::{CliCommand, CommandKind, GraphFormat, parse_args};

    #[test]
    fn parses_json_flag_before_path() {
        let command = parse_args(["summary".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(
            command,
            CliCommand {
                kind: CommandKind::Summary,
                path: "specs".to_string(),
                output_dir: None,
                json: true,
                graph_format: GraphFormat::Text,
                init_id: None,
                init_title: None,
                init_modules: vec![],
                init_capabilities: vec![]
            }
        );
    }

    #[test]
    fn parses_json_flag_after_path() {
        let command = parse_args(["inspect".to_string(), "specs".to_string(), "--json".to_string()]).expect("args should parse");
        assert_eq!(
            command,
            CliCommand {
                kind: CommandKind::Inspect,
                path: "specs".to_string(),
                output_dir: None,
                json: true,
                graph_format: GraphFormat::Text,
                init_id: None,
                init_title: None,
                init_modules: vec![],
                init_capabilities: vec![]
            }
        );
    }

    #[test]
    fn rejects_extra_positional_arguments() {
        let error = parse_args(["summary".to_string(), "specs".to_string(), "extra".to_string()]).expect_err("extra args should fail");
        assert!(error.contains("usage: fastspec"));
    }

    #[test]
    fn parses_validate_command() {
        let command = parse_args(["validate".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(
            command,
            CliCommand {
                kind: CommandKind::Validate,
                path: "specs".to_string(),
                output_dir: None,
                json: true,
                graph_format: GraphFormat::Text,
                init_id: None,
                init_title: None,
                init_modules: vec![],
                init_capabilities: vec![]
            }
        );
    }

    #[test]
    fn parses_graph_command() {
        let command = parse_args(["graph".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(
            command,
            CliCommand {
                kind: CommandKind::Graph,
                path: "specs".to_string(),
                output_dir: None,
                json: true,
                graph_format: GraphFormat::Json,
                init_id: None,
                init_title: None,
                init_modules: vec![],
                init_capabilities: vec![]
            }
        );
    }

    #[test]
    fn parses_plan_command() {
        let command = parse_args(["plan".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(
            command,
            CliCommand {
                kind: CommandKind::Plan,
                path: "specs".to_string(),
                output_dir: None,
                json: true,
                graph_format: GraphFormat::Text,
                init_id: None,
                init_title: None,
                init_modules: vec![],
                init_capabilities: vec![]
            }
        );
    }

    #[test]
    fn parses_generate_command_with_out_dir() {
        let command =
            parse_args(["generate".to_string(), "--json".to_string(), "--out".to_string(), "out".to_string(), "specs".to_string()])
                .expect("args should parse");
        assert_eq!(
            command,
            CliCommand {
                kind: CommandKind::Generate,
                path: "specs".to_string(),
                output_dir: Some("out".to_string()),
                json: true,
                graph_format: GraphFormat::Text,
                init_id: None,
                init_title: None,
                init_modules: vec![],
                init_capabilities: vec![]
            }
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
    #[test]
    fn parses_graph_format_flag() {
        let command = parse_args(["graph".to_string(), "--format".to_string(), "mermaid".to_string(), "specs".to_string()])
            .expect("args should parse");
        assert_eq!(command.graph_format, GraphFormat::Mermaid);
        assert!(!command.json);
    }

    #[test]
    fn graph_json_flag_sets_json_format() {
        let command = parse_args(["graph".to_string(), "--json".to_string(), "specs".to_string()]).expect("args should parse");
        assert_eq!(command.graph_format, GraphFormat::Json);
        assert!(command.json);
    }

    #[test]
    fn format_flag_rejected_on_non_graph_commands() {
        let error = parse_args(["plan".to_string(), "--format".to_string(), "mermaid".to_string(), "specs".to_string()])
            .expect_err("--format should be rejected on plan");
        assert!(error.contains("--format is only valid for the graph command"));
    }
}
