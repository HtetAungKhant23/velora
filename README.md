# Velora

Velora is an image processing API built with Go, demonstrating **Hexagonal Architecture** and **Tactical Domain-Driven Design (DDD)** patterns. This project serves as a practical reference for building maintainable, testable, and loosely-coupled applications.

---

## Overview

Velora provides user authentication and image upload capabilities. While the business logic is straightforward, the architecture emphasizes:

- **Separation of concerns** between domain logic and infrastructure
- **Dependency inversion** where high-level modules don't depend on low-level modules
- **Encapsulation** of domain invariants within domain objects
- **Testability** through interface-based design

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        ADAPTERS (Infrastructure)                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │   Primary    │  │   Secondary  │  │     Secondary        │  │
│  │   (Inbound)  │  │   (Outbound) │  │     (Outbound)       │  │
│  │              │  │              │  │                      │  │
│  │ HTTP Handler │  │ Repository   │  │ Storage, Token,      │  │
│  │ (chi router) │  │ (PostgreSQL) │  │ Image Processor      │  │
│  └──────┬───────┘  └──────┬───────┘  └──────────┬───────────┘  │
└─────────┼──────────────────┼─────────────────────┼──────────────┘
          │                  │                     │
          │                  │                     │
          ▼                  │                     │
┌────────────────────────────┼─────────────────────┼──────────────┐
│                       PORTS (Interfaces)                        │
│  ┌──────────────┐          │                     │              │
│  │   Primary    │          │                     │              │
│  │  UseCase     │          ▼                     ▼              │
│  │  (inbound)   │  ┌───────────────┐   ┌───────────────────┐   │
│  └──────────────┘  │   Secondary   │   │     Secondary     │   │
│                    │  Repository   │   │  Storage, Token,  │   │
│                    │   (outbound)  │   │   Processor       │   │
│                    └───────────────┘   └───────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        DOMAIN (Core)                            │
│  ┌──────────────────┐  ┌──────────────────┐                    │
│  │   User Entity    │  │   Image Entity   │                    │
│  │                  │  │                  │                    │
│  │  - Email (VO)    │  │  - Format (VO)   │                    │
│  │  - Password (VO) │  │  - Dimensions    │                    │
│  │  - UserID        │  │  - FileSize (VO) │                    │
│  │                  │  │  - StoragePath   │                    │
│  └──────────────────┘  └──────────────────┘                    │
│                                                                 │
│  ┌──────────────────────────────────────────┐                  │
│  │         Application Services             │                  │
│  │   (UserService, ImageService)            │                  │
│  └──────────────────────────────────────────┘                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Project Structure

```
velora/
├── cmd/
│   └── main.go                 # Application entry point & wiring
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   │
│   ├── core/                   # DOMAIN LAYER (no external dependencies)
│   │   ├── domain/
│   │   │   ├── user/           # User aggregate
│   │   │   │   ├── user.go         # Entity
│   │   │   │   ├── email.go        # Value Object
│   │   │   │   └── password.go     # Value Object
│   │   │   ├── image/          # Image aggregate
│   │   │   │   ├── image.go        # Entity
│   │   │   │   ├── format.go       # Value Object
│   │   │   │   ├── dimensions.go   # Value Object
│   │   │   │   ├── file_size.go    # Value Object
│   │   │   │   └── storage_path.go # Value Object
│   │   │   └── shared/
│   │   │       └── error.go        # Domain errors
│   │   │
│   │   ├── ports/              # INTERFACE DEFINITIONS
│   │   │   ├── user_port.go        # Primary port (UseCase)
│   │   │   ├── user_repository.go  # Secondary port
│   │   │   ├── image_port.go       # Primary port
│   │   │   ├── image_repository.go # Secondary port
│   │   │   ├── image_storage_port.go
│   │   │   ├── image_processor.go
│   │   │   └── token_port.go
│   │   │
│   │   └── services/           # APPLICATION SERVICES
│   │       ├── user_service.go
│   │       └── image_service.go
│   │
│   └── adapters/               # ADAPTERS LAYER (infrastructure)
│       ├── handler/            # Primary adapters (HTTP)
│       │   ├── auth_handler.go
│       │   ├── image_handler.go
│       │   ├── health_handler.go
│       │   ├── router.go
│       │   └── middleware/
│       │       ├── auth_guard.go
│       │       └── response.go
│       │
│       ├── repository/         # Secondary adapters (persistence)
│       │   ├── db.go
│       │   ├── user_repository.go
│       │   └── image_repository.go
│       │
│       ├── storage/            # Secondary adapters (file storage)
│       │   └── local_storage.go
│       │
│       ├── token/              # Secondary adapters (JWT)
│       │   └── jwt_token_service.go
│       │
│       └── processor/          # Secondary adapters (image processing)
│           └── image_processor.go
│
├── migrations/                 # Database migrations
│   ├── 000001_user_schema.up.sql
│   └── 000002_image_schema.up.sql
│
├── storage/                    # Local file storage
└── .env                        # Environment configuration
```

---

## Dependency Flow

The architecture enforces a strict dependency rule: **dependencies point inward**.

```
┌────────────────────────────────────────────────────────────┐
│                    cmd/main.go                             │
│              (Composition Root - Wiring)                   │
└──────────────────────────┬─────────────────────────────────┘
                           │
                           │ creates and injects
                           ▼
┌────────────────────────────────────────────────────────────┐
│                    ADAPTERS                                │
│                                                            │
│  Primary: HTTP Handlers ────┐                              │
│                              │                              │
│  Secondary: Repository ◄─────┼────── implements ────────┐  │
│  Secondary: Storage ◄────────┼────── implements        │  │
│  Secondary: Token ◄──────────┼────── implements        │  │
│                              │                          │  │
└──────────────────────────────┼──────────────────────────┼──┘
                               │ depends on               │
                               ▼                          │
┌────────────────────────────────────────────────────────────┐
│                      PORTS                                 │
│                                                            │
│  UseCase interfaces ◄───────┐                              │
│  Repository interfaces ◄────┼──── defined by ───────────┐ │
│  Storage interfaces ◄───────┤                           │ │
│  Token interfaces ◄─────────┘                           │ │
│                                                         │ │
└─────────────────────────────────────────────────────────┼─┘
                                                          │
                                        implemented by ───┘
                                                          │
┌─────────────────────────────────────────────────────────┼─┐
│                      DOMAIN                              │ │
│                                                          │ │
│  Entities: User, Image ◄─────────── uses ───────────────┘ │
│  Value Objects: Email, Password, Format, etc.             │
│  Services: UserService, ImageService                      │
│                                                            │
│  (Zero dependencies on external world)                     │
└────────────────────────────────────────────────────────────┘
```

---

## Benefits of This Architecture

| Benefit                    | How This Project Achieves It                                       |
| -------------------------- | ------------------------------------------------------------------ |
| **Testability**            | Interfaces allow easy mocking; domain has no external dependencies |
| **Maintainability**        | Clear separation of concerns; changes are localized                |
| **Flexibility**            | Swap PostgreSQL for MongoDB by implementing UserRepository         |
| **Framework Independence** | Domain layer uses no HTTP or DB frameworks                         |
| **Domain Focus**           | Business rules are explicit and centralized                        |
| **Explicit Dependencies**  | All dependencies are injected; no hidden globals                   |

---

## **Build with ❤️ using Go**

---
