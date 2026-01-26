package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v3"

	"github.com/vulhub/vulhub-cli/internal/cli/ui"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// checkResult represents the result of a single check
type checkResult struct {
	Name    string
	Status  checkStatus
	Message string
	Hint    string
}

type checkStatus int

const (
	statusOK checkStatus = iota
	statusWarning
	statusError
	statusSkipped
)

// Doctor creates the doctor command
func (c *Commands) Doctor() *cli.Command {
	return &cli.Command{
		Name:    "doctor",
		Usage:   "Check system environment and diagnose potential issues",
		Aliases: []string{"doc"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Show detailed information for each check",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return c.runDoctor(ctx, cmd.Bool("verbose"))
		},
	}
}

func (c *Commands) runDoctor(ctx context.Context, verbose bool) error {
	table := ui.NewTable()

	fmt.Println("Vulhub Doctor - System Environment Check")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	var results []checkResult
	var hasErrors bool
	var hasWarnings bool

	// 1. Docker checks
	fmt.Println("Docker Environment")
	fmt.Println(strings.Repeat("-", 30))

	dockerResults := c.checkDocker(ctx, verbose)
	results = append(results, dockerResults...)
	c.printCheckResults(dockerResults)
	fmt.Println()

	// 2. Configuration checks
	fmt.Println("Configuration")
	fmt.Println(strings.Repeat("-", 30))

	configResults := c.checkConfiguration(ctx, verbose)
	results = append(results, configResults...)
	c.printCheckResults(configResults)
	fmt.Println()

	// 3. Network connectivity checks
	fmt.Println("Network Connectivity")
	fmt.Println(strings.Repeat("-", 30))

	networkResults := c.checkNetworkConnectivity(ctx, verbose)
	results = append(results, networkResults...)
	c.printCheckResults(networkResults)
	fmt.Println()

	// 4. Docker registry checks
	fmt.Println("Docker Registry")
	fmt.Println(strings.Repeat("-", 30))

	registryResults := c.checkDockerRegistry(ctx, verbose)
	results = append(results, registryResults...)
	c.printCheckResults(registryResults)
	fmt.Println()

	// Summary
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Summary")
	fmt.Println(strings.Repeat("-", 30))

	okCount := 0
	warnCount := 0
	errCount := 0
	skipCount := 0

	for _, r := range results {
		switch r.Status {
		case statusOK:
			okCount++
		case statusWarning:
			warnCount++
			hasWarnings = true
		case statusError:
			errCount++
			hasErrors = true
		case statusSkipped:
			skipCount++
		}
	}

	fmt.Printf("  %s %d passed\n", ui.SuccessStyle.Render("[OK]"), okCount)
	if warnCount > 0 {
		fmt.Printf("  %s %d warnings\n", ui.WarningStyle.Render("[WARN]"), warnCount)
	}
	if errCount > 0 {
		fmt.Printf("  %s %d errors\n", ui.ErrorStyle.Render("[ERR]"), errCount)
	}
	if skipCount > 0 {
		fmt.Printf("  %s %d skipped\n", ui.MutedStyle.Render("[SKIP]"), skipCount)
	}
	fmt.Println()

	// Print hints for issues
	if hasErrors || hasWarnings {
		fmt.Println("Suggestions:")
		fmt.Println(strings.Repeat("-", 30))
		hintSet := make(map[string]bool) // Deduplicate hints
		for _, r := range results {
			if r.Hint != "" && (r.Status == statusError || r.Status == statusWarning) {
				if !hintSet[r.Hint] {
					fmt.Printf("  - %s: %s\n", r.Name, r.Hint)
					hintSet[r.Hint] = true
				}
			}
		}
		fmt.Println()
	}

	// Network troubleshooting suggestions
	if c.hasNetworkIssues(results) {
		c.printNetworkTroubleshooting()
	}

	if hasErrors {
		return fmt.Errorf("some checks failed, please fix the issues above")
	}

	table.PrintSuccess("All critical checks passed! Vulhub is ready to use.")
	return nil
}

func (c *Commands) printCheckResults(results []checkResult) {
	for _, r := range results {
		var statusStr string
		switch r.Status {
		case statusOK:
			statusStr = ui.SuccessStyle.Render("[OK]")
		case statusWarning:
			statusStr = ui.WarningStyle.Render("[WARN]")
		case statusError:
			statusStr = ui.ErrorStyle.Render("[ERR]")
		case statusSkipped:
			statusStr = ui.MutedStyle.Render("[SKIP]")
		}
		fmt.Printf("  %s %s: %s\n", statusStr, r.Name, r.Message)
	}
}

// checkDocker checks Docker-related requirements
func (c *Commands) checkDocker(ctx context.Context, verbose bool) []checkResult {
	var results []checkResult

	// Check if docker command exists
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		results = append(results, checkResult{
			Name:    "Docker Installation",
			Status:  statusError,
			Message: "Docker is not installed",
			Hint:    "Install Docker from https://docs.docker.com/get-docker/",
		})
		// Skip remaining Docker checks
		results = append(results, checkResult{
			Name:    "Docker Daemon",
			Status:  statusSkipped,
			Message: "Skipped (Docker not installed)",
		})
		results = append(results, checkResult{
			Name:    "Docker Compose",
			Status:  statusSkipped,
			Message: "Skipped (Docker not installed)",
		})
		return results
	}

	msg := "Docker is installed"
	if verbose {
		msg = fmt.Sprintf("Docker is installed at %s", dockerPath)
	}
	results = append(results, checkResult{
		Name:    "Docker Installation",
		Status:  statusOK,
		Message: msg,
	})

	// Check if docker daemon is running
	cmd := exec.CommandContext(ctx, dockerPath, "info")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		hint := "Start Docker daemon"
		switch runtime.GOOS {
		case "windows", "darwin":
			hint = "Start Docker Desktop from the application menu"
		default:
			hint = "Run: sudo systemctl start docker"
		}
		results = append(results, checkResult{
			Name:    "Docker Daemon",
			Status:  statusError,
			Message: "Docker daemon is not running",
			Hint:    hint,
		})
		results = append(results, checkResult{
			Name:    "Docker Compose",
			Status:  statusSkipped,
			Message: "Skipped (Docker daemon not running)",
		})
		return results
	}

	results = append(results, checkResult{
		Name:    "Docker Daemon",
		Status:  statusOK,
		Message: "Docker daemon is running",
	})

	// Check docker compose
	cmd = exec.CommandContext(ctx, dockerPath, "compose", "version")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		results = append(results, checkResult{
			Name:    "Docker Compose",
			Status:  statusError,
			Message: "Docker Compose is not available",
			Hint:    "Please upgrade Docker to the latest version or install Docker Compose plugin: https://docs.docker.com/compose/install/",
		})
	} else {
		msg := "Docker Compose is available"
		if verbose {
			version := strings.TrimSpace(stdout.String())
			if version != "" {
				msg = version
			}
		}
		results = append(results, checkResult{
			Name:    "Docker Compose",
			Status:  statusOK,
			Message: msg,
		})
	}

	return results
}

// checkConfiguration checks vulhub configuration files
func (c *Commands) checkConfiguration(ctx context.Context, verbose bool) []checkResult {
	var results []checkResult
	paths := c.Config.Paths()

	// Check config directory
	configDir := paths.ConfigDir()
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		results = append(results, checkResult{
			Name:    "Config Directory",
			Status:  statusWarning,
			Message: fmt.Sprintf("%s does not exist", configDir),
			Hint:    "Run 'vulhub init' to initialize",
		})
	} else {
		results = append(results, checkResult{
			Name:    "Config Directory",
			Status:  statusOK,
			Message: fmt.Sprintf("%s exists", configDir),
		})
	}

	// Check config.toml
	configFile := paths.ConfigFile()
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		results = append(results, checkResult{
			Name:    "Config File",
			Status:  statusWarning,
			Message: "config.toml not found",
			Hint:    "Run 'vulhub init' to create configuration",
		})
	} else {
		// Try to parse config file
		var cfg types.Config
		if _, err := toml.DecodeFile(configFile, &cfg); err != nil {
			results = append(results, checkResult{
				Name:    "Config File",
				Status:  statusError,
				Message: fmt.Sprintf("config.toml is corrupted: %v", err),
				Hint:    "Delete the file and run 'vulhub init' to recreate",
			})
		} else {
			results = append(results, checkResult{
				Name:    "Config File",
				Status:  statusOK,
				Message: "config.toml is valid",
			})
		}
	}

	// Check environments.toml
	envsFile := paths.EnvironmentsFile()
	if _, err := os.Stat(envsFile); os.IsNotExist(err) {
		results = append(results, checkResult{
			Name:    "Environments File",
			Status:  statusWarning,
			Message: "environments.toml not found",
			Hint:    "Run 'vulhub init' or 'vulhub syncup' to download environment list",
		})
	} else {
		// Try to parse environments file
		var envList types.EnvironmentList
		if _, err := toml.DecodeFile(envsFile, &envList); err != nil {
			results = append(results, checkResult{
				Name:    "Environments File",
				Status:  statusError,
				Message: fmt.Sprintf("environments.toml is corrupted: %v", err),
				Hint:    "Delete the file and run 'vulhub syncup' to re-download",
			})
		} else {
			results = append(results, checkResult{
				Name:    "Environments File",
				Status:  statusOK,
				Message: fmt.Sprintf("environments.toml is valid (%d environments)", len(envList.Environment)),
			})
		}
	}

	// Check environments directory
	envsDir := paths.EnvironmentsDir()
	if info, err := os.Stat(envsDir); os.IsNotExist(err) {
		results = append(results, checkResult{
			Name:    "Environments Directory",
			Status:  statusOK,
			Message: "environments/ will be created when needed",
		})
	} else if !info.IsDir() {
		results = append(results, checkResult{
			Name:    "Environments Directory",
			Status:  statusError,
			Message: "environments exists but is not a directory",
			Hint:    "Remove the file and run 'vulhub init'",
		})
	} else {
		downloaded, err := c.Environment.ListDownloaded(ctx)
		if err != nil {
			results = append(results, checkResult{
				Name:    "Environments Directory",
				Status:  statusError,
				Message: fmt.Sprintf("Failed to list downloaded environments: %v", err),
			})
		} else {
			results = append(results, checkResult{
				Name:    "Environments Directory",
				Status:  statusOK,
				Message: fmt.Sprintf("%d environments downloaded", len(downloaded)),
			})
		}
	}

	return results
}

// checkNetworkConnectivity checks network connectivity to required services
func (c *Commands) checkNetworkConnectivity(ctx context.Context, verbose bool) []checkResult {
	var results []checkResult

	// Define endpoints to check
	endpoints := []struct {
		name     string
		url      string
		host     string
		critical bool
		hint     string
	}{
		// GitHub endpoints (required for downloading environment configs)
		{
			name:     "GitHub API",
			url:      "https://api.github.com",
			host:     "api.github.com",
			critical: true,
			hint:     "Required for fetching environment list. Try using a proxy or VPN.",
		},
		{
			name:     "GitHub Raw Content",
			url:      "https://raw.githubusercontent.com",
			host:     "raw.githubusercontent.com",
			critical: true,
			hint:     "Required for downloading environment files. Try using a proxy or VPN.",
		},
		// Docker Hub endpoints (required for pulling images)
		{
			name:     "Docker Hub",
			url:      "https://hub.docker.com",
			host:     "hub.docker.com",
			critical: false,
			hint:     "Docker Hub website. Consider using a registry mirror.",
		},
		{
			name:     "Docker Registry",
			url:      "https://registry-1.docker.io/v2/",
			host:     "registry-1.docker.io",
			critical: true,
			hint:     "Docker image registry. Consider using a registry mirror.",
		},
		{
			name:     "Docker Auth",
			url:      "https://auth.docker.io",
			host:     "auth.docker.io",
			critical: true,
			hint:     "Docker authentication service. Consider using a registry mirror.",
		},
		{
			name:     "Docker CDN",
			url:      "https://production.cloudflare.docker.com",
			host:     "production.cloudflare.docker.com",
			critical: false,
			hint:     "Docker image layer CDN. Consider using a registry mirror.",
		},
	}

	// Use the configured HTTP client (with proxy if set)
	for _, ep := range endpoints {
		// First check DNS resolution
		start := time.Now()
		_, err := net.DefaultResolver.LookupHost(ctx, ep.host)
		dnsTime := time.Since(start)

		if err != nil {
			status := statusWarning
			if ep.critical {
				status = statusError
			}
			results = append(results, checkResult{
				Name:    ep.name,
				Status:  status,
				Message: fmt.Sprintf("DNS resolution failed: %v", err),
				Hint:    ep.hint,
			})
			continue
		}

		// Then check HTTP connectivity
		start = time.Now()
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, ep.url, nil)
		resp, err := c.HTTPClient.Do(req)
		httpTime := time.Since(start)

		if err != nil {
			status := statusWarning
			if ep.critical {
				status = statusError
			}
			results = append(results, checkResult{
				Name:    ep.name,
				Status:  status,
				Message: fmt.Sprintf("Connection failed: %v", err),
				Hint:    ep.hint,
			})
			continue
		}
		resp.Body.Close()

		// Check response time
		totalTime := dnsTime + httpTime
		msg := "OK"
		status := statusOK

		if verbose {
			msg = fmt.Sprintf("OK (DNS: %dms, HTTP: %dms)", dnsTime.Milliseconds(), httpTime.Milliseconds())
		}

		if totalTime > 5*time.Second {
			status = statusWarning
			msg = fmt.Sprintf("Slow connection (%dms)", totalTime.Milliseconds())
		}

		results = append(results, checkResult{
			Name:    ep.name,
			Status:  status,
			Message: msg,
		})
	}

	return results
}

// checkDockerRegistry checks Docker registry connectivity
func (c *Commands) checkDockerRegistry(ctx context.Context, verbose bool) []checkResult {
	var results []checkResult

	// Check registry connectivity via Docker
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		results = append(results, checkResult{
			Name:    "Registry Connectivity",
			Status:  statusSkipped,
			Message: "Skipped (Docker not installed)",
		})
		return results
	}

	// Check if daemon is running
	cmd := exec.CommandContext(ctx, dockerPath, "info")
	if err := cmd.Run(); err != nil {
		results = append(results, checkResult{
			Name:    "Docker Pull Test",
			Status:  statusSkipped,
			Message: "Skipped (Docker daemon not running)",
		})
		return results
	}

	// Check Docker daemon registry mirrors configuration first
	cmd = exec.CommandContext(ctx, dockerPath, "info", "--format", "{{range .RegistryConfig.Mirrors}}{{.}} {{end}}")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Run()

	mirrors := strings.TrimSpace(stdout.String())
	if mirrors != "" {
		results = append(results, checkResult{
			Name:    "Registry Mirrors",
			Status:  statusOK,
			Message: fmt.Sprintf("Configured: %s", mirrors),
		})
	} else {
		results = append(results, checkResult{
			Name:    "Registry Mirrors",
			Status:  statusOK,
			Message: "No mirrors configured (using Docker Hub directly)",
		})
	}

	// Remove local hello-world image first to ensure a fresh pull test
	// This ensures we actually test network connectivity, not local cache
	testImage := "hello-world:latest"
	removeCmd := exec.CommandContext(ctx, dockerPath, "rmi", "-f", testImage)
	removeCmd.Run() // Ignore error - image might not exist

	// Try to pull a minimal test image
	// Use a context with timeout for the pull operation
	pullCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cmd = exec.CommandContext(pullCtx, dockerPath, "pull", testImage)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Run()
	pullTime := time.Since(start)

	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if strings.Contains(errMsg, "timeout") || errors.Is(pullCtx.Err(), context.DeadlineExceeded) {
			results = append(results, checkResult{
				Name:    "Docker Pull Test",
				Status:  statusError,
				Message: "Image pull timed out (>60s)",
				Hint:    "Docker registry may be slow or unreachable. Consider using a registry mirror.",
			})
		} else {
			results = append(results, checkResult{
				Name:    "Docker Pull Test",
				Status:  statusError,
				Message: fmt.Sprintf("Failed to pull test image: %s", errMsg),
				Hint:    "Check Docker daemon logs and network connectivity.",
			})
		}
	} else {
		msg := "Successfully pulled test image"
		status := statusOK

		if verbose {
			msg = fmt.Sprintf("Successfully pulled %s in %dms", testImage, pullTime.Milliseconds())
		}

		if pullTime > 30*time.Second {
			status = statusWarning
			msg = fmt.Sprintf("Slow pull (%ds). Consider using a registry mirror.", int(pullTime.Seconds()))
		}

		results = append(results, checkResult{
			Name:    "Docker Pull Test",
			Status:  status,
			Message: msg,
		})

		// Clean up the test image to avoid leaving artifacts
		cleanupCmd := exec.CommandContext(ctx, dockerPath, "rmi", "-f", testImage)
		cleanupCmd.Run() // Ignore error
	}

	return results
}

// hasNetworkIssues checks if results indicate network connectivity issues
func (c *Commands) hasNetworkIssues(results []checkResult) bool {
	for _, r := range results {
		if r.Status == statusError || r.Status == statusWarning {
			name := strings.ToLower(r.Name)
			msg := strings.ToLower(r.Message)
			// Check if it's a network-related issue
			if strings.Contains(name, "github") || strings.Contains(name, "docker") {
				if strings.Contains(msg, "timeout") ||
					strings.Contains(msg, "connection") ||
					strings.Contains(msg, "slow") ||
					strings.Contains(msg, "dns") ||
					strings.Contains(msg, "failed") {
					return true
				}
			}
		}
	}
	return false
}

// printNetworkTroubleshooting prints general network troubleshooting suggestions
func (c *Commands) printNetworkTroubleshooting() {
	fmt.Println()
	fmt.Println(ui.WarningStyle.Render("Network Troubleshooting"))
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println()

	fmt.Println("1. Configure Docker Registry Mirror (for slow Docker pulls):")
	fmt.Println("   Edit Docker daemon configuration:")
	switch runtime.GOOS {
	case "windows":
		fmt.Println("   - Docker Desktop: Settings -> Docker Engine")
	case "darwin":
		fmt.Println("   - Docker Desktop: Preferences -> Docker Engine")
	default:
		fmt.Println("   - Edit /etc/docker/daemon.json")
	}
	fmt.Println()
	fmt.Println("   Add registry mirrors:")
	fmt.Println("   {")
	fmt.Println("     \"registry-mirrors\": [")
	fmt.Println("       \"https://mirror.gcr.io\",")
	fmt.Println("       \"https://docker-cf.registry.cyou\"")
	fmt.Println("     ]")
	fmt.Println("   }")
	fmt.Println()
	fmt.Println("   Then restart Docker daemon.")
	fmt.Println()

	fmt.Println("2. Use a Proxy (for GitHub/Docker connectivity issues):")
	fmt.Println("   Set HTTP_PROXY and HTTPS_PROXY environment variables:")
	fmt.Println()
	switch runtime.GOOS {
	case "windows":
		fmt.Println("   set HTTP_PROXY=http://your-proxy:port")
		fmt.Println("   set HTTPS_PROXY=http://your-proxy:port")
	default:
		fmt.Println("   export HTTP_PROXY=http://your-proxy:port")
		fmt.Println("   export HTTPS_PROXY=http://your-proxy:port")
	}
	fmt.Println()
	fmt.Println("   For Docker, also configure proxy in Docker daemon settings.")
	fmt.Println()

	fmt.Println("3. Check DNS Settings:")
	fmt.Println("   Try using public DNS servers like:")
	fmt.Println("   - Google: 8.8.8.8, 8.8.4.4")
	fmt.Println("   - Cloudflare: 1.1.1.1, 1.0.0.1")
	fmt.Println()

	fmt.Println("4. Authenticate with GitHub (for API rate limits):")
	fmt.Println("   Run 'vulhub github-auth' to increase API rate limit.")
	fmt.Println()
}
