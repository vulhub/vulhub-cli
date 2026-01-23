# Command Reference

This document provides detailed information about all vulhub-cli commands, their options, and internal behavior.

## Global Options

These options are available for all commands:

| Option | Description |
|--------|-------------|
| `--verbose` | Enable verbose/debug output |
| `--config <path>` | Specify a custom configuration file path |
| `--help, -h` | Show help information |
| `--version, -v` | Show version information |

---

## init

Initialize vulhub-cli configuration and download the environment list.

### Usage

```bash
vulhub init [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--force, -f` | Force overwrite existing configuration |

### Behavior

1. **Check existing configuration**: If already initialized, prompts for confirmation (unless `--force` is used)
2. **Create config directory**: Creates `~/.vulhub/` if it doesn't exist
3. **Create default config**: Generates `config.toml` with default settings
4. **Download environment list**: Fetches `environments.toml` from GitHub
5. **Create environments directory**: Creates `~/.vulhub/environments/` for downloaded environments

### Example

```bash
# First-time initialization
vulhub init

# Force re-initialization
vulhub init --force
```

---

## syncup

Synchronize the local environment list with the latest version from GitHub.

### Usage

```bash
vulhub syncup
```

### Options

None.

### Behavior

1. **Check initialization**: Fails if vulhub-cli is not initialized
2. **Load current list**: Reads the existing `environments.toml`
3. **Download latest list**: Fetches the current version from GitHub
4. **Compare and update**: Saves the new list and updates sync timestamp
5. **Report changes**: Shows the difference in environment count

### Example

```bash
vulhub syncup
# Output:
# ℹ Downloading latest environment list from GitHub...
# ✓ Environment list updated successfully!
# Previous: 175 environments
# Current:  180 environments
# Added:    5 new environments
```

---

## start

Start a vulnerability environment.

### Usage

```bash
vulhub start [options] <keyword>
```

### Options

| Option | Description |
|--------|-------------|
| `--yes, -y` | Skip confirmation prompts |
| `--pull` | Pull images before starting |
| `--build` | Build images before starting |
| `--force-recreate` | Force recreate containers even if unchanged |

### Arguments

| Argument | Description |
|----------|-------------|
| `keyword` | CVE number, path, or application name to match |

### Behavior

1. **Check Docker**: Verifies Docker and Docker Compose are available
2. **Ensure initialization**: Prompts to initialize if not done
3. **Check sync status**: Prompts to sync if environment list is outdated (>7 days)
4. **Resolve keyword**: Matches the keyword to environments
   - Exact CVE match (e.g., `CVE-2021-44228`)
   - Exact path match (e.g., `log4j/CVE-2021-44228`)
   - Fuzzy match on app name, path, or description
5. **Handle multiple matches**: If multiple environments match, displays a selection menu
6. **Download environment**: If not already downloaded, fetches files from GitHub
7. **Start containers**: Runs `docker compose up -d` with specified options
8. **Display status**: Shows running containers and exposed ports

### Example

```bash
# Start by CVE number
vulhub start CVE-2021-44228

# Start with image pull
vulhub start --pull log4j

# Start non-interactively (fails if multiple matches)
vulhub start -y struts2/s2-045
```

---

## stop

Stop a running vulnerability environment.

### Usage

```bash
vulhub stop [options] <keyword>
```

### Options

| Option | Description |
|--------|-------------|
| `--yes, -y` | Skip confirmation prompts |

### Arguments

| Argument | Description |
|----------|-------------|
| `keyword` | CVE number, path, or application name to match |

### Behavior

1. **Ensure initialization**: Checks if vulhub-cli is initialized
2. **Check sync status**: Prompts to sync if outdated
3. **Resolve keyword**: Matches the keyword to environments
4. **Handle multiple matches**: Displays selection menu if needed
5. **Stop containers**: Runs `docker compose stop`

### Example

```bash
vulhub stop CVE-2021-44228
```

---

## restart

Restart a vulnerability environment.

### Usage

```bash
vulhub restart [options] <keyword>
```

### Options

| Option | Description |
|--------|-------------|
| `--yes, -y` | Skip confirmation prompts |

### Arguments

| Argument | Description |
|----------|-------------|
| `keyword` | CVE number, path, or application name to match |

### Behavior

1. **Ensure initialization**: Checks if vulhub-cli is initialized
2. **Check sync status**: Prompts to sync if outdated
3. **Resolve keyword**: Matches the keyword to environments
4. **Handle multiple matches**: Displays selection menu if needed
5. **Restart containers**: Runs `docker compose restart`
6. **Display status**: Shows container status and ports

### Example

```bash
vulhub restart log4j
```

---

## clean

Completely remove an environment, including containers, volumes, and local files.

### Usage

```bash
vulhub clean [options] <keyword>
```

### Options

| Option | Description |
|--------|-------------|
| `--yes, -y` | Skip confirmation prompts |

### Arguments

| Argument | Description |
|----------|-------------|
| `keyword` | CVE number, path, or application name to match |

### Behavior

1. **Ensure initialization**: Checks if vulhub-cli is initialized
2. **Check sync status**: Prompts to sync if outdated
3. **Resolve keyword**: Matches the keyword to environments
4. **Handle multiple matches**: Displays selection menu if needed
5. **Stop and remove containers**: Runs `docker compose down -v`
6. **Remove local files**: Deletes the environment directory from `~/.vulhub/environments/`

### Example

```bash
# Interactive cleanup
vulhub clean CVE-2021-44228

# Non-interactive cleanup
vulhub clean -y log4j/CVE-2021-44228
```

---

## list

List all downloaded vulnerability environments.

### Usage

```bash
vulhub list
vulhub ls  # Alias
```

### Options

None.

### Behavior

1. **Ensure initialization**: Checks if vulhub-cli is initialized
2. **Scan environments directory**: Lists all environments in `~/.vulhub/environments/`
3. **Check status**: For each environment, checks if containers are running
4. **Display table**: Shows environment path, status, and container information

### Example

```bash
vulhub list
# Output:
# Downloaded Environments:
#   PATH                      STATUS    CONTAINERS
#   log4j/CVE-2021-44228      running   1
#   struts2/s2-045            stopped   0
```

---

## list-available

List all available vulnerability environments from the environment list.

### Usage

```bash
vulhub list-available
vulhub la  # Alias
```

### Options

None.

### Behavior

1. **Ensure initialization**: Checks if vulhub-cli is initialized
2. **Load environment list**: Reads `environments.toml`
3. **Display paginated list**: Shows all environments with CVE, app name, and path

### Example

```bash
vulhub list-available
```

---

## status

Show the status of vulnerability environments.

### Usage

```bash
vulhub status [keyword]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `keyword` | (Optional) Specific environment to check |

### Behavior

**Without keyword:**
1. Scans all downloaded environments
2. Shows only running environments with their container status

**With keyword:**
1. Resolves the keyword to a specific environment
2. Shows detailed status including:
   - Environment path and CVE
   - Download status
   - Container states and ports

### Example

```bash
# Show all running environments
vulhub status

# Show specific environment status
vulhub status log4j
```

---

## search

Search for vulnerability environments by keyword.

### Usage

```bash
vulhub search <keyword>
```

### Arguments

| Argument | Description |
|----------|-------------|
| `keyword` | Search term (CVE, app name, technology, etc.) |

### Behavior

1. **Ensure initialization**: Checks if vulhub-cli is initialized
2. **Search environments**: Matches keyword against:
   - CVE numbers
   - Application names
   - Environment paths
   - Descriptions
3. **Rank results**: Orders by match relevance
4. **Display results**: Shows matching environments with highlights

### Example

```bash
vulhub search spring
# Output:
# Search results for 'spring':
#   spring/CVE-2022-22965    Spring Framework RCE
#   spring/CVE-2022-22963    Spring Cloud Function SpEL Injection
#   ...
```

---

## info

Show detailed information about a vulnerability environment.

### Usage

```bash
vulhub info [options] <keyword>
```

### Options

| Option | Description |
|--------|-------------|
| `--yes, -y` | Skip confirmation prompts |
| `--no-readme` | Do not show README content |

### Arguments

| Argument | Description |
|----------|-------------|
| `keyword` | CVE number, path, or application name to match |

### Behavior

1. **Ensure initialization**: Checks if vulhub-cli is initialized
2. **Check sync status**: Prompts to sync if outdated
3. **Resolve keyword**: Matches the keyword to environments
4. **Handle multiple matches**: Displays selection menu if needed
5. **Fetch information**: Downloads README and docker-compose.yml from GitHub
6. **Display with pager**: Shows information in a scrollable view

### Displayed Information

- Environment path
- CVE number(s)
- Application name and version
- Download status
- README content (exploitation instructions)
- Docker Compose configuration

### Example

```bash
vulhub info CVE-2021-44228

# Without README
vulhub info --no-readme log4j
```

---

## github-auth

Authenticate with GitHub to increase API rate limits.

### Usage

```bash
vulhub github-auth [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--status` | Show current authentication status |
| `--remove` | Remove saved authentication |

### Behavior

**Default (authenticate):**
1. Checks if already authenticated
2. Initiates OAuth Device Flow
3. Displays authorization URL and code
4. Opens browser automatically
5. Polls for authorization completion
6. Saves access token to configuration

**With `--status`:**
- Shows whether authenticated
- Displays masked token if present
- Shows current rate limit tier

**With `--remove`:**
- Prompts for confirmation
- Removes saved token from configuration

### Example

```bash
# Authenticate
vulhub github-auth

# Check status
vulhub github-auth --status

# Remove authentication
vulhub github-auth --remove
```

---

## Keyword Resolution

Many commands accept a `<keyword>` argument. The resolution process works as follows:

### Match Priority

1. **Exact CVE Match**: If keyword matches a CVE exactly (e.g., `CVE-2021-44228`)
2. **Exact Path Match**: If keyword matches an environment path exactly (e.g., `log4j/CVE-2021-44228`)
3. **Fuzzy Match**: Searches across CVE, app name, path, and description

### Multiple Matches

When multiple environments match a keyword:
- **Interactive mode** (default): Displays a selection menu
- **Non-interactive mode** (`--yes`): Returns an error

### Examples

```bash
# Exact CVE - single match
vulhub start CVE-2021-44228

# Exact path - single match
vulhub start log4j/CVE-2021-44228

# Fuzzy match - may have multiple matches
vulhub start log4j  # Prompts if multiple log4j vulnerabilities exist
```

---

## Error Handling

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| "Docker is not installed" | Docker not found in PATH | Install Docker |
| "Docker daemon is not running" | Docker service not started | Start Docker Desktop or service |
| "vulhub-cli is not initialized" | `init` not run | Run `vulhub init` |
| "no environment found matching" | Keyword didn't match | Use `vulhub search` to find correct name |
| "rate limit exceeded" | GitHub API limit reached | Run `vulhub github-auth` |

### Rate Limit Handling

When a GitHub API rate limit error occurs:
1. CLI detects the rate limit error
2. Prompts user to authenticate with GitHub
3. If user authenticates, automatically retries the operation
4. No need to re-run the command manually
