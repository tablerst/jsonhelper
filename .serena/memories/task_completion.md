# Task Completion

- Run `cargo fmt --all -- --check` after Rust edits.
- Run `cargo test --workspace` for normal code changes.
- Run `cargo clippy --workspace --all-targets` when behavior or public API changes.
- Run `cd vscode; npm test` after TypeScript extension edits.
- Run `cd vscode; npm run package` when WASM, extension metadata, or VSIX packaging changes.
- For parser changes, add focused JSON5/JSONC/repr fixture or unit coverage and assert diagnostic JSON shape.
- Check `git status --short --untracked-files=all` before final response; do not touch unrelated user changes.
