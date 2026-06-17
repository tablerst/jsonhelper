# Core

- ParseLens is a Rust workspace for previewing JSON-family and repr-like diagnostic text as readable JSON.
- Core crates:
  - `crates/parselens-core/`: JSON, JSONC, and JSON5 parse/format API.
  - `crates/parselens-repr/`: lossy Python `repr`-like text to diagnostic JSON.
  - `crates/parselens-cli/`: `parselens` CLI binary.
  - `crates/parselens-wasm/`: `wasm-bindgen` wrapper used by the editor integration.
- VS Code extension package: `vscode/`.
- Docs: `README.md`, `README.zh.md`, `vscode/README.md`, `AGENTS.md`.
- Read for toolchain/commands: `mem:tech_stack`, `mem:suggested_commands`.
- Read for style/task closeout: `mem:conventions`, `mem:task_completion`.
