# StoryLoom Project Mandates

This document serves as the foundational guide for all AI-assisted development on the StoryLoom project. These instructions take absolute precedence over general defaults.

## Project Overview
StoryLoom is a comprehensive storytelling and world-building platform designed for writers, world-builders, and narrative designers.

## Tech Stack
- **Language:** Go (Golang) 1.21+
- **API Framework:** [go-swagger](https://github.com/go-swagger/go-swagger) (OpenAPI 2.0)
- **Web Server:** Gin or Echo (integrated with go-swagger middleware)
- **Database:** SQLite
- **ORM:** GORM with SQLite driver
- **Containerization:** Docker

## Architecture & Layout
- **API-First Design:** `api/swagger.yaml` is the source of truth. Always update the specification before modifying API-related code.
- **Directory Structure:**
  - `/api`: Swagger/OpenAPI specifications.
  - `/cmd`: Entry points for applications.
  - `/internal`: Private application and library code.
    - `/internal/domain`: Core domain models and business logic.
    - `/internal/world`: World-building module.
    - `/internal/character`: Character management module.
    - `/internal/plot`: Plot development module.
    - `/internal/conflict`: Conflict tracking module.
    - `/internal/scene`: Scene-by-scene orchestration module.
    - `/internal/timeline`: Chronological event management module.
  - `/pkg`: Public library code.
  - `/scripts`: Shell and Batch scripts for development/deployment.

## Engineering Conventions
- **Code Style:** Strictly adhere to [gofumpt](https://github.com/mvdan/gofumpt) for formatting. All Go code must be formatted with `gofumpt -l -w .`.
- **Domain-Driven Design (DDD):** Organize logic by the six core modules (World, Character, Plot, Conflict, Scene, Timeline).
- **Code Generation:** Use `swagger generate server` for API boilerplate. Do not manually edit generated files (typically in `restapi/`).
- **Persistence:** Use GORM for all database interactions. Keep SQLite logic portable.
- **Testing:** Unit tests should accompany all business logic in `/internal`. Use `go test ./...`.

## Development Workflow
1. **API Update:** Modify `api/swagger.yaml`.
2. **Generation:** Run `swagger generate server -f api/swagger.yaml`.
3. **Implementation:** Implement handlers and business logic in `/internal`.
4. **Formatting:** Run `gofumpt -l -w .` before committing or testing.
5. **Validation:** Run tests and ensure `swagger validate api/swagger.yaml` passes.

## AI Instructions
- Always verify `go-swagger` is installed before attempting code generation.
- Prioritize clean, idiomatic Go code following [Effective Go](https://go.dev/doc/effective_go).
- When asked to add features, start by proposing changes to the `swagger.yaml` first.
