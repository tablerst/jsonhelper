# Tech Stack

- Main language: Rust 2024 workspace.
- Crates: `parselens-core`, `parselens-repr`, `parselens-cli`, `parselens-wasm`.
- CLI binary: `parselens`.
- Editor integration: VS Code extension in `vscode/`, written in TypeScript and backed by a bundled `wasm-bindgen` module.
- Repository URL: `https://github.com/tablerst/ParseLens.git`.
- Use Cargo for Rust build/test/lint and npm scripts under `vscode/` for extension build/test/package.
