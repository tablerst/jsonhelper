use std::fs;
use std::io::{self, Read};
use std::path::PathBuf;

use anyhow::{Context, Result, bail};
use clap::{Parser, Subcommand, ValueEnum};
use json5helper_core::Syntax;

#[derive(Debug, Parser)]
#[command(name = "json5helper")]
#[command(about = "Parse and format JSON, JSONC, JSON5, and Python repr-like diagnostics.")]
struct Cli {
    #[command(subcommand)]
    command: Command,
}

#[derive(Debug, Subcommand)]
enum Command {
    /// Parse and print canonical JSON.
    Parse(JsonCommand),
    /// Format input as pretty JSON.
    Fmt(JsonCommand),
    /// Convert Python repr-like object graphs into diagnostic JSON.
    ReprJson(ReprCommand),
}

#[derive(Debug, Parser)]
struct JsonCommand {
    /// Input syntax.
    #[arg(long, value_enum, default_value_t = CliSyntax::Json5)]
    syntax: CliSyntax,
    /// Read input from this file instead of stdin.
    input: Option<PathBuf>,
    /// Emit compact JSON instead of pretty JSON.
    #[arg(long)]
    compact: bool,
}

#[derive(Debug, Parser)]
struct ReprCommand {
    /// Read input from this file instead of stdin.
    input: Option<PathBuf>,
    /// Emit compact JSON instead of pretty JSON.
    #[arg(long)]
    compact: bool,
}

#[derive(Debug, Clone, Copy, ValueEnum)]
enum CliSyntax {
    Json,
    Jsonc,
    Json5,
}

impl From<CliSyntax> for Syntax {
    fn from(value: CliSyntax) -> Self {
        match value {
            CliSyntax::Json => Syntax::Json,
            CliSyntax::Jsonc => Syntax::Jsonc,
            CliSyntax::Json5 => Syntax::Json5,
        }
    }
}

fn main() -> Result<()> {
    let cli = Cli::parse();
    match cli.command {
        Command::Parse(command) | Command::Fmt(command) => run_json(command),
        Command::ReprJson(command) => run_repr(command),
    }
}

fn run_json(command: JsonCommand) -> Result<()> {
    let input = read_input(command.input.as_ref())?;
    let output = format_json_text(&input, command.syntax, command.compact)?;
    println!("{output}");
    Ok(())
}

fn run_repr(command: ReprCommand) -> Result<()> {
    let input = read_input(command.input.as_ref())?;
    let output = format_repr_text(&input, command.compact)?;
    println!("{output}");
    Ok(())
}

fn format_json_text(input: &str, syntax: CliSyntax, compact: bool) -> Result<String> {
    Ok(json5helper_core::format(input, syntax.into(), !compact)?)
}

fn format_repr_text(input: &str, compact: bool) -> Result<String> {
    let value = repr_json::parse_repr(input)?;
    Ok(if compact {
        serde_json::to_string(&value)?
    } else {
        serde_json::to_string_pretty(&value)?
    })
}

fn read_input(path: Option<&PathBuf>) -> Result<String> {
    if let Some(path) = path {
        fs::read_to_string(path).with_context(|| format!("failed to read {}", path.display()))
    } else {
        let mut input = String::new();
        io::stdin()
            .read_to_string(&mut input)
            .context("failed to read stdin")?;
        if input.is_empty() {
            bail!("no input provided; pass a file path or pipe text on stdin");
        }
        Ok(input)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn maps_cli_syntax_to_core_syntax() {
        assert_eq!(Syntax::from(CliSyntax::Json), Syntax::Json);
        assert_eq!(Syntax::from(CliSyntax::Jsonc), Syntax::Jsonc);
        assert_eq!(Syntax::from(CliSyntax::Json5), Syntax::Json5);
    }

    #[test]
    fn formats_json5_compact() {
        let output = format_json_text(
            "{unquoted: 'value', trailing: [1, 2,]}",
            CliSyntax::Json5,
            true,
        )
        .unwrap();

        assert_eq!(output, r#"{"trailing":[1,2],"unquoted":"value"}"#);
    }

    #[test]
    fn formats_jsonc_pretty() {
        let output =
            format_json_text(r#"{"value":1,}// comment"#, CliSyntax::Jsonc, false).unwrap();

        assert!(output.contains("\n  \"value\": 1\n"));
    }

    #[test]
    fn formats_repr_compact() {
        let output = format_repr_text("AgentExecutor(verbose=True)", true).unwrap();

        assert_eq!(output, r#"{"$type":"AgentExecutor","verbose":true}"#);
    }

    #[test]
    fn formats_repr_pretty() {
        let output = format_repr_text("AgentExecutor(verbose=True)", false).unwrap();

        assert!(output.contains("\n  \"verbose\": true\n"));
    }
}
