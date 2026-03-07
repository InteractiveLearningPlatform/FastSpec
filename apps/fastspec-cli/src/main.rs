use std::env;
use std::path::Path;
use std::process::ExitCode;

use fastspec_core::{parse_spec_path, validate_spec_tree};

fn main() -> ExitCode {
    let mut args = env::args().skip(1);
    match (args.next().as_deref(), args.next(), args.next()) {
        (Some("inspect"), Some(path), None) => inspect_path(Path::new(&path)),
        (Some("summary"), Some(path), None) => print_summary(Path::new(&path)),
        _ => {
            eprintln!("usage: fastspec <summary|inspect> <path>");
            ExitCode::from(2)
        }
    }
}

fn print_summary(path: &Path) -> ExitCode {
    match validate_spec_tree(path) {
        Ok(summaries) => {
            for summary in summaries {
                println!("{}\t{}\t{}\t{}", summary.kind.as_str(), summary.id, summary.title, summary.path.display());
            }
            ExitCode::SUCCESS
        }
        Err(error) => {
            eprintln!("{error}");
            ExitCode::from(1)
        }
    }
}

fn inspect_path(path: &Path) -> ExitCode {
    match parse_spec_path(path) {
        Ok(documents) => {
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
        Err(error) => {
            eprintln!("{error}");
            ExitCode::from(1)
        }
    }
}
