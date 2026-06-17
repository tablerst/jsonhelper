# Suggested Commands

- `cargo test --workspace` - run all Rust tests.
- `cargo fmt --all` - format Rust code.
- `cargo fmt --all -- --check` - verify Rust formatting.
- `cargo clippy --workspace --all-targets` - run Rust lints.
- `cargo run -p parselens-cli -- fmt --syntax json5 <file>` - format JSON5 input.
- `cargo run -p parselens-cli -- repr-json <file>` - convert repr-like diagnostics.
- `cd vscode; npm test` - compile the extension and run focused TypeScript tests.
- `cd vscode; npm run package` - build WASM, compile TypeScript, and produce `parselens-vscode.vsix`.
- PowerShell repo inspection examples: `rg --files`, `git status --short --untracked-files=all`, `Get-Content -Raw <path>`.
