# ParseLens

ParseLens turns structured and semi-structured text into readable JSON previews.

The current implementation covers:

- JSON, JSONC, and JSON5 parsing into `serde_json::Value`
- canonical or pretty JSON output
- a CLI binary named `parselens`
- a VS Code right-click preview extension
- a diagnostic parser that converts Python `repr`-like object graphs into lossy JSON
- a WebAssembly wrapper for editor integration

## Workspace

```text
crates/
  parselens-core/   # JSON, JSONC, JSON5 parse and format API
  parselens-cli/    # CLI binary
  parselens-wasm/   # WebAssembly wrapper for editor integration
  parselens-repr/   # Python repr-like text to diagnostic JSON
vscode/             # VS Code extension package
```

## Commands

```bash
cargo test --workspace
cargo fmt --all
cargo clippy --workspace --all-targets
```

Parse and format JSON5:

```bash
parselens fmt --syntax json5 example.json5
```

Read from stdin:

```bash
echo "{unquoted: 'value', trailing: [1, 2,]}" | parselens fmt --syntax json5
```

Convert Python `repr`-like diagnostics:

```bash
echo "AgentExecutor(verbose=True)" | parselens repr-json
```

When running from source, use the CLI crate:

```bash
cargo run -p parselens-cli -- fmt --syntax json5 example.json5
cargo run -p parselens-cli -- repr-json agent.repr
```

## VS Code Extension

The `vscode/` package provides an explicit preview workflow instead of a
formatter. It does not register Format Document, does not run on save, and never
mutates the source document.

Use it from the editor context menu:

1. Select text, or leave the cursor in a document to use the whole document.
2. Right-click and choose `ParseLens: Parse Preview`.
3. Pick `Preview as JSON`, `Preview as JSONC`, `Preview as JSON5`, or
   `Preview as Python repr` from the submenu.
4. The result opens beside the current editor in a read-only JSON preview
   editor with normal VS Code folding and syntax highlighting.

The command is also available from the Command Palette as mode-specific
`ParseLens` preview commands, which is useful when checking whether the
extension activated correctly.

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
`Run ParseLens VS Code Extension` to build the extension and open an Extension
Development Host. This default launch keeps the normal extension-host
environment so compatibility issues remain visible. If unrelated installed
extensions pollute the logs or make debugging noisy, use
`Run ParseLens VS Code Extension (Isolated)` instead; that profile disables
other extensions and stores its temporary VS Code state under `.vscode-dev/`.

`npm run build` compiles the Rust WebAssembly wrapper and the TypeScript
extension. Install these one-time prerequisites when needed:

```bash
rustup target add wasm32-unknown-unknown
cargo install wasm-bindgen-cli --version 0.2.125 --locked
```

To package the extension for another PC:

```bash
cd vscode
npm run package
```

This creates `vscode/parselens-vscode.vsix`. Install it with VS Code's
`Extensions: Install from VSIX...` command, or:

```bash
code --install-extension parselens-vscode.vsix --force
```

For VS Code Insiders, use:

```bash
code-insiders --install-extension parselens-vscode.vsix --force
```

If the command only shows "Activating Extensions..." and no `ParseLens` output
channel appears, close the old Extension Development Host or reload the target
VS Code window. VS Code can keep running an older extension host after the
extension is rebuilt.

## Coverage Goal

The target for the Rust implementation is at least 80% line coverage. Use
`cargo llvm-cov --workspace --fail-under-lines 80` once `cargo-llvm-cov` is
installed.

## WebAssembly

The core crate keeps parsing logic independent from filesystem concerns. The
current VS Code extension uses a bundled `wasm-bindgen` backend so users do not
need Rust, the `parselens` CLI, or `wasmtime` installed.
