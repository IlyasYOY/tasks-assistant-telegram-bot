# Repository Conventions

This document describes the conventions and best‑practice guidelines that
should be followed when contributing to this repository.  

Adhering to these conventions helps keep the codebase consistent, readable, and
easy to maintain.

## 1. Go Code Style

| Aspect                     | Convention                                                                 |
|----------------------------|-----------------------------------------------------------------------------|
| **Naming**                 | - Packages: short, lower‑case, no underscores. <br> - Exported identifiers: `CamelCase`. <br> - Unexported identifiers: `camelCase`. |
| **Error handling**         | Wrap errors with context using `fmt.Errorf("…: %w", err)`.                  |
| **Imports**                | Group imports in three blocks: standard library, third‑party, internal.   |
| **Context**                | Pass `context.Context` to functions that perform I/O or long‑running work. |
| **Testing**                | Place tests in `*_test.go` files, use the `testing` package and `testify` for assertions. |

## 3. Commit Messages

Follow the **Conventional Commits** format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

| Type      | Description                              |
|-----------|------------------------------------------|
| `feat`    | New feature                              |
| `fix`     | Bug fix                                  |
| `docs`    | Documentation only changes                |
| `style`   | Formatting, missing semicolons, etc.     |
| `refactor`| Code change that neither fixes a bug nor adds a feature |
| `test`    | Adding or fixing tests                   |
| `chore`   | Build process, auxiliary tools, etc.     |

Example:

```
feat(handler): add unknown command handler

- Respond with a friendly message when an unknown command is received.
- Add unit tests for the new handler.
```

## 4. Testing Guidelines

- **Unit tests** should be fast, deterministic, and not depend on external services.
- Use table‑driven tests where appropriate.
  - Don't parameterize tests with error.
- Keep test files in the same package as the code they test (or use `package xxx_test` for black‑box testing).

## 5. Database Migrations

- Migrations are written in SQL and live under `internal/store/migrations`.
- Use Goose’s `Up`/`Down` blocks (`-- +goose Up` / `-- +goose Down`).
- Migration filenames must start with a timestamp: `YYYYMMDD_hhmmss_description.sql`.
- After adding a migration, run `go test ./...` (the test suite will apply migrations on a temporary in‑memory DB).

## 6. Documentation

- Keep the `README.md` up‑to‑date with installation, configuration, and usage instructions.
- Inline code comments should be clear and concise; avoid redundant explanations.
