# Vulhub CLI

A command-line tool for managing [Vulhub](https://github.com/vulhub/vulhub) vulnerability environments. Start, stop, and manage security lab environments without learning Docker Compose commands or downloading the complete Vulhub repository.

> **Warning**
> This project is currently in **experimental stage** and has not been officially released. The command-line interface, configuration file format, and other aspects may still change in future versions. Please be aware of potential breaking changes when using or upgrading this tool.

## Features

- **Easy Environment Management** - Start vulnerability labs with a single command
- **Smart Search** - Find environments by CVE number, application name, or fuzzy matching
- **Automatic Downloads** - Environments are downloaded on-demand from GitHub
- **GitHub Authentication** - Built-in OAuth flow to avoid API rate limits
- **Cross-Platform** - Works on Windows, macOS, and Linux

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) with Docker Compose V2
- [Go 1.21+](https://golang.org/dl/) (for building from source)

## Installation

### From Source

```bash
git clone https://github.com/vulhub/vulhub-cli.git
cd vulhub-cli
go build -o vulhub ./cmd/vulhub
```

### Add to PATH (Optional)

```bash
# Linux/macOS
sudo mv vulhub /usr/local/bin/

# Windows (PowerShell as Administrator)
Move-Item vulhub.exe C:\Windows\System32\
```

## Quick Start

```bash
# Initialize vulhub-cli (downloads environment list)
vulhub init

# Check system environment (Docker, network, etc.)
vulhub doctor

# Search for environments
vulhub search log4j

# Start an environment
vulhub start CVE-2021-44228

# Check environment status
vulhub status

# Stop an environment
vulhub stop CVE-2021-44228

# Completely remove an environment
vulhub clean CVE-2021-44228
```

## Commands Overview

| Command | Description |
|---------|-------------|
| `init` | Initialize configuration and download environment list |
| `syncup` | Update environment list from GitHub |
| `start` | Start a vulnerability environment |
| `stop` | Stop a running environment |
| `restart` | Restart an environment |
| `clean` | Completely remove an environment |
| `status` | Show status of downloaded environments (aliases: `ls`, `list`) |
| `list-available` | List all available environments |
| `search` | Search for environments |
| `info` | Show environment details |
| `github-auth` | Authenticate with GitHub |
| `doctor` | Check system environment and diagnose potential issues |

For detailed command documentation, see the [docs](./docs/) directory.

## Keyword Matching

Commands that accept a `[keyword]` argument support multiple matching modes:

- **Exact CVE**: `CVE-2021-44228`
- **Exact Path**: `log4j/CVE-2021-44228`
- **Fuzzy Match**: `log4j`, `struts`, `spring`

When multiple environments match, you'll be prompted to select one interactively.

## GitHub Rate Limits

GitHub API has rate limits (60 requests/hour for unauthenticated users). If you hit rate limits, the CLI will automatically prompt you to authenticate:

```bash
# Or authenticate proactively
vulhub github-auth
```

This increases the limit to 5,000 requests/hour.

## Configuration

Configuration is stored in `~/.vulhub/`:

```
~/.vulhub/
├── config.toml          # Main configuration
├── environments.toml    # Cached environment list
└── environments/        # Downloaded environment files
```

## Documentation

- [Getting Started](./docs/01-getting-started.md)
- [Command Reference](./docs/02-commands.md)
- [Configuration](./docs/03-configuration.md)
- [GitHub Authentication](./docs/04-authentication.md)
- [Architecture](./docs/05-architecture.md)
- [Troubleshooting](./docs/06-troubleshooting.md)

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Vulhub](https://github.com/vulhub/vulhub) - Pre-built vulnerable docker environments
