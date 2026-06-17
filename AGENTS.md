# Repository Guidelines

## Project Structure & Module Organization

This repository is a Rust workspace for ParseLens. New development should happen under `crates/`: `parselens-core` contains JSON, JSONC, and JSON5 parsing/formatting APIs; `parselens-repr` converts Python `repr`-like diagnostics into lossy JSON; `parselens-cli` provides the `parselens` binary; `parselens-wasm` exposes the editor integration wrapper. The VS Code extension lives under `vscode/`.

## Build, Test, and Development Commands

- `cargo test --workspace`: run all Rust tests.
- `cargo fmt --all`: format Rust code.
- `cargo clippy --workspace --all-targets`: run Rust lints.
- `cargo run -p parselens-cli -- fmt --syntax json5 <file>`: format JSON5 as pretty JSON.
- `cargo run -p parselens-cli -- repr-json <file>`: convert repr-like diagnostics to JSON.
- `cd vscode && npm test`: compile the TypeScript extension and run focused extension tests.
- `cd vscode && npm run package`: build the WASM wrapper, compile TypeScript, and create `parselens-vscode.vsix`.

## Coding Style & Naming Conventions

Use idiomatic Rust 2024 and keep public APIs small. Prefer explicit error types in libraries and `anyhow` only at CLI boundaries. Keep parser/formatter logic out of the CLI crate so future WASI or editor integrations can depend on `parselens-core` directly. Use snake_case for modules, functions, and variables; use PascalCase for types. Use `ParseLens` for user-facing product text and `parselens` for crate names, binary names, command ids, file names, and URI schemes.

## Testing Guidelines

Use Rust unit tests for crate-local behavior and add fixture-style cases for JSON5/JSONC syntax: comments, trailing commas, single quotes, unquoted keys, hexadecimal numbers, and malformed input. The target is at least 80% line coverage via `cargo llvm-cov --workspace --fail-under-lines 80`. Repr parsing tests should assert diagnostic JSON shape, not Python object equivalence.

## Dev Docs Usage

- `README.md` is the product and development entry point. Update it when setup,
  capabilities, IPC behavior, roadmap, or security expectations change.
- If `dev_docs/` is introduced or used, treat it as the project working-doc
  layer:
  - `dev_docs/exec/`: active design docs, execution plans, acceptance notes, and
    current tracking.
  - `dev_docs/archive/`: historical context only; do not treat archive notes as
    the default source of truth.
- Prefer updating an existing primary guide before creating parallel docs.
- Keep docs organized by module/topic, for example `ui`, `ipc`, `backend`,
  `storage`, `desktop`, or `testing`.
- Chinese prose is acceptable and preferred for local design/tracking docs when
  it improves clarity for the maintainers. Keep API names, command names, file
  paths, type names, error codes, log keys, and external protocol terms in
  English when that avoids ambiguity.
- When archiving or moving docs, update nearby indexes and internal links in the
  same change.

## Serena Usage

- If Serena tools are available, call `serena.activate_project` once for this
  project before substantive work unless it is already active.
- Read the Serena Instructions Manual once per context before using Serena for
  project work.
- Prefer Serena symbolic tools for code exploration and edits when they fit the
  task. Use normal shell reads for non-code docs and configs.
- Use Serena memories when they are likely to contain relevant repo conventions,
  command guidance, prior architecture decisions, or task completion rules.

## SubAgent Usage

- Treat requests to use or parallelize with SubAgents, including vague wording
  like "结合 SubAgent 并行推进", as permission to split execution work when it is
  safe and useful.
- Before delegating substantial work, identify the critical path and independent
  workstreams. Keep urgent, tightly coupled, or high-risk edits local.
- Delegate bounded implementation, exploration, or verification slices with
  clear ownership and non-overlapping write scopes.
- Prefer implementation workers for well-scoped paths/modules and explorer
  workers when the boundary is unclear or discovery is the main value.
- Prompts to SubAgents should state owned paths, expected output, relevant
  assumptions, validation commands, and whether edits are allowed.
- Reuse an existing SubAgent only when the follow-up stays within the same
  bounded context; otherwise spawn a new narrowly scoped one.

## Commit & Pull Request Guidelines

Recent history uses short Conventional Commit-style subjects such as `feat: Create top-level functions`. Continue with `feat:`, `fix:`, `docs:`, or similar prefixes. Pull requests should describe behavior changes, list validation commands, and call out compatibility impacts for CLI, library API, VS Code extension users, and future WASI consumers.
