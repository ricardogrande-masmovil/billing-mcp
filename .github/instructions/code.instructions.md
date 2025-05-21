---
applyTo: '**'
---

# Go Coding & Best Practices

## General Conventions
- Follow Go idioms (Effective Go, CodeReviewComments).
- Prioritize clarity, readability, conciseness.
- **Error Handling:** Use `fmt.Errorf` with `%w`; return errors.
- **Logging:** Use `zerolog.Logger` with context fields.
- **Dependencies:** Go Modules (`go.mod`, `go.sum`).
- **UUIDs:** Use for primary keys.

## GORM & Database
- GORM models in `infrastructure/.../sql/model.go`.
- Converters in `.../sql/converter.go`.
- Access via repository implementations.

## MCP Implementation
- Tools: `api/mcp/tools.go`.
- Handlers: `api/mcp/mcp.go`, `internal/<domain>/ports/mcp.go`.
- Map tools to domain services.

## When Assisting
- Identify layer: domain, infra, port, api.
- Keep domain logic pure; infra adapters implement interfaces.
- Use `config.Config` for config values.

## File Guidance
- `cmd/main.go`: startup, shutdown.
- `cmd/di/wire.go`: DI wiring.
- `config/config.go`: config structure.
- `internal/*/domain/service.go`: business logic.
- `internal/*/infra/persistence/repository.go`: data access.
- `api/mcp/*`: MCP definitions.
