# json5helper

`json5helper` is being migrated from an early Go prototype to a Rust workspace for parsing and formatting JSON-family configuration formats.

The Rust implementation currently targets:

- JSON, JSONC, and JSON5 parsing into `serde_json::Value`
- canonical or pretty JSON output
- a CLI binary named `json5helper`
- a VS Code right-click preview extension
- a separate diagnostic parser that converts Python `repr`-like object graphs into lossy JSON
- a WebAssembly wrapper for editor integration

The previous Go implementation remains in the repository as historical reference under `internal/`, `pkg/`, `test/`, and `jsonhelper.go`.

## Workspace

```text
crates/
  json5helper-core/   # JSON, JSONC, JSON5 parse and format API
  json5helper-cli/    # CLI binary
  json5helper-wasm/   # WebAssembly wrapper for editor integration
  repr-json/          # Python repr-like text to diagnostic JSON
vscode/               # VS Code extension package
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

## VS Code Extension

The `vscode/` package provides an explicit preview workflow instead of a
formatter. It does not register Format Document, does not run on save, and never
mutates the source document.

Use it from the editor context menu:

1. Select text, or leave the cursor in a document to use the whole document.
2. Right-click and choose `Json5helper: Parse Preview`.
3. Pick `Parse as JSON`, `Parse as JSONC`, `Parse as JSON5`, or
   `Parse as Python repr`.
4. The result opens beside the current editor as an editable untitled JSON
   document.

JSONC and JSON5 previews are intentionally lossy: comments, single quotes,
unquoted keys, trailing commas, and other source-level syntax are converted into
canonical JSON.

Extension development commands:

```bash
cd vscode
npm install
npm run build
```

From the repository root, use the VS Code launch configuration
`Run Json5helper VS Code Extension` to build the extension and open an Extension
Development Host.

`npm run build` compiles the Rust WebAssembly wrapper and the TypeScript
extension. Install these one-time prerequisites when needed:

```bash
rustup target add wasm32-unknown-unknown
cargo install wasm-bindgen-cli --version 0.2.125 --locked
```

## Coverage Goal

The target for the Rust implementation is at least 80% line coverage. Use `cargo llvm-cov --workspace --fail-under-lines 80` once `cargo-llvm-cov` is installed.

## WebAssembly

The core crate keeps parsing logic independent from filesystem concerns. The
current VS Code extension uses a bundled `wasm-bindgen` backend so users do not
need Rust, the `json5helper` CLI, or `wasmtime` installed.
