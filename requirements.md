# Requirements Document

## Introduction

This document defines the requirements for the Vulhub CLI tool. The tool aims to simplify the usage of the Vulhub security lab platform, allowing users to start specific vulnerability environments without learning Docker Compose commands or downloading the complete Vulhub project.

## Technology Stack

### Programming Language

- **Go** (1.21+)

### Third-Party Libraries

If the following functionalities are required, the listed third-party libraries SHOULD be preferred to implement them:

| Category | Library | Purpose |
|----------|---------|---------|
| CLI Framework | `github.com/urfave/cli/v3` | Command-line interface implementation |
| Dependency Injection | `github.com/uber-go/fx` | Dependency management and lifecycle |
| TOML Parsing | `github.com/BurntSushi/toml` | Parsing `environments.toml` and `config.toml` |
| HTTP Client | `resty.dev/v3` | Enhanced HTTP operations |
| GitHub Client | `github.com/google/go-github` | GitHub API integration |
| Utility Functions | `github.com/samber/lo` | Functional programming helpers |

### Standard Library Packages

The following Go standard library packages SHALL be used:

| Category | Packages |
|----------|----------|
| File Operations | `os`, `path/filepath`, `io` |
| Concurrency | `context`, `sync` |
| Logging | `log/slog` |
| Testing | `testing` |

### Library Selection Guidelines

1. When implementing any feature, MUST first check if the required functionality can be achieved using the libraries listed above
2. If a library is not listed but is required, preference SHALL be given to well-maintained, widely-adopted libraries
3. Standard library packages SHALL be preferred over third-party alternatives when functionality is equivalent

---

## Requirements

### Requirement 1: Initialization Commands

**User Story:** As a user, I want to initialize the Vulhub CLI tool and sync the environment list, so that I can have the latest vulnerability environments available locally.

#### Commands

- `vulhub init` - Initialize the CLI tool configuration
- `vulhub syncup` - Sync/update the environment list from remote

#### Acceptance Criteria

1. WHEN user executes `vulhub init` THEN system SHALL create initialization configuration files in user's home directory (e.g., `~/.vulhub/config.toml`)
2. WHEN user executes `vulhub init` THEN system SHALL download the official environment list from https://github.com/vulhub/vulhub/blob/master/environments.toml
3. WHEN configuration files already exist THEN system SHALL prompt user whether to overwrite existing configuration
4. WHEN user executes `vulhub syncup` THEN system SHALL fetch the latest environment list from remote and update local cache
5. WHEN initialization completes successfully THEN system SHALL display success message with basic usage instructions

#### Configuration Files

The following files SHALL be created in `~/.vulhub/` directory:

- `config.toml` - Main configuration file
- `environments.toml` - Cached environment list from vulhub repository
- `environments/` - Directory for caching downloaded environment files

---

### Requirement 2: CLI Commands and Arguments

**User Story:** As a user, I want to use simple command-line commands to manage vulnerability environments.

#### Commands

| Command | Description |
|---------|-------------|
| `vulhub start [keyword]` | Start a vulnerability environment |
| `vulhub stop [keyword]` | Stop a running vulnerability environment |
| `vulhub down [keyword]` | Completely remove an environment (containers, volumes, and local files) |
| `vulhub restart [keyword]` | Restart a vulnerability environment |
| `vulhub clean [keyword]` | Clean up environment files and Docker resources |
| `vulhub list` | List all downloaded vulnerability environments |
| `vulhub list-available` | List all available vulnerability environments |
| `vulhub status` | Display status of all environments |
| `vulhub search [keyword]` | Search for vulnerability environments |
| `vulhub info [keyword]` | Display detailed information about a vulnerability environment |
| `vulhub github-auth` | Authenticate with GitHub using OAuth Device Flow |
| `vulhub doctor` | Check system environment and diagnose potential issues |

#### Acceptance Criteria

1. WHEN user executes any command THEN system SHALL parse command-line arguments correctly
2. WHEN user provides `--help` or `-h` flag THEN system SHALL display help information for the command
3. WHEN user provides invalid arguments THEN system SHALL display error message with correct usage
4. WHEN command requires a keyword but none is provided THEN system SHALL prompt user or display error message
5. This requirement SHALL only implement command-line parsing logic; backend execution logic is NOT in scope

#### Global Flags

- `--help, -h` - Display help information
- `--version, -v` - Display version information
- `--verbose` - Enable verbose output
- `--config <path>` - Specify custom configuration file path

---

### Requirement 3: Keyword Resolution Logic

**User Story:** As a user, I want to specify environments using flexible keywords (CVE number, path, or application name), and the system should help me find the correct environment.

#### Keyword Types

1. **Exact CVE Number** - e.g., `CVE-2021-44228`
2. **Exact Path** - e.g., `log4j/CVE-2021-44228`
3. **Fuzzy Keyword** - e.g., `log4j`, `apache`, `tomcat`

#### Acceptance Criteria

1. WHEN user provides an exact CVE number (format: `CVE-YYYY-XXXXX`) THEN system SHALL search environments.toml for matching CVE
2. WHEN user provides an exact path that exists in environments.toml THEN system SHALL directly use that environment
3. WHEN exact CVE or path is found THEN system SHALL directly execute the corresponding operation (start/stop/etc.)
4. WHEN user provides a fuzzy keyword (application name, partial match) THEN system SHALL list all matching environments
5. WHEN multiple environments match the keyword THEN system SHALL display an interactive selection menu
6. WHEN no environment matches the keyword THEN system SHALL display error message and suggest using `vulhub search` command
7. Interactive selection SHALL support keyboard navigation (arrow keys) and number input

#### Matching Priority

1. Exact CVE match (highest priority)
2. Exact path match
3. Application name match (case-insensitive)
4. Partial name/description match (lowest priority)

---

### Requirement 4: Docker Compose Bridge Package

**User Story:** As a developer, I want a standalone package that provides a bridge to Docker Compose commands, so that the main CLI can use Docker Compose functionality without tight coupling.

#### Package Name

`dockercompose` (internal package)

#### Interface Design

```go
package dockercompose

type ComposeClient interface {
    // Start starts the Docker Compose environment in the specified directory
    Start(ctx context.Context, workDir string, options StartOptions) error
    
    // Stop stops the Docker Compose environment
    Stop(ctx context.Context, workDir string, options StopOptions) error
    
    // Restart restarts the Docker Compose environment
    Restart(ctx context.Context, workDir string, options RestartOptions) error
    
    // Status returns the status of containers in the environment
    Status(ctx context.Context, workDir string) ([]ContainerStatus, error)
    
    // Logs retrieves logs from containers
    Logs(ctx context.Context, workDir string, options LogsOptions) (io.ReadCloser, error)
    
    // Down stops and removes containers, networks, and volumes
    Down(ctx context.Context, workDir string, options DownOptions) error
    
    // Pull pulls images for the environment
    Pull(ctx context.Context, workDir string) error
}
```

#### Acceptance Criteria

1. Package SHALL provide a clean interface for executing Docker Compose commands
2. Package SHALL wrap `docker compose` CLI commands (not docker-compose v1)
3. Package SHALL handle command output (stdout/stderr) appropriately
4. Package SHALL return structured errors with meaningful messages
5. Package SHALL support context cancellation for long-running operations
6. Package SHALL check if Docker daemon is running before executing commands
7. Package SHALL be independently testable with mock implementations

---

### Requirement 5: GitHub Bridge Package

**User Story:** As a developer, I want a standalone package that provides a bridge to GitHub API, so that the main CLI can download files and access repository information.

#### Package Name

`github` (internal package)

#### Interface Design

```go
package github

type GitHubClient interface {
    // DownloadFile downloads a single file from a repository
    DownloadFile(ctx context.Context, owner, repo, path, ref string) ([]byte, error)
    
    // DownloadDirectory downloads all files in a directory from a repository
    DownloadDirectory(ctx context.Context, owner, repo, path, ref, destDir string) error
    
    // GetFileContent gets the content of a file (for small files)
    GetFileContent(ctx context.Context, owner, repo, path, ref string) (string, error)
    
    // GetLatestRelease gets the latest release information
    GetLatestRelease(ctx context.Context, owner, repo string) (*Release, error)
    
    // ListDirectoryContents lists contents of a directory in the repository
    ListDirectoryContents(ctx context.Context, owner, repo, path, ref string) ([]ContentEntry, error)
}
```

#### Acceptance Criteria

1. Package SHALL provide a clean interface for GitHub API operations
2. Package SHALL use GitHub REST API for file downloads and repository access
3. Package SHALL support downloading raw file content from repositories
4. Package SHALL handle rate limiting gracefully with appropriate retry logic
5. Package SHALL support optional GitHub token authentication for higher rate limits
6. Package SHALL cache downloaded files to reduce API calls
7. Package SHALL return structured errors with meaningful messages
8. Package SHALL be independently testable with mock implementations

#### Configuration

- GitHub API base URL: `https://api.github.com`
- Raw content URL: `https://raw.githubusercontent.com`
- Optional authentication via `GITHUB_TOKEN` environment variable

---

### Requirement 6: Subcommand Implementation

**User Story:** As a user, I want all CLI subcommands to work correctly, utilizing the Docker Compose and GitHub packages for their operations.

#### Subcommand: `init`

1. Check if configuration directory exists (`~/.vulhub/`)
2. Create configuration directory and default config file if not exists
3. Call GitHub package to download `environments.toml` from vulhub repository
4. Save environment list to local cache
5. Display success message

#### Subcommand: `syncup`

1. Call GitHub package to fetch latest `environments.toml`
2. Compare with local cached version
3. Update local cache if remote version is newer
4. Display update summary (new/updated/removed environments)

#### Subcommand: `start [keyword]`

1. Resolve keyword using Requirement 3 logic
2. If multiple matches, display interactive selection
3. Call GitHub package to download environment files (docker-compose.yml, etc.)
4. Call Docker Compose package to start the environment
5. Display access information (URLs, ports)

#### Subcommand: `stop [keyword]`

1. Resolve keyword using Requirement 3 logic
2. Call Docker Compose package to stop the environment
3. Display confirmation message

#### Subcommand: `restart [keyword]`

1. Resolve keyword using Requirement 3 logic
2. Call Docker Compose package to restart the environment
3. Display confirmation message

#### Subcommand: `clean [keyword]`

1. Resolve keyword using Requirement 3 logic
2. If multiple matches, display interactive selection
3. Call Docker Compose package to stop and remove containers, networks, and volumes (`docker compose down -v`)
4. Remove local environment files (docker-compose.yml, etc.)
5. Display confirmation message

#### Subcommand: `list`

1. Scan local cache directory for downloaded environments
2. Call Docker Compose package to get status of each environment
3. Display list of environments with their status (running/stopped)

#### Subcommand: `status`

1. Call Docker Compose package to get status of all running containers
2. Display detailed status information (container names, ports, uptime)

#### Subcommand: `search [keyword]`

1. Load environment list from local cache
2. Search by CVE, application name, tags, and description
3. Display matching environments in formatted table

#### Subcommand: `info [keyword]`

1. Resolve keyword to specific environment
2. Display detailed information:
   - Environment name
   - CVE number(s)
   - Application name
   - Tags
   - Path in vulhub repository
   - README content (if available)

#### Subcommand: `github-auth`

1. If `--status` flag is provided, display current authentication status
2. If `--remove` flag is provided, prompt for confirmation and remove saved token
3. Otherwise, initiate OAuth Device Flow authentication:
   a. Request device code from GitHub OAuth endpoint
   b. Display user code and verification URL
   c. Automatically open browser to verification URL
   d. Poll for access token while user completes authorization
   e. Save access token to configuration file
4. Display success message with rate limit information

#### Subcommand: `doctor`

1. Check Docker environment:
   - Verify Docker is installed and in PATH
   - Verify Docker daemon is running (`docker info`)
   - Verify Docker Compose is available (`docker compose version`)
2. Check configuration:
   - Verify config directory exists (`~/.vulhub/`)
   - Validate `config.toml` syntax
   - Validate `environments.toml` syntax
   - Check environments directory status
3. Check network connectivity:
   - Test DNS resolution and HTTP connectivity to GitHub API (`api.github.com`)
   - Test connectivity to GitHub Raw Content (`raw.githubusercontent.com`)
   - Test connectivity to Docker Hub (`hub.docker.com`)
   - Test connectivity to Docker Registry (`registry-1.docker.io`)
   - Test connectivity to Docker Auth (`auth.docker.io`)
   - Test connectivity to Docker CDN (`production.cloudflare.docker.com`)
4. Check Docker registry:
   - Display configured registry mirrors
   - Perform actual image pull test using `hello-world:latest`
5. Display summary with pass/warning/error counts
6. If `--fix` flag is provided, attempt to fix issues automatically (e.g., create missing directories)
7. If `--verbose` flag is provided, display detailed timing and path information

#### Acceptance Criteria

1. All subcommands SHALL use the Docker Compose package (Requirement 4) for Docker operations
2. All subcommands SHALL use the GitHub package (Requirement 5) for repository access
3. All subcommands SHALL handle errors gracefully and display user-friendly messages
4. All subcommands SHALL respect the `--verbose` flag for detailed output
5. All subcommands SHALL support `--help` flag for command-specific help
6. Long-running operations SHALL display progress indicators
7. Interactive selections SHALL be skippable with `--yes` or `-y` flag for automation

---

### Requirement 7: GitHub OAuth Device Flow Authentication

**User Story:** As a user, I want to authenticate with GitHub easily without manually creating tokens, so that I can avoid API rate limits with minimal effort.

#### OAuth App Configuration

- **OAuth App Name**: Vulhub CLI
- **Client ID**: `Ov23liDeiHCLOTtZxFY4`
- **Authorization Flow**: Device Flow (RFC 8628)
- **Required Scope**: `public_repo`

#### Command

- `vulhub github-auth` - Authenticate with GitHub using OAuth Device Flow

#### Command Flags

| Flag | Description |
|------|-------------|
| `--status` | Display current authentication status |
| `--remove` | Remove saved GitHub authentication |

#### Acceptance Criteria

1. WHEN user executes `vulhub github-auth` THEN system SHALL initiate OAuth Device Flow
2. WHEN device code is received THEN system SHALL display user code and verification URL
3. WHEN verification URL is displayed THEN system SHALL attempt to open default browser automatically
4. WHEN waiting for authorization THEN system SHALL poll GitHub at the specified interval
5. WHEN user completes authorization THEN system SHALL save the access token to config file
6. WHEN user denies authorization THEN system SHALL display appropriate message
7. WHEN device code expires THEN system SHALL inform user and suggest retrying
8. WHEN user executes `vulhub github-auth --status` THEN system SHALL display authentication status with masked token
9. WHEN user executes `vulhub github-auth --remove` THEN system SHALL prompt for confirmation before removing token
10. WHEN any command encounters GitHub API rate limit AND user is not authenticated THEN system SHALL prompt to start OAuth flow

#### OAuth Device Flow Sequence

1. **Request Device Code**
   - POST to `https://github.com/login/device/code`
   - Parameters: `client_id`, `scope=public_repo`
   - Response: `device_code`, `user_code`, `verification_uri`, `expires_in`, `interval`

2. **User Authorization**
   - Display `user_code` to user
   - Direct user to `verification_uri` (https://github.com/login/device)
   - Open browser automatically if possible

3. **Poll for Access Token**
   - POST to `https://github.com/login/oauth/access_token`
   - Parameters: `client_id`, `device_code`, `grant_type=urn:ietf:params:oauth:grant-type:device_code`
   - Poll at `interval` seconds until success or expiration
   - Handle `authorization_pending`, `slow_down`, `expired_token`, `access_denied` responses

4. **Store Token**
   - Save access token to `~/.vulhub/config.toml` under `[github]` section
   - Token format: OAuth access token (starts with `gho_`)

#### Implementation Notes

- OAuth implementation is in `internal/github/oauth.go`
- Device flow callbacks allow UI customization
- Browser opening is platform-specific (darwin: `open`, windows: `cmd /c start`, linux: `xdg-open`)
- Environment variable `GITHUB_TOKEN` takes precedence over saved OAuth token
