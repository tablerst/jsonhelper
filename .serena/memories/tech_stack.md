# Tech Stack

- Language: Go.
- Module path: `github.com/tablerst/jsonhelper`.
- `go.mod` currently pins `go 1.23.1`.
- No third-party module dependencies currently; `go.sum` is empty.
- Standard Go tooling only; no Makefile, task runner, linter config, or CI config observed.
- Serena semantic tools require `gopls`; current environment lacked `gopls` during onboarding, so direct shell inspection was used.