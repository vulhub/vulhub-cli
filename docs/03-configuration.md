# Configuration

This document describes the configuration options and file structure for vulhub-cli.

## Configuration Directory

By default, vulhub-cli stores all configuration and data in `~/.vulhub/`:

```
~/.vulhub/
├── config.toml          # Main configuration file
├── environments.toml    # Cached list of available environments
└── environments/        # Downloaded environment files
    ├── log4j/
    │   └── CVE-2021-44228/
    │       ├── docker-compose.yml
    │       └── ...
    └── struts2/
        └── s2-045/
            └── ...
```

### Platform-Specific Paths

| Platform | Default Path |
|----------|--------------|
| Linux | `~/.vulhub/` |
| macOS | `~/.vulhub/` |
| Windows | `%USERPROFILE%\.vulhub\` |

## Configuration File (config.toml)

The main configuration file uses TOML format:

```toml
# GitHub repository settings
[github]
owner = "vulhub"
repo = "vulhub"
branch = "master"
token = ""  # GitHub access token (set by github-auth)

# Sync settings
[sync]
last_sync = "2024-01-15T10:30:00Z"
auto_sync_days = 7  # Prompt to sync after this many days
```

### Configuration Options

#### [github] Section

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `owner` | string | `"vulhub"` | GitHub repository owner |
| `repo` | string | `"vulhub"` | GitHub repository name |
| `branch` | string | `"master"` | Git branch to use |
| `token` | string | `""` | GitHub access token |

#### [sync] Section

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `last_sync` | datetime | - | Timestamp of last sync |
| `auto_sync_days` | integer | `7` | Days before prompting to sync |

## Environment Variables

Environment variables take precedence over configuration file settings:

| Variable | Description |
|----------|-------------|
| `GITHUB_TOKEN` | GitHub access token (overrides config.toml) |

### Example

```bash
# Use a specific GitHub token
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
vulhub start log4j
```

## Environments List (environments.toml)

This file contains the cached list of available vulnerability environments. It is automatically downloaded from the Vulhub repository during `init` and `syncup`.

### Structure

```toml
[[environment]]
path = "log4j/CVE-2021-44228"
app = "Apache Log4j"
cve = ["CVE-2021-44228"]
description = "Apache Log4j2 Remote Code Execution"

[[environment]]
path = "struts2/s2-045"
app = "Apache Struts2"
cve = ["CVE-2017-5638"]
description = "Apache Struts2 Remote Code Execution"

# ... more environments
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `path` | string | Environment directory path |
| `app` | string | Application name |
| `cve` | array | List of CVE numbers |
| `description` | string | Brief description |

## Downloaded Environments

When you start an environment for the first time, vulhub-cli downloads the necessary files from GitHub and stores them locally:

```
~/.vulhub/environments/<app>/<vulnerability>/
├── docker-compose.yml    # Docker Compose configuration
├── README.md             # Documentation (if available)
├── README.zh-cn.md       # Chinese documentation (if available)
└── ...                   # Other files (Dockerfiles, configs, etc.)
```

### Managing Downloaded Environments

**List downloaded environments:**
```bash
vulhub list
```

**Remove a specific environment:**
```bash
vulhub clean <keyword>
```

**Manual cleanup:**
```bash
rm -rf ~/.vulhub/environments/<app>/<vulnerability>
```

## Custom Configuration Path

You can specify a custom configuration file location using the `--config` flag:

```bash
vulhub --config /path/to/custom/config.toml start log4j
```

This only changes the config file location; the environments directory remains relative to the config file location.

## Resetting Configuration

To completely reset vulhub-cli:

```bash
# Remove all configuration and downloaded environments
rm -rf ~/.vulhub

# Re-initialize
vulhub init
```

## Troubleshooting

### Configuration Not Found

If vulhub-cli reports "not initialized":

```bash
vulhub init
```

### Token Not Working

1. Check if environment variable is set:
   ```bash
   echo $GITHUB_TOKEN
   ```

2. Re-authenticate:
   ```bash
   vulhub github-auth --remove
   vulhub github-auth
   ```

### Sync Issues

Force a fresh sync:

```bash
vulhub syncup
```

### Corrupted Configuration

Reset and reinitialize:

```bash
rm ~/.vulhub/config.toml
vulhub init --force
```
