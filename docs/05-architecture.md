# Architecture

This document describes the internal architecture of vulhub-cli for developers and contributors.

## Overview

vulhub-cli is built with a modular architecture using dependency injection. The codebase is organized into distinct packages, each with a specific responsibility.

```
cmd/vulhub/           # Application entry point
internal/
├── cli/              # CLI layer (commands, UI)
│   ├── commands/     # Command implementations
│   └── ui/           # User interface components
├── config/           # Configuration management
├── environment/      # Environment lifecycle management
├── github/           # GitHub API client
├── compose/          # Docker Compose wrapper
└── resolver/         # Keyword resolution
pkg/
└── types/            # Shared type definitions
```

## Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.21+ |
| CLI Framework | urfave/cli/v3 |
| Dependency Injection | uber-go/fx |
| GitHub API | google/go-github |
| Configuration | BurntSushi/toml |
| Utilities | samber/lo |

## Package Responsibilities

### cmd/vulhub

The application entry point. Responsibilities:
- Set up signal handling
- Initialize the fx application container
- Wire up dependencies
- Run the CLI

### internal/cli

The CLI layer handles user interaction.

#### internal/cli/commands

Each command is implemented as a method on the `Commands` struct:

```go
type Commands struct {
    Config       config.Manager
    Environment  environment.Manager
    Resolver     resolver.Resolver
    Downloader   *github.Downloader
    GitHubClient *github.GitHubClient
}

func (c *Commands) Start() *cli.Command { ... }
func (c *Commands) Stop() *cli.Command { ... }
```

Common patterns are extracted into helper functions:
- `ensureInitialized()` - Check/perform initialization
- `checkAndPromptSync()` - Check/perform sync
- `withRateLimitRetry()` - Handle rate limits with automatic retry

#### internal/cli/ui

UI components for terminal interaction:
- `Table` - Formatted table output
- `Selector` - Interactive selection menus
- `Pager` - Scrollable content display
- `Spinner` - Progress indicators

### internal/config

Configuration management responsibilities:
- Load/save configuration files
- Manage paths (config dir, environments dir)
- Track sync status

Key interface:
```go
type Manager interface {
    Get() *types.Config
    Set(cfg *types.Config)
    Save(ctx context.Context) error
    IsInitialized() bool
    NeedSync() bool
    LoadEnvironments(ctx context.Context) (*types.EnvironmentList, error)
    SaveEnvironments(ctx context.Context, list *types.EnvironmentList) error
    // ...
}
```

### internal/environment

Environment lifecycle management:

```go
type Manager interface {
    Start(ctx context.Context, env types.Environment, options StartOptions) error
    Stop(ctx context.Context, env types.Environment) error
    Restart(ctx context.Context, env types.Environment) error
    Down(ctx context.Context, env types.Environment) error
    Status(ctx context.Context, env types.Environment) (*types.EnvironmentStatus, error)
    GetInfo(ctx context.Context, env types.Environment) (*types.EnvironmentInfo, error)
    EnsureDownloaded(ctx context.Context, env types.Environment) (string, error)
    // ...
}
```

### internal/github

GitHub API interaction:

#### Client

Low-level GitHub API operations:
```go
type Client interface {
    DownloadFile(ctx context.Context, owner, repo, path, ref string) ([]byte, error)
    DownloadDirectory(ctx context.Context, owner, repo, path, ref, destDir string) error
    GetFileContent(ctx context.Context, owner, repo, path, ref string) (string, error)
    ListDirectoryContents(ctx context.Context, owner, repo, path, ref string) ([]types.ContentEntry, error)
}
```

#### Downloader

High-level download operations:
```go
type Downloader struct { ... }

func (d *Downloader) DownloadEnvironmentsList(ctx context.Context) ([]byte, error)
func (d *Downloader) DownloadEnvironment(ctx context.Context, env types.Environment, destDir string) error
func (d *Downloader) GetEnvironmentReadme(ctx context.Context, env types.Environment) (string, error)
```

#### OAuth

GitHub OAuth Device Flow implementation:
```go
func RequestDeviceCode(ctx context.Context) (*DeviceCodeResponse, error)
func PollForAccessToken(ctx context.Context, deviceCode string, interval int) (string, error)
```

### internal/compose

Docker Compose wrapper:

```go
type Client interface {
    Start(ctx context.Context, workDir string, options StartOptions) error
    Stop(ctx context.Context, workDir string, options StopOptions) error
    Restart(ctx context.Context, workDir string, options RestartOptions) error
    Down(ctx context.Context, workDir string, options DownOptions) error
    Status(ctx context.Context, workDir string) ([]types.ContainerStatus, error)
    CheckDocker(ctx context.Context) error
}
```

Executes `docker compose` commands via subprocess.

### internal/resolver

Keyword resolution for environment matching:

```go
type Resolver interface {
    Resolve(ctx context.Context, keyword string) (*ResolveResult, error)
    IsCVEFormat(keyword string) bool
    IsPathFormat(keyword string) bool
}
```

Match types:
1. **Exact CVE** - Direct CVE number match
2. **Exact Path** - Direct path match
3. **Fuzzy** - Partial match on app, path, or description

## Dependency Injection

The application uses uber-go/fx for dependency injection:

```go
app := fx.New(
    fx.Provide(func() *slog.Logger { return slog.Default() }),
    config.Module,
    github.Module,
    compose.Module,
    resolver.Module,
    environment.Module,
    fx.Provide(cli.NewApp),
    fx.Populate(&cliApp),
)
```

Each module provides its dependencies:
```go
// config/module.go
var Module = fx.Module("config",
    fx.Provide(NewConfigManager),
)

// github/module.go
var Module = fx.Module("github",
    fx.Provide(NewClientFromConfig),
    fx.Provide(NewDownloaderFromConfig),
)
```

## Data Flow

### Start Command Flow

```
User: vulhub start log4j
         │
         ▼
    ┌─────────────┐
    │   CLI       │  Parse arguments
    │  Command    │
    └─────┬───────┘
          │
          ▼
    ┌─────────────┐
    │  Resolver   │  Match "log4j" to environments
    └─────┬───────┘
          │
          ▼
    ┌─────────────┐
    │ Environment │  Check if downloaded
    │  Manager    │
    └─────┬───────┘
          │ (if not downloaded)
          ▼
    ┌─────────────┐
    │  GitHub     │  Download from repository
    │ Downloader  │
    └─────┬───────┘
          │
          ▼
    ┌─────────────┐
    │  Compose    │  docker compose up -d
    │  Client     │
    └─────────────┘
```

### Rate Limit Retry Flow

```
    ┌─────────────┐
    │  Operation  │  Any GitHub API call
    └─────┬───────┘
          │
          ▼
    ┌─────────────┐
    │ Rate Limit  │  Check for 403/rate limit
    │   Check     │
    └─────┬───────┘
          │ (if rate limited and no token)
          ▼
    ┌─────────────┐
    │   OAuth     │  Prompt and authenticate
    │    Flow     │
    └─────┬───────┘
          │
          ▼
    ┌─────────────┐
    │  Refresh    │  Update client token
    │   Client    │
    └─────┬───────┘
          │
          ▼
    ┌─────────────┐
    │   Retry     │  Automatic retry
    │ Operation   │
    └─────────────┘
```

## Error Handling

### Error Types

- **User errors**: Missing arguments, invalid keywords
- **Docker errors**: Docker not running, compose failures
- **GitHub errors**: Rate limits, network issues
- **Configuration errors**: Not initialized, corrupted config

### Rate Limit Handling

The `withRateLimitRetry` wrapper handles rate limits automatically:

```go
func (c *Commands) withRateLimitRetry(ctx context.Context, fn func() error) error {
    err := fn()
    if err == nil {
        return nil
    }

    if !github.IsRateLimitError(err) {
        return err
    }

    // Already has token? Return error
    if c.Config.Get().GitHub.Token != "" {
        return fmt.Errorf("rate limit exceeded even with token: %w", err)
    }

    // Prompt for authentication
    if c.PromptTokenSetup(ctx) {
        c.refreshGitHubClient()
        return fn()  // Automatic retry
    }

    return err
}
```

## Testing

### Unit Tests

```bash
go test ./...
```

### Integration Tests

Requires Docker:
```bash
go test -tags=integration ./...
```

## Adding New Commands

1. Create a new file in `internal/cli/commands/`:
```go
func (c *Commands) MyCommand() *cli.Command {
    return &cli.Command{
        Name:  "my-command",
        Usage: "Description",
        Action: func(ctx context.Context, cmd *cli.Command) error {
            return c.runMyCommand(ctx)
        },
    }
}

func (c *Commands) runMyCommand(ctx context.Context) error {
    // Implementation
}
```

2. Register in `commands.go`:
```go
func (c *Commands) All() []*cli.Command {
    return []*cli.Command{
        // ...
        c.MyCommand(),
    }
}
```

## Contributing

1. Follow existing code patterns
2. Add tests for new functionality
3. Update documentation
4. Run `go fmt` and `golangci-lint`
