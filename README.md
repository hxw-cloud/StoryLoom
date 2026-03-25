# StoryLoom

**StoryLoom** is a professional storytelling and world-building platform designed as a "digital editor" for long-form novel writers. It emphasizes structural logic and creative realism, helping writers maintain consistency across complex narratives.

## 🌟 Core Features (PRD v3.0)

StoryLoom is divided into six core domain modules:

1. **World-Building & Realism (`world`)**: Build coherent settings with a professional detail library and historical event tracking. Ensures "Logic Settings" are respected throughout the story.
2. **Character System (`character`)**: Deep character profiles, Point of View (POV) settings, motivation/conflict tracking, and dynamic relationship mapping.
3. **Plot & Structure (`plot`)**: Card-based dynamic outlining, customizable templates (e.g., Save the Cat, Hero's Journey), and pacing intensity tracking.
4. **Conflict & Theme (`conflict`)**: Multi-dimensional conflict tracking (e.g., Man vs. Self, Man vs. Society) to ensure the core theme is consistently delivered.
5. **Scene Orchestration (`scene`)**: Guided scene writing templates validating POV, sensory details, and goal/obstacle resolution.
6. **Timeline & Information (`timeline`)**: A global master timeline synchronizing character movements, events, and tracking "who knows what and when" (secret management).

## 🛠️ Tech Stack

- **Language**: Go (Golang) 1.21+
- **API Framework**: go-swagger (OpenAPI 2.0)
- **Database**: SQLite (Local, portable persistence)
- **ORM**: GORM
- **Architecture**: Domain-Driven Design (DDD) & API-First

## 🚀 Getting Started

### Prerequisites
- [Go 1.21+](https://go.dev/doc/install)
- [go-swagger](https://goswagger.io/install.html) (`go install github.com/go-swagger/go-swagger/cmd/swagger@latest`)
- [gofumpt](https://github.com/mvdan/gofumpt) (`go install mvdan.cc/gofumpt@latest`)

### Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/hxw-cloud/StoryLoom.git
   cd StoryLoom
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the server:**
   ```bash
   go run cmd/storyloom-server/main.go --port=8080
   ```
   The API will be available at `http://localhost:8080/api/v1`.

### Development Workflow

StoryLoom strictly follows an **API-First Design**:
1. Modify the `api/swagger.yaml` file to define new endpoints or models.
2. Generate the server boilerplate:
   ```bash
   swagger generate server -f api/swagger.yaml -A storyloom
   ```
3. Implement the business logic in the respective `/internal` module.
4. Format your code before committing:
   ```bash
   gofumpt -l -w .
   ```

## 📄 License
[MIT License](LICENSE)
