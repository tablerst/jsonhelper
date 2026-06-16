# json5helper

`json5helper` is being migrated from an early Go prototype to a Rust workspace for parsing and formatting JSON-family configuration formats.

The Rust implementation currently targets:

- JSON, JSONC, and JSON5 parsing into `serde_json::Value`
- canonical or pretty JSON output
- a CLI binary named `json5helper`
- a separate diagnostic parser that converts Python `repr`-like object graphs into lossy JSON
- future WASI builds for editor integration

The previous Go implementation remains in the repository as historical reference under `internal/`, `pkg/`, `test/`, and `jsonhelper.go`.

## Workspace

```text
crates/
  json5helper-core/   # JSON, JSONC, JSON5 parse and format API
  json5helper-cli/    # CLI binary
  repr-json/          # Python repr-like text to diagnostic JSON
```

## Commands

```bash
cargo test --workspace
cargo fmt --all
cargo clippy --workspace --all-targets
```

Parse and format JSON5:

```bash
json5helper fmt --syntax json5 example.json5
```

Read from stdin:

```bash
echo "{unquoted: 'value', trailing: [1, 2,]}" | json5helper fmt --syntax json5
```

Convert Python `repr`-like diagnostics:

```bash
echo "AgentExecutor(verbose=True)" | json5helper repr-json
```

## Coverage Goal

The target for the Rust implementation is at least 80% line coverage. Use `cargo llvm-cov --workspace --fail-under-lines 80` once `cargo-llvm-cov` is installed.

## WASI

The core crate keeps parsing logic independent from filesystem concerns so it can later be wrapped for `wasm32-wasip2` or consumed by a VS Code extension. The VS Code extension itself is not implemented in this repository yet.
