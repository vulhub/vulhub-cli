# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

vulhub-cli is a Go command-line tool that simplifies interaction with the Vulhub security lab platform. It allows users to start, stop, and manage vulnerability environments without learning Docker Compose commands or downloading the complete Vulhub repository.

**Status**: Implementation in progress - core functionality is implemented.

## Technology Stack

- **Language**: Go 1.21+
- **CLI Framework**: `github.com/urfave/cli/v3`
- **Web Framework**: `github.com/gin-gonic/gin`
- **Dependency Injection**: `github.com/uber-go/fx`
- **TOML Parsing**: `github.com/BurntSushi/toml`
- **HTTP Client**: `resty.dev/v3`
- **GitHub API**: `github.com/google/go-github`
- **Utilities**: `github.com/samber/lo`
- **Logging**: `log/slog` (standard library)

## Architecture

The project follows a modular architecture with dependency injection:

```
cmd/vulhub/           # Main entry point
internal/
├── api/              # Web API server (Gin-based REST API)
├── cli/              # CLI command implementations (urfave/cli)
├── compose/          # Docker Compose bridge (wraps `docker compose` CLI)
├── github/           # GitHub API bridge (downloads from vulhub/vulhub repo)
├── config/           # Configuration management
├── environment/      # Environment lifecycle management
└── vulhub/           # Vulhub-specific business logic
```

### Key Patterns

1. **Bridge Pattern**: `compose/` and `github/` packages provide clean interfaces to external systems
2. **Dependency Injection**: Use uber-go/fx for lifecycle management and dependency wiring
3. **Interface-based Design**: Define interfaces in requirements.md (see Requirement 4 & 5)

### Configuration

User configuration stored in `~/.vulhub/`:
- `config.toml` - Main configuration
- `environments.toml` - Cached environment list from vulhub repository
- `environments/` - Downloaded environment files (docker-compose.yml, etc.)

## CLI Commands

| Command | Description |
|---------|-------------|
| `vulhub init` | Initialize configuration and download environment list |
| `vulhub syncup` | Update environment list from GitHub |
| `vulhub start [keyword]` | Start a vulnerability environment |
| `vulhub stop [keyword]` | Stop a running environment |
| `vulhub restart [keyword]` | Restart an environment |
| `vulhub clean [keyword]` | Completely remove an environment (containers, volumes, local files) |
| `vulhub status [keyword]` | Show status of downloaded environments (aliases: `ls`, `list`) |
| `vulhub list-available` | List all available environments |
| `vulhub search [keyword]` | Search for environments |
| `vulhub info [keyword]` | Show environment details |
| `vulhub github-auth` | Authenticate with GitHub using OAuth Device Flow |
| `vulhub doctor` | Check system environment and diagnose potential issues |
| `vulhub serve` | Start the web API server |

Keywords support: exact CVE numbers (`CVE-2021-44228`), exact paths (`log4j/CVE-2021-44228`), or fuzzy matching (`log4j`).

## Web API

The `vulhub serve` command starts a REST API server that exposes all core functionality via HTTP endpoints.

### Starting the Server

```bash
# Start on default port 8080
vulhub serve

# Start on custom host and port
vulhub serve --host 127.0.0.1 --port 3000
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/status` | System status |
| POST | `/api/v1/syncup` | Sync environment list from GitHub |
| GET | `/api/v1/environments` | List all available environments |
| GET | `/api/v1/environments/downloaded` | List downloaded environments |
| GET | `/api/v1/environments/running` | List running environments |
| GET | `/api/v1/environments/info/*path` | Get environment info |
| GET | `/api/v1/environments/status/*path` | Get environment status |
| POST | `/api/v1/environments/start/*path` | Start an environment |
| POST | `/api/v1/environments/stop/*path` | Stop an environment |
| POST | `/api/v1/environments/restart/*path` | Restart an environment |
| DELETE | `/api/v1/environments/clean/*path` | Clean an environment |

### Example Usage

```bash
# List all environments
curl http://localhost:8080/api/v1/environments

# Start an environment
curl -X POST http://localhost:8080/api/v1/environments/start/log4j/CVE-2021-44228

# Stop an environment
curl -X POST http://localhost:8080/api/v1/environments/stop/log4j/CVE-2021-44228
```

## Build Commands

Once Go project is initialized:

```bash
# Build
go build -o vulhub ./cmd/vulhub

# Test
go test ./...

# Test with verbose output
go test -v ./...

# Run single test
go test -v -run TestFunctionName ./path/to/package

# Lint (if golangci-lint is installed)
golangci-lint run
```
