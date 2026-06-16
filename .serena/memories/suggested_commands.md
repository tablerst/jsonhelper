# Suggested Commands

- `go test ./...` - run all tests; passed during onboarding.
- `go test ./... -cover` - run tests with coverage; current top-level packages show 0.0% statement coverage because tests are in separate `test` package.
- `go vet ./...` - static analysis; passed during onboarding.
- `gofmt -w <files>` - format changed Go files.
- `go list ./...` - list/validate packages in module.
- PowerShell repo inspection examples: `rg --files`, `git status --short --untracked-files=all`, `Get-Content -Raw <path>`.