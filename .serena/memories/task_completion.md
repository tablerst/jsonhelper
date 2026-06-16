# Task Completion

- Run `gofmt -w` on any changed Go files.
- Run `go test ./...` for normal code changes.
- Run `go vet ./...` when behavior or package structure changes.
- For parser/lexer/encoder changes, consider `go test ./... -cover` and add focused JSON5/JSONC edge-case tests.
- Check `git status --short --untracked-files=all` before final response; do not touch unrelated user changes.