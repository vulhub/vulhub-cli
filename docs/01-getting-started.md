# Getting Started

This guide will help you get started with vulhub-cli, from installation to running your first vulnerability environment.

## Prerequisites

Before using vulhub-cli, ensure you have:

1. **Docker** installed and running
   - Download from [docker.com](https://docs.docker.com/get-docker/)
   - Docker Compose V2 is required (included with Docker Desktop)

2. **Go 1.21+** (only for building from source)
   - Download from [golang.org](https://golang.org/dl/)

## Installation

### Building from Source

```bash
# Clone the repository
git clone https://github.com/vulhub/vulhub-cli.git
cd vulhub-cli

# Build the binary
go build -o vulhub ./cmd/vulhub

# Verify installation
./vulhub --version
```

### Adding to System PATH

To use `vulhub` from any directory:

**Linux/macOS:**
```bash
sudo mv vulhub /usr/local/bin/
```

**Windows (PowerShell as Administrator):**
```powershell
Move-Item vulhub.exe C:\Windows\System32\
```

## First-Time Setup

### Step 1: Initialize vulhub-cli

Run the `init` command to create the configuration directory and download the environment list:

```bash
vulhub init
```

This will:
1. Create the configuration directory (`~/.vulhub/`)
2. Download the list of available environments from GitHub
3. Create the default configuration file

**Expected output:**
```
ℹ Creating default configuration...
ℹ Downloading environment list from GitHub...
✓ Initialization complete! Found 180 environments.

Configuration directory: /home/user/.vulhub/

Quick start:
  vulhub search log4j         # Search for environments
  vulhub info CVE-2021-44228  # View environment details
  vulhub start CVE-2021-44228 # Start an environment
```

### Step 2: Search for an Environment

Find vulnerability environments using keywords:

```bash
# Search by application name
vulhub search log4j

# Search by CVE number
vulhub search CVE-2021

# Search by technology
vulhub search spring
```

### Step 3: View Environment Details

Get detailed information about an environment:

```bash
vulhub info CVE-2021-44228
```

This displays:
- CVE information
- Application and version affected
- README content with exploitation instructions
- Docker Compose configuration

### Step 4: Start an Environment

Start the vulnerability environment:

```bash
vulhub start CVE-2021-44228
```

This will:
1. Download the environment files (if not already downloaded)
2. Pull the required Docker images
3. Start the containers
4. Display the exposed ports

**Example output:**
```
ℹ Starting environment: log4j/CVE-2021-44228

Containers:
  NAME                  STATUS    PORTS
  log4j-cve-2021-44228  running   0.0.0.0:8983->8983/tcp

✓ Environment 'log4j/CVE-2021-44228' started successfully!
```

### Step 5: Access the Environment

Open your browser and navigate to the exposed port:
```
http://localhost:8983
```

### Step 6: Stop the Environment

When you're done:

```bash
# Stop the environment (containers remain)
vulhub stop CVE-2021-44228

# Or completely remove it (containers, volumes, and files)
vulhub clean CVE-2021-44228
```

## Common Workflows

### Workflow 1: Quick Lab Setup

```bash
vulhub start struts2/s2-045
# ... do your testing ...
vulhub clean struts2/s2-045
```

### Workflow 2: Multiple Environments

```bash
# Start multiple environments
vulhub start CVE-2021-44228
vulhub start CVE-2017-5638

# Check what's running
vulhub status

# Stop all when done
vulhub stop CVE-2021-44228
vulhub stop CVE-2017-5638
```

### Workflow 3: Exploring Available Environments

```bash
# List all available environments
vulhub list-available

# Search for specific technology
vulhub search wordpress

# Get details before starting
vulhub info wordpress/pwnscriptum
```

## Troubleshooting

### Docker Not Running

If you see "Docker daemon is not running":

```bash
# Linux
sudo systemctl start docker

# macOS/Windows
# Open Docker Desktop application
```

### Rate Limit Exceeded

If you see "GitHub API rate limit exceeded":

```bash
vulhub github-auth
```

This authenticates with GitHub and increases your rate limit from 60 to 5,000 requests/hour.

### Environment Not Found

If a keyword doesn't match:

```bash
# Use search to find the correct name
vulhub search <partial-name>

# Use the exact path shown in search results
vulhub start app-name/CVE-XXXX-XXXXX
```

## Next Steps

- Read the [Command Reference](./commands.md) for detailed command documentation
- Learn about [Configuration](./configuration.md) options
- Set up [GitHub Authentication](./authentication.md) to avoid rate limits
