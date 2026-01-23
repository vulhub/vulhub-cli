# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

vulhub-cli is a Go command-line tool that simplifies interaction with the Vulhub security lab platform. It allows users to start, stop, and manage vulnerability environments without learning Docker Compose commands or downloading the complete Vulhub repository.

**Status**: Implementation in progress - core functionality is implemented.

## Technology Stack

- **Language**: Go 1.21+
- **CLI Framework**: `github.com/urfave/cli/v3`
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
| `vulhub down [keyword]` | Completely remove an environment (containers, volumes, local files) |
| `vulhub restart [keyword]` | Restart an environment |
| `vulhub list` | List all downloaded environments |
| `vulhub list-available` | List all available environments |
| `vulhub search [keyword]` | Search for environments |
| `vulhub info [keyword]` | Show environment details |
| `vulhub clean [keyword]` | Clean up resources |
| `vulhub github-auth` | Authenticate with GitHub using OAuth Device Flow |

Keywords support: exact CVE numbers (`CVE-2021-44228`), exact paths (`log4j/CVE-2021-44228`), or fuzzy matching (`log4j`).

## GitHub Authentication

The tool downloads files from GitHub, which has API rate limits (60 requests/hour for unauthenticated users).

### OAuth Device Flow Authentication

The CLI uses GitHub OAuth Device Flow for authentication, providing a seamless experience:

```bash
vulhub github-auth           # Start OAuth authentication flow
vulhub github-auth --status  # Check current authentication status
vulhub github-auth --remove  # Remove saved authentication
```

**How it works:**
1. Run `vulhub github-auth`
2. A browser opens automatically to GitHub's device authorization page
3. Enter the displayed code on the GitHub page
4. Authorize the Vulhub CLI OAuth App
5. Authentication completes automatically

### Automatic Rate Limit Handling

When a rate limit error occurs and the user is not authenticated, the CLI automatically prompts to start the OAuth flow.

### Token Storage

- OAuth token is saved in `~/.vulhub/config.toml`
- Environment variable `GITHUB_TOKEN` takes precedence over saved token

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
