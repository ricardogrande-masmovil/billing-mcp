---
applyTo: '**/*_test.go'
---

# Testing Instructions for Billing MCP Server

## 1. General Testing Principles
- Write unit tests for new functionality, especially for:
    - Domain logic (`internal/<domain>/domain/`).
    - Critical infrastructure components (e.g., repository methods).
- Place tests in `_test.go` files within the same package as the code being tested.
- Aim for clear, focused tests that verify specific behaviors.

## 2. Tools & Libraries
- Use `github.com/stretchr/testify/assert` for non-fatal assertions.
- Use `github.com/stretchr/testify/require` for assertions that should halt the test on failure.

## 3. What to Test
- **Domain Services:** Business logic, edge cases, validation.
- **Repositories:** Mock DB or use test DB to verify queries & mapping.
- **Converters:** Accurate mapping between domain and infra/API models.
- **MCP Handlers/Ports:** Request parsing, service calls, response formatting. Mock domain services.

## 4. Best Practices
- **Isolation:** Keep unit tests isolated. Use mocks/stubs.
- **Readability:** Descriptive test names.
- **Coverage:** Good coverage of critical paths.
- **Maintenance:** Update tests with code changes.

## 5. When Assisting with Tests
- Suggest/generate tests for new/modified functions.
- Help identify test cases (positive, negative, edge).
- Assist in creating mocks.
