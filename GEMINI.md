# StoryLoom Project Mandates

This document serves as the foundational guide for all AI-assisted development on the StoryLoom project. These instructions take absolute precedence over general defaults.

## Project Overview
StoryLoom is a professional SaaS platform designed for long-form novel writers. It focuses on "Structural Logic" and "Creative Realism," providing a "digital editor" experience to help writers build coherent worlds, deep characters, and engaging plots.

## Tech Stack
- **Language:** Go (Golang) 1.21+
- **API Framework:** [go-swagger](https://github.com/go-swagger/go-swagger) (OpenAPI 2.0)
- **Web Server:** Gin or Echo (integrated with go-swagger middleware)
- **Database:** SQLite (Default for portability and speed)
- **ORM:** GORM with SQLite driver
- **AI Integration:** LLM APIs (GPT-4/Claude) for logic checking, plot suggestions, and character voice analysis.
- **Containerization:** Docker

## Architecture & Layout (DDD-Focused)
- **API-First Design:** `api/swagger.yaml` is the source of truth.
- **Directory Structure:**
  - `/api`: Swagger/OpenAPI specifications.
  - `/cmd`: Entry points for the server and CLI tools.
  - `/internal`: Core business logic organized by PRD modules:
    - `/internal/world`: **World-Building & Realism** (Logic settings, professional detail library).
    - `/internal/character`: **Character System** (POV, motivations, dynamic relationships, voice analysis).
    - `/internal/plot`: **Plot & Structure** (Card-based outlining, pacing intensity charts).
    - `/internal/conflict`: **Conflict & Theme** (Multi-dimensional conflict tracking, theme consistency).
    - `/internal/scene`: **Scene Orchestration** (POV validation, sensory detail templates, goals/obstacles).
    - `/internal/timeline`: **Timeline & Information Tracking** (Global timeline, character trackers, secret management).
    - `/internal/audit`: Logic and consistency checking (The "Digital Editor" engine).
  - `/pkg`: Shared utilities (database, AI client wrappers).
  - `/scripts`: Deployment and automation scripts.

## Engineering Conventions
- **Code Style:** Strictly adhere to [gofumpt](https://github.com/mvdan/gofumpt). All Go code must be formatted with `gofumpt -l -w .`.
- **Domain-Driven Design (DDD):** Logic must be encapsulated within its respective module in `/internal`.
- **Logic Validation:** Every module should have an "Editor Check" (AI-driven or rule-based) to ensure consistency (e.g., character POV consistency, timeline paradoxes).
- **Persistence:** Use GORM. Ensure schema migrations are handled cleanly for SQLite.
- **Performance (NFR):** API response times for scene/card operations should be < 500ms.

## Development Workflow
1. **API Update:** Modify `api/swagger.yaml` based on PRD REQs.
2. **Generation:** Run `swagger generate server -f api/swagger.yaml`.
3. **Implementation:** Implement domain logic in `/internal`.
4. **Formatting:** Run `gofumpt -l -w .`.
5. **Validation:** Run `go test ./...` and `swagger validate api/swagger.yaml`.

## AI Instructions
- **Priority:** Focus on the MVP modules (World, Character, Plot, Scene, Timeline).
- **Logic First:** When implementing features, consider how the "Digital Editor" would validate the new data (e.g., "Is this character in two places at once in the timeline?").
- **Clean Code:** Follow [Effective Go](https://go.dev/doc/effective_go). Use clear naming for domain entities as defined in the PRD (e.g., `POV`, `Cliffhanger`, `CoreQuestion`).
- **Surgical Edits:** When modifying code, ensure minimal disruption to other modules.
