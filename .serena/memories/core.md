# Core

- Go module `github.com/tablerst/jsonhelper`; root package in `jsonhelper.go` re-exports `Parse`/`Encode` via `pkg/jsonutil`.
- Public wrapper package: `pkg/jsonutil/`.
- Implementation-only packages under `internal/`: `lexer`, `parser`, `encoder`, `utils`.
- Current tests live in top-level `test/` package; no per-package test files yet in `internal/*` or `pkg/*`.
- Docs: `README.md`, `README.zh.md`; current README text may display mojibake under default PowerShell reads, verify encoding before editing.
- Read for toolchain/commands: `mem:tech_stack`, `mem:suggested_commands`.
- Read for style/task closeout: `mem:conventions`, `mem:task_completion`.