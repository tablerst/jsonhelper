# Conventions

- Use `gofmt`; do not hand-align or use spaces for Go indentation.
- Keep root API small: root package delegates through `pkg/jsonutil`; implementation belongs under `internal/*`.
- Package names are role-oriented: `lexer`, `parser`, `encoder`, `jsonutil`.
- Preserve public API stability for exported `Parse` and `Encode` helpers unless explicitly changing API.
- Test files use Go standard `testing`, `*_test.go`, and `TestXxx` names.
- Recent git history uses short Conventional Commit-style subjects: `feat: ...`, `fix: ...`.