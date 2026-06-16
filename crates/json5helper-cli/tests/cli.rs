use std::io::Write;
use std::path::{Path, PathBuf};
use std::process::{Command, Stdio};

use serde_json::Value;

fn bin() -> PathBuf {
    PathBuf::from(env!("CARGO_BIN_EXE_json5helper"))
}

fn fixture(name: &str) -> PathBuf {
    Path::new(env!("CARGO_MANIFEST_DIR"))
        .join("tests")
        .join("fixtures")
        .join(name)
}

fn run(args: &[&str]) -> (bool, String, String) {
    let output = Command::new(bin())
        .args(args)
        .output()
        .expect("failed to run json5helper");

    (
        output.status.success(),
        String::from_utf8(output.stdout).expect("stdout was not utf-8"),
        String::from_utf8(output.stderr).expect("stderr was not utf-8"),
    )
}

fn run_with_stdin(args: &[&str], stdin: &str) -> (bool, String, String) {
    let mut child = Command::new(bin())
        .args(args)
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .expect("failed to spawn json5helper");

    child
        .stdin
        .as_mut()
        .expect("stdin was not piped")
        .write_all(stdin.as_bytes())
        .expect("failed to write stdin");

    let output = child.wait_with_output().expect("failed to wait for child");
    (
        output.status.success(),
        String::from_utf8(output.stdout).expect("stdout was not utf-8"),
        String::from_utf8(output.stderr).expect("stderr was not utf-8"),
    )
}

#[test]
fn formats_json5_file_as_compact_json() {
    let path = fixture("sample.json5");
    let (success, stdout, stderr) = run(&[
        "fmt",
        "--syntax",
        "json5",
        "--compact",
        path.to_str().unwrap(),
    ]);

    assert!(success, "{stderr}");
    let value: Value = serde_json::from_str(&stdout).unwrap();
    assert_eq!(value["metadata"]["hexadecimal"], 912559);
    assert_eq!(value["metadata"]["ratio"], 0.75);
    assert_eq!(value["metadata"]["max"], 8675309.0);
    assert_eq!(value["features"]["trailingComma"][1], "arrays");
    assert_eq!(value["features"]["unicode"], "雪");
    assert_eq!(value["nested"]["array"][1]["label"], "second");
}

#[test]
fn formats_jsonc_file_as_pretty_json() {
    let path = fixture("sample.jsonc");
    let (success, stdout, stderr) = run(&["fmt", "--syntax", "jsonc", path.to_str().unwrap()]);

    assert!(success, "{stderr}");
    let value: Value = serde_json::from_str(&stdout).unwrap();
    assert_eq!(value["name"], "jsonc");
    assert_eq!(value["settings"]["threshold"], 3.5);
    assert_eq!(
        value["settings"]["paths"][1],
        "crates/json5helper-core/src/lib.rs"
    );
    assert_eq!(value["items"][0]["tags"][1], "jsonc");
}

#[test]
fn parses_json_from_stdin_as_compact_json() {
    let (success, stdout, stderr) = run_with_stdin(
        &["parse", "--syntax", "json", "--compact"],
        r#"{"ok":true}"#,
    );

    assert!(success, "{stderr}");
    assert_eq!(stdout.trim_end(), r#"{"ok":true}"#);
}

#[test]
fn converts_repr_file_as_compact_json() {
    let path = fixture("agent.repr");
    let (success, stdout, stderr) = run(&["repr-json", "--compact", path.to_str().unwrap()]);

    assert!(success, "{stderr}");
    let value: Value = serde_json::from_str(&stdout).unwrap();
    assert_eq!(value["$type"], "AgentExecutor");
    assert_eq!(value["verbose"], true);
    assert_eq!(value["max_iterations"], 3);
    assert_eq!(value["agent"]["runnable"]["$op"], "pipe");
    assert_eq!(
        value["agent"]["runnable"]["items"][0]["mapper"]["agent_scratchpad"]["$type"],
        "RunnableLambda"
    );
    assert_eq!(
        value["agent"]["runnable"]["items"][0]["mapper"]["agent_scratchpad"]["$args"][0]["$repr"],
        "lambda x: message_formatter(x['intermediate_steps'])"
    );
    assert_eq!(
        value["agent"]["runnable"]["items"][1]["messages"][1]["variable_name"],
        "chat_history"
    );
    assert_eq!(value["callbacks"][0]["tag"], "debug");
    assert_eq!(
        value["metadata"]["function"],
        "<function _get_type at 0x0000025FFE70A0E0>"
    );
}

#[test]
fn reports_no_input_when_stdin_is_empty() {
    let (success, stdout, stderr) = run_with_stdin(&["fmt", "--syntax", "json5"], "");

    assert!(!success);
    assert!(stdout.is_empty());
    assert!(stderr.contains("no input provided"));
}
