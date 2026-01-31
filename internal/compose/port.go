package compose

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	maxPortAttempts     = 100
	originalComposeFile = "docker-compose.yml"
)

// dockerContainer represents a container from docker ps --format json
type dockerContainer struct {
	Ports string `json:"Ports"`
}

// portChecker checks port availability using Docker API
type portChecker struct {
	usedPorts map[int]bool
}

// newPortChecker creates a new port checker by querying Docker for used ports
func newPortChecker() *portChecker {
	pc := &portChecker{
		usedPorts: make(map[int]bool),
	}
	pc.loadDockerPorts()
	return pc
}

// loadDockerPorts queries Docker for all ports used by running containers
func (pc *portChecker) loadDockerPorts() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "ps", "--format", "{{json .}}")
	output, err := cmd.Output()
	if err != nil {
		slog.Debug("failed to get docker containers", "error", err)
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		var container dockerContainer
		if err := json.Unmarshal([]byte(line), &container); err != nil {
			slog.Debug("failed to parse container info", "error", err)
			continue
		}

		ports := parseDockerPorts(container.Ports)
		for _, port := range ports {
			pc.usedPorts[port] = true
		}
	}

	slog.Debug("loaded docker ports", "count", len(pc.usedPorts))
}

// parseDockerPorts parses the Ports string from docker ps
func parseDockerPorts(portsStr string) []int {
	var ports []int
	if portsStr == "" {
		return ports
	}

	parts := strings.Split(portsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if idx := strings.Index(part, "->"); idx != -1 {
			hostPart := part[:idx]
			if colonIdx := strings.LastIndex(hostPart, ":"); colonIdx != -1 {
				portStr := hostPart[colonIdx+1:]
				if port, err := strconv.Atoi(portStr); err == nil {
					ports = append(ports, port)
				}
			}
		}
	}

	return ports
}

// isAvailable checks if a port is available
func (pc *portChecker) isAvailable(port int) bool {
	return !pc.usedPorts[port]
}

// markUsed marks a port as used
func (pc *portChecker) markUsed(port int) {
	pc.usedPorts[port] = true
}

// findAvailable finds an available port starting from startPort
func (pc *portChecker) findAvailable(startPort int) (int, error) {
	for i := 0; i < maxPortAttempts; i++ {
		port := startPort + i
		if port > 65535 {
			break
		}
		if pc.isAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found after %d attempts starting from %d", maxPortAttempts, startPort)
}

// resolvePortConflicts checks for port conflicts and modifies docker-compose.yml if needed
func resolvePortConflicts(workDir string) error {
	composePath := filepath.Join(workDir, originalComposeFile)

	// Read original compose file
	data, err := os.ReadFile(composePath)
	if err != nil {
		slog.Debug("failed to read compose file", "error", err)
		return nil
	}

	// Parse YAML while preserving structure
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		slog.Debug("failed to parse compose file", "error", err)
		return nil
	}

	// Create port checker
	checker := newPortChecker()

	// Find and modify port mappings
	hasChanges := false
	if root.Kind == yaml.DocumentNode && len(root.Content) > 0 {
		hasChanges = processServicesNode(root.Content[0], checker)
	}

	if !hasChanges {
		return nil
	}

	// Write back to original file
	output, err := yaml.Marshal(&root)
	if err != nil {
		slog.Warn("failed to marshal compose file", "error", err)
		return nil
	}

	if err := os.WriteFile(composePath, output, 0644); err != nil {
		slog.Warn("failed to write compose file", "error", err)
		return nil
	}

	slog.Info("updated compose file with resolved ports", "path", composePath)
	return nil
}

// processServicesNode finds the services node and processes port mappings
func processServicesNode(root *yaml.Node, checker *portChecker) bool {
	if root.Kind != yaml.MappingNode {
		return false
	}

	hasChanges := false

	for i := 0; i < len(root.Content)-1; i += 2 {
		keyNode := root.Content[i]
		valueNode := root.Content[i+1]

		if keyNode.Value == "services" && valueNode.Kind == yaml.MappingNode {
			for j := 0; j < len(valueNode.Content)-1; j += 2 {
				serviceName := valueNode.Content[j].Value
				serviceNode := valueNode.Content[j+1]
				if processServicePorts(serviceName, serviceNode, checker) {
					hasChanges = true
				}
			}
		}
	}

	return hasChanges
}

// processServicePorts processes the ports of a single service
func processServicePorts(serviceName string, serviceNode *yaml.Node, checker *portChecker) bool {
	if serviceNode.Kind != yaml.MappingNode {
		return false
	}

	hasChanges := false

	for i := 0; i < len(serviceNode.Content)-1; i += 2 {
		keyNode := serviceNode.Content[i]
		valueNode := serviceNode.Content[i+1]

		if keyNode.Value == "ports" && valueNode.Kind == yaml.SequenceNode {
			for _, portNode := range valueNode.Content {
				if processPortNode(serviceName, portNode, checker) {
					hasChanges = true
				}
			}
		}
	}

	return hasChanges
}

// processPortNode processes a single port mapping node
func processPortNode(serviceName string, portNode *yaml.Node, checker *portChecker) bool {
	switch portNode.Kind {
	case yaml.ScalarNode:
		return processShortSyntaxPort(serviceName, portNode, checker)
	case yaml.MappingNode:
		return processLongSyntaxPort(serviceName, portNode, checker)
	}
	return false
}

// processShortSyntaxPort processes a short syntax port like "8080:80"
func processShortSyntaxPort(serviceName string, portNode *yaml.Node, checker *portChecker) bool {
	portStr := strings.Trim(portNode.Value, "\"'")

	// Extract protocol
	protocol := ""
	if idx := strings.LastIndex(portStr, "/"); idx != -1 {
		protocol = portStr[idx:]
		portStr = portStr[:idx]
	}

	parts := strings.Split(portStr, ":")
	if len(parts) < 2 {
		return false
	}

	var hostIP, hostPortStr, containerPortStr string
	switch len(parts) {
	case 2:
		hostPortStr = parts[0]
		containerPortStr = parts[1]
	case 3:
		hostIP = parts[0]
		hostPortStr = parts[1]
		containerPortStr = parts[2]
	default:
		return false
	}

	if strings.Contains(hostPortStr, "-") {
		return false
	}

	hostPort, err := strconv.Atoi(hostPortStr)
	if err != nil {
		return false
	}

	if checker.isAvailable(hostPort) {
		checker.markUsed(hostPort)
		return false
	}

	newPort, err := checker.findAvailable(hostPort + 1)
	if err != nil {
		slog.Warn("failed to find available port", "service", serviceName, "original", hostPort, "error", err)
		return false
	}

	var newValue string
	if hostIP != "" {
		newValue = fmt.Sprintf("%s:%d:%s", hostIP, newPort, containerPortStr)
	} else {
		newValue = fmt.Sprintf("%d:%s", newPort, containerPortStr)
	}
	if protocol != "" {
		newValue += protocol
	}

	portNode.Value = newValue
	checker.markUsed(newPort)

	slog.Info("port conflict resolved",
		"service", serviceName,
		"original", hostPort,
		"new", newPort,
	)

	return true
}

// processLongSyntaxPort processes a long syntax port mapping
func processLongSyntaxPort(serviceName string, portNode *yaml.Node, checker *portChecker) bool {
	var publishedNode *yaml.Node
	var hostPort int

	for i := 0; i < len(portNode.Content)-1; i += 2 {
		keyNode := portNode.Content[i]
		valueNode := portNode.Content[i+1]

		if keyNode.Value == "published" {
			publishedNode = valueNode
			published := valueNode.Value
			if published == "" {
				return false
			}
			if strings.Contains(published, "-") {
				return false
			}
			if matched, _ := regexp.MatchString(`^\$\{.*\}$`, published); matched {
				return false
			}

			var err error
			hostPort, err = strconv.Atoi(published)
			if err != nil {
				return false
			}
			break
		}
	}

	if publishedNode == nil || hostPort == 0 {
		return false
	}

	if checker.isAvailable(hostPort) {
		checker.markUsed(hostPort)
		return false
	}

	newPort, err := checker.findAvailable(hostPort + 1)
	if err != nil {
		slog.Warn("failed to find available port", "service", serviceName, "original", hostPort, "error", err)
		return false
	}

	publishedNode.Value = strconv.Itoa(newPort)
	checker.markUsed(newPort)

	slog.Info("port conflict resolved",
		"service", serviceName,
		"original", hostPort,
		"new", newPort,
	)

	return true
}
