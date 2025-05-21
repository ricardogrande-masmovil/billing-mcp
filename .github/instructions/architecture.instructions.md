---
applyTo: '**'
---

# Project Architecture: Hexagonal (Ports & Adapters) with DDD

This document outlines the architectural principles for the `billing-mcp` project, designed to guide GitHub Copilot in understanding and assisting with code generation and modifications. The architecture is based on the Hexagonal (Ports and Adapters) pattern combined with concepts from Domain-Driven Design (DDD).

## Core Principles

1.  **Domain-Centric:** The business logic (domain) is the core of the application. It should be pure, with no dependencies on infrastructure, frameworks, or specific delivery mechanisms.
2.  **Dependency Inversion:** Dependencies flow inwards. Outer layers (infrastructure, API) depend on abstractions (interfaces/ports) defined by or for the inner layers (application, domain).
3.  **Separation of Concerns:** Clear separation between business logic, application logic, and infrastructure concerns.

## Layers & Directory Structure

The project is structured into the following main layers, reflected in the directory structure:

1.  **Domain Layer (`internal/<module>/domain/`)**
    *   **Purpose:** Contains the heart of the business logic. This is the innermost part of the hexagon.
    *   **Contents:**
        *   **Models (`model/`)**: Domain entities, value objects, and aggregates (e.g., `Movement`, `Invoice`). These represent business concepts and rules. UUIDs are typically used for primary keys of entities.
        *   **Services (`service.go`)**:
            *   **Domain Services:** Encapsulate core domain logic that doesn't naturally fit within an entity or value object (e.g., complex calculations or rules involving multiple entities).
            *   **Repository Interfaces (Output Ports):** Defines contracts for data persistence, abstracting the actual storage mechanism (e.g., `MovementRepository` interface).
    *   **Key Characteristics:**
        *   Pure Go code.
        *   No dependencies on specific frameworks (like GORM, Echo) or outer layers.
        *   Defines interfaces for its needs (output ports like repositories).

2.  **Application Layer (`internal/<module>/app/`)**
    *   **Purpose:** Orchestrates domain objects and calls Domain Services to fulfill specific application use cases. These services act as the primary entry point for driving adapters (input ports).
    *   **Contents:**
        *   Application Service implementations (e.g., `CreateMovementUseCase` in `internal/movements/app/usecase.go` or `InvoiceApplicationService` in `internal/invoices/app/service.go`).
        *   These services implement the application's use cases, coordinating domain entities and domain services.
        *   May define their own input/output models (DTOs) for use cases if they differ from domain models, often placed in a `dto/` or `model/` subdirectory within the `app/` layer or within the specific use case file.
    *   **Dependencies:** Depends on Domain Services and Repository Interfaces (defined in the domain layer). Does not depend on outer layers like infrastructure or specific delivery mechanisms.

3.  **Adapter Layers (Ports & Adapters)**

    *   **A. Input/Driving Adapters (`internal/<module>/ports/` and `api/`)**
        *   **Purpose:** Connect external actors (users, other systems, HTTP requests, MCP calls) to the application services (input ports).
        *   **Directory Structure:**
            *   `internal/<module>/ports/<protocol>/` (e.g., `internal/invoices/ports/mcp.go`): Handles protocol-specific request/response mapping and calls application services. Contains port-specific models/DTOs and converters.
            *   `api/mcp/mcp.go`, `api/mcp/tools.go`: Defines MCP tools and may contain central MCP request handling logic that delegates to specific domain port handlers.
        *   **Responsibilities:**
            *   Parse incoming requests.
            *   Validate input (can be basic validation, complex validation might be in domain).
            *   Call appropriate application service methods.
            *   Format responses.
        *   **Dependencies:** Depend on application service interfaces.

    *   **B. Output/Driven Adapters (`internal/<module>/infrastructure/`)**
        *   **Purpose:** Implementations of output ports defined by the domain/application layer (e.g., repository interfaces, external service clients).
        *   **Directory Structure:**
            *   `internal/<module>/infrastructure/persistence/`: Contains data persistence implementations.
                *   `repository.go`: Implements the domain's repository interface (e.g., `MovementSQLRepository` implementing `MovementRepository`).
                *   `sql/model.go`: Infrastructure-specific data models (e.g., GORM structs). These are distinct from domain models.
                *   `sql/converter.go`: Maps between domain models and infrastructure (e.g., GORM) models.
                *   `sql/sql_client.go`: (If used) Low-level database interaction logic.
            *   `internal/<module>/infrastructure/<external_service>/`: Clients for other external services.
        *   **Responsibilities:** Interact with external systems like databases, message queues, third-party APIs.
        *   **Dependencies:** Implement interfaces defined in the domain layer. May use specific libraries (e.g., GORM, HTTP clients).

4.  **Configuration (`config/`)**
    *   `config.go`: Defines the application configuration structure.
    *   Values are accessed via `config.Config`.

5.  **Dependency Injection (`cmd/di/`)**
    *   `wire.go`: Uses `google/wire` to define providers and assemble the application graph.
    *   `wire_gen.go`: Generated by Wire.
    *   Initializes and injects dependencies across layers.

6.  **Main (`cmd/main.go`)**
    *   Application entry point.
    *   Initializes configuration, logging, DI container, and starts the server/application.

## Key Guidelines for Copilot

*   **Identify the Layer:** When asked to add or modify code, first determine which architectural layer it belongs to.
*   **Domain Purity:** Keep the `internal/<module>/domain/` layer free of infrastructure or framework-specific code. Business logic should be expressed in plain Go.
*   **Interfaces are Key (Ports):**
    *   For new external interactions (database, API calls), define an interface (port) in the domain or application layer.
    *   Implement this interface with an adapter in the `infrastructure` layer.
*   **Data Models:**
    *   Domain models (`internal/<module>/domain/model/`) represent core business concepts.
    *   Infrastructure models (`internal/<module>/infrastructure/persistence/sql/model.go`) are specific to the persistence technology (e.g., GORM).
    *   Port models (`internal/<module>/ports/<protocol>/model.go`) are specific to the API/transport layer (e.g., request/response DTOs).
    *   Use converters (`converter.go`) to map between these models at the boundaries of layers.
*   **Error Handling:** Use `fmt.Errorf` with `%w` for wrapping errors, especially when crossing layer boundaries. Return errors upwards.
*   **Logging:** Use the `zerolog.Logger` provided via DI, adding contextual fields.
*   **MCP Implementation:**
    *   Tools are defined in `api/mcp/tools.go`.
    *   Handlers are typically in `api/mcp/mcp.go` (central) and/or `internal/<domain>/ports/mcp.go` (domain-specific MCP logic).
    *   MCP handlers should adapt incoming requests and call the relevant domain/application services.

By following these guidelines, we aim to maintain a clean, decoupled, and testable codebase.
