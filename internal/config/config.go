package config

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Manager defines the interface for configuration management
type Manager interface {
	// Load loads the configuration from disk
	Load(ctx context.Context) error

	// Save saves the configuration to disk
	Save(ctx context.Context) error

	// Get returns the current configuration
	Get() *types.Config

	// Set updates the configuration
	Set(cfg *types.Config)

	// LoadEnvironments loads the environment list from disk
	LoadEnvironments(ctx context.Context) (*types.EnvironmentList, error)

	// SaveEnvironments saves the environment list to disk
	SaveEnvironments(ctx context.Context, envs *types.EnvironmentList) error

	// Paths returns the paths manager
	Paths() *Paths

	// IsInitialized checks if the configuration is initialized
	IsInitialized() bool
}

// ConfigManager implements the Manager interface
type ConfigManager struct {
	paths  *Paths
	config *types.Config
	mu     sync.RWMutex
}

// NewConfigManager creates a new ConfigManager
func NewConfigManager(paths *Paths) *ConfigManager {
	return &ConfigManager{
		paths:  paths,
		config: nil,
	}
}

// Load loads the configuration from disk
func (m *ConfigManager) Load(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.paths.ConfigExists() {
		// Return default config if file doesn't exist
		cfg := types.DefaultConfig()
		m.config = &cfg
		return nil
	}

	data, err := os.ReadFile(m.paths.ConfigFile())
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg types.Config
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		cfg.GitHub.Token = token
	}

	m.config = &cfg
	return nil
}

// Save saves the configuration to disk
func (m *ConfigManager) Save(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.config == nil {
		return fmt.Errorf("no configuration to save")
	}

	if err := m.paths.EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	f, err := os.Create(m.paths.ConfigFile())
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	// Don't save the token to file (it should come from env var)
	cfgToSave := *m.config
	cfgToSave.GitHub.Token = ""

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(cfgToSave); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// Get returns the current configuration
func (m *ConfigManager) Get() *types.Config {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.config == nil {
		cfg := types.DefaultConfig()
		return &cfg
	}
	return m.config
}

// Set updates the configuration
func (m *ConfigManager) Set(cfg *types.Config) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config = cfg
}

// LoadEnvironments loads the environment list from disk
func (m *ConfigManager) LoadEnvironments(ctx context.Context) (*types.EnvironmentList, error) {
	if !m.paths.EnvironmentsFileExists() {
		return nil, fmt.Errorf("environments file not found, please run 'vulhub init' first")
	}

	data, err := os.ReadFile(m.paths.EnvironmentsFile())
	if err != nil {
		return nil, fmt.Errorf("failed to read environments file: %w", err)
	}

	var envs types.EnvironmentList
	if _, err := toml.Decode(string(data), &envs); err != nil {
		return nil, fmt.Errorf("failed to parse environments file: %w", err)
	}

	return &envs, nil
}

// GetEnvironments returns the environments slice (helper method)
func (m *ConfigManager) GetEnvironments(ctx context.Context) ([]types.Environment, error) {
	envList, err := m.LoadEnvironments(ctx)
	if err != nil {
		return nil, err
	}
	return envList.Environment, nil
}

// SaveEnvironments saves the environment list to disk
func (m *ConfigManager) SaveEnvironments(ctx context.Context, envs *types.EnvironmentList) error {
	if err := m.paths.EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	f, err := os.Create(m.paths.EnvironmentsFile())
	if err != nil {
		return fmt.Errorf("failed to create environments file: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(envs); err != nil {
		return fmt.Errorf("failed to encode environments: %w", err)
	}

	return nil
}

// Paths returns the paths manager
func (m *ConfigManager) Paths() *Paths {
	return m.paths
}

// IsInitialized checks if the configuration is initialized
func (m *ConfigManager) IsInitialized() bool {
	return m.paths.ConfigExists() && m.paths.EnvironmentsFileExists()
}
