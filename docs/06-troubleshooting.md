# Troubleshooting

This document provides solutions for common issues detected by the `vulhub doctor` command and other problems you may encounter while using vulhub-cli.

## Using the Doctor Command

The `vulhub doctor` command performs comprehensive system checks to diagnose potential issues:

```bash
vulhub doctor
```

### Options

| Option | Description |
|--------|-------------|
| `--verbose, -v` | Show detailed information for each check |

### Example Output

```bash
$ vulhub doctor

Vulhub Doctor - System Environment Check
==================================================

Docker Environment
------------------------------
  [OK] Docker Installation: Docker is installed
  [OK] Docker Daemon: Docker daemon is running
  [OK] Docker Compose: Docker Compose is available

Configuration
------------------------------
  [OK] Config Directory: C:\Users\anywhere\.vulhub exists
  [OK] Config File: config.toml is valid
  [OK] Environments File: environments.toml is valid (307 environments)
  [OK] Environments Directory: 0 environments downloaded

Network Connectivity
------------------------------
  [OK] GitHub API: OK
  [OK] GitHub Raw Content: OK
  [OK] Docker Hub: OK
  [OK] Docker Registry: OK
  [OK] Docker Auth: OK
  [OK] Docker CDN: OK

Docker Registry
------------------------------
  [OK] Registry Mirrors: No mirrors configured (using Docker Hub directly)
  [OK] Docker Pull Test: Successfully pulled test image

==================================================
Summary
------------------------------
  [OK] 15 passed

[OK] All critical checks passed! Vulhub is ready to use.
```

---

## Docker Issues

### Docker Not Installed

**Symptom**: `[ERR] Docker Installation: Docker is not installed`

**Solution**: Install Docker for your operating system:

- **Windows/macOS**: Download and install [Docker Desktop](https://docs.docker.com/get-docker/)
- **Linux**: Follow the official installation guide for your distribution:

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install docker.io docker-compose-plugin

# CentOS/RHEL/Fedora
sudo dnf install docker docker-compose-plugin

# Start and enable Docker
sudo systemctl enable docker
sudo systemctl start docker
```

---

### Docker Daemon Not Running

**Symptom**: `[ERR] Docker Daemon: Docker daemon is not running`

**Solution**:

**Windows/macOS (Docker Desktop)**:
1. Open Docker Desktop from the application menu
2. Wait for the Docker icon in the system tray to show "Docker is running"

**Linux**:
```bash
# Start Docker daemon
sudo systemctl start docker

# Check status
sudo systemctl status docker

# Enable auto-start on boot
sudo systemctl enable docker
```

**Common causes**:
- Docker Desktop not started after system reboot
- Insufficient system resources (memory/disk)
- Permission issues (Linux: ensure your user is in the `docker` group)

```bash
# Add user to docker group (Linux)
sudo usermod -aG docker $USER
# Log out and log back in for changes to take effect
```

---

### Docker Compose Not Available

**Symptom**: `[ERR] Docker Compose: Docker Compose is not available`

**Solution**:

Docker Compose v2 is included with all modern Docker installations (Docker Desktop and Docker Engine). If you see this error:

1. **Update Docker to the latest version**:
   - **Windows/macOS**: Update Docker Desktop from the application menu
   - **Linux**: Follow the [official Docker installation guide](https://docs.docker.com/engine/install/) to install or update

2. **Verify installation**:
   ```bash
   docker compose version
   ```

3. **If using legacy `docker-compose` (v1)**: The standalone `docker-compose` command is deprecated. Upgrade Docker to get the built-in `docker compose` (with space) command that vulhub-cli uses.

---

## Configuration Issues

### Config Directory Does Not Exist

**Symptom**: `[WARN] Config Directory: /home/user/.vulhub does not exist`

**Solution**:

```bash
vulhub init
```

---

### Config File Not Found or Corrupted

**Symptom**: 
- `[WARN] Config File: config.toml not found`
- `[ERR] Config File: config.toml is corrupted: <error>`

**Solution**:

```bash
# Re-initialize configuration
vulhub init --force
```

If you need to preserve custom settings, backup and recreate:

```bash
# Backup corrupted file
mv ~/.vulhub/config.toml ~/.vulhub/config.toml.bak

# Re-initialize
vulhub init
```

---

### Environments File Not Found or Corrupted

**Symptom**:
- `[WARN] Environments File: environments.toml not found`
- `[ERR] Environments File: environments.toml is corrupted: <error>`

**Solution**:

```bash
# Re-download environment list
vulhub syncup

# Or if syncup fails, delete and re-initialize
rm ~/.vulhub/environments.toml
vulhub init
```

---

### Environments Directory Issues

**Symptom**: `[ERR] Environments Directory: environments exists but is not a directory`

**Solution**:

```bash
# Remove the file and re-initialize
rm ~/.vulhub/environments
vulhub init
```

---

## Network Connectivity Issues

The doctor command checks connectivity to these critical services:

| Service | URL | Purpose |
|---------|-----|---------|
| GitHub API | `api.github.com` | Fetching environment list |
| GitHub Raw Content | `raw.githubusercontent.com` | Downloading environment files |
| Docker Hub | `hub.docker.com` | Docker Hub website |
| Docker Registry | `registry-1.docker.io` | Pulling Docker images |
| Docker Auth | `auth.docker.io` | Docker authentication |
| Docker CDN | `production.cloudflare.docker.com` | Docker image layers |

### DNS Resolution Failed

**Symptom**: `[ERR] <Service>: DNS resolution failed: <error>`

**Possible causes**:

1. **DNS server issues** - Your configured DNS server may be unreliable or unreachable
2. **DNS pollution** - In some regions (e.g., China), DNS queries may be intercepted or poisoned by network filtering (GFW), returning incorrect IP addresses for services like GitHub

**Solutions**:

1. **Use reliable public DNS servers** such as Google DNS (`8.8.8.8`) or Cloudflare DNS (`1.1.1.1`). The configuration method varies by operating system.

2. **If DNS pollution is suspected**, consider using:
   - Encrypted DNS (DoH/DoT)
   - A VPN or proxy service
   - Local DNS tools that bypass pollution (e.g., `dnscrypt-proxy`, `smartdns`)

---

### Connection Failed or Timeout

**Symptom**: 
- `[ERR] <Service>: Connection failed: <error>`
- `[WARN] <Service>: Slow connection (<N>ms)`

**Solutions**:

1. **Configure proxy for vulhub-cli:**

   Using `--proxy` flag:
   ```bash
   vulhub --proxy http://127.0.0.1:8080 syncup
   vulhub --proxy socks5://127.0.0.1:1080 start log4j
   ```

   Using environment variable:
   ```bash
   export VULHUB_PROXY=http://127.0.0.1:8080
   vulhub syncup
   ```

   Using config file (`~/.vulhub/config.toml`):
   ```toml
   [network]
   proxy = "http://127.0.0.1:8080"
   ```

   Note: This only affects vulhub-cli's network requests. Docker image pulls require separate proxy configuration (see [Docker Registry Issues](#docker-registry-issues)).

2. **Use a VPN** if services are blocked in your region

3. **Check firewall settings** - Ensure outbound HTTPS (port 443) is allowed

---

### GitHub API Rate Limiting

**Symptom**: GitHub API requests failing with rate limit errors

**Solution**: Authenticate with GitHub to increase rate limits:

```bash
vulhub github-auth
```

This uses OAuth Device Flow to authenticate and increases your API rate limit from 60 to 5,000 requests per hour.

---

## Docker Registry Issues

### Docker Pull Test Failed or Slow

**Symptom**:
- `[ERR] Docker Pull Test: Image pull timed out (>60s)`
- `[ERR] Docker Pull Test: Failed to pull test image: <error>`
- `[WARN] Docker Pull Test: Slow pull (<N>s). Consider using a registry mirror.`

**Possible causes**:

1. **Network issues** - Docker Hub may be slow or unreachable from your region
2. **Firewall/proxy blocking** - Corporate networks may block Docker registry traffic

**Solutions**:

1. **Configure registry mirrors** - Add a Docker Hub mirror to speed up image pulls. Edit Docker daemon configuration to add:

   ```json
   {
     "registry-mirrors": ["https://mirror.example.com"]
   }
   ```

   - **Docker Desktop**: Settings/Preferences -> Docker Engine
   - **Linux**: Edit `/etc/docker/daemon.json` and restart Docker

2. **Configure Docker proxy** - If you're behind a corporate proxy, configure Docker to use it through Docker Desktop settings or systemd service configuration

3. **Use a VPN** - If Docker Hub is blocked in your region

---

## Common Error Messages

### "vulhub is not initialized"

Run the init command first:

```bash
vulhub init
```

### "No matching environment found"

The environment keyword doesn't match any available environment:

```bash
# List all available environments
vulhub list-available

# Search for environments
vulhub search log4j
```

### "Multiple environments match"

Your keyword matches multiple environments. Use a more specific keyword:

```bash
# Use exact CVE number
vulhub start CVE-2021-44228

# Or use full path
vulhub start log4j/CVE-2021-44228
```

### "Environment not downloaded"

Download the environment first:

```bash
vulhub start <environment>  # This will download automatically
```

### Permission Denied (Linux)

If you see permission errors when running Docker commands:

```bash
# Add your user to the docker group
sudo usermod -aG docker $USER

# Log out and log back in, or run:
newgrp docker
```

---

## Getting Help

If you continue to experience issues:

1. Run `vulhub doctor -v` for detailed diagnostics
2. Check the [Vulhub GitHub Issues](https://github.com/vulhub/vulhub/issues)
3. Ensure Docker and Docker Compose are up to date
4. Try running `vulhub init --force` to reset configuration
