# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

vulhub-cli is a Go command-line tool that simplifies interaction with the Vulhub security lab platform. It allows users to start, stop, and manage vulnerability environments without learning Docker Compose commands or downloading the complete Vulhub repository.

**Status**: Design phase - `requirements.md` contains detailed specifications but implementation has not started.

## Technology Stack

- **Language**: Go 1.21+
- **CLI Framework**: `github.com/urfave/cli/v3`
- **Dependency Injection**: `github.com/uber-go/fx`
- **TOML Parsing**: `github.com/BurntSushi/toml`
- **HTTP Client**: `github.com/go-resty/resty/v3`
- **GitHub API**: `github.com/google/go-github`
- **Utilities**: `github.com/samber/lo`
- **Logging**: `log/slog` (standard library)

## Architecture

The project follows a modular architecture with dependency injection:

```
cmd/vulhub/           # Main entry point
internal/
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
| `vulhub list` | List running environments |
| `vulhub list-available` | List all available environments |
| `vulhub search [keyword]` | Search for environments |
| `vulhub info [keyword]` | Show environment details |
| `vulhub clean [keyword]` | Clean up resources |

Keywords support: exact CVE numbers (`CVE-2021-44228`), exact paths (`log4j/CVE-2021-44228`), or fuzzy matching (`log4j`).

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
