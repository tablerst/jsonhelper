# Repository Guidelines

## Project Structure & Module Organization

This repository is migrating from an early Go prototype to a Rust workspace. New development should happen under `crates/`: `json5helper-core` contains JSON, JSONC, and JSON5 parsing/formatting APIs; `repr-json` converts Python `repr`-like diagnostics into lossy JSON; `json5helper-cli` provides the `json5helper` binary. The old Go implementation remains as reference in `jsonhelper.go`, `pkg/`, `internal/`, and `test/`.

## Build, Test, and Development Commands

- `cargo test --workspace`: run all Rust tests.
- `cargo fmt --all`: format Rust code.
- `cargo clippy --workspace --all-targets`: run Rust lints.
- `cargo run -p json5helper-cli -- fmt --syntax json5 <file>`: format JSON5 as pretty JSON.
- `cargo run -p json5helper-cli -- repr-json <file>`: convert repr-like diagnostics to JSON.
- `go test ./...`: run the legacy Go tests when touching historical Go code.

## Coding Style & Naming Conventions

Use idiomatic Rust 2024 and keep public APIs small. Prefer explicit error types in libraries and `anyhow` only at CLI boundaries. Keep parser/formatter logic out of the CLI crate so future WASI or editor integrations can depend on `json5helper-core` directly. Use snake_case for modules, functions, and variables; use PascalCase for types.

## Testing Guidelines

Use Rust unit tests for crate-local behavior and add fixture-style cases for JSON5/JSONC syntax: comments, trailing commas, single quotes, unquoted keys, hexadecimal numbers, and malformed input. The target is at least 80% line coverage via `cargo llvm-cov --workspace --fail-under-lines 80`. Repr parsing tests should assert diagnostic JSON shape, not Python object equivalence.

## Commit & Pull Request Guidelines

Recent history uses short Conventional Commit-style subjects such as `feat: Create top-level functions`. Continue with `feat:`, `fix:`, `docs:`, or similar prefixes. Pull requests should describe behavior changes, list validation commands, and call out compatibility impacts for CLI, library API, and future WASI consumers.
