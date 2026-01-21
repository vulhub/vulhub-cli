package config

import (
	"os"
	"path/filepath"
)

const (
	// DefaultConfigDirName is the default configuration directory name
	DefaultConfigDirName = ".vulhub"

	// ConfigFileName is the main configuration file name
	ConfigFileName = "config.toml"

	// EnvironmentsFileName is the environments list file name
	EnvironmentsFileName = "environments.toml"

	// EnvironmentsDirName is the directory name for downloaded environments
	EnvironmentsDirName = "environments"
)

// Paths provides path management for vulhub-cli
type Paths struct {
	configDir string
}

// NewPaths creates a new Paths instance with the default config directory
func NewPaths() (*Paths, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &Paths{
		configDir: filepath.Join(homeDir, DefaultConfigDirName),
	}, nil
}

// NewPathsWithDir creates a new Paths instance with a custom config directory
func NewPathsWithDir(configDir string) *Paths {
	return &Paths{
		configDir: configDir,
	}
}

// ConfigDir returns the configuration directory path
func (p *Paths) ConfigDir() string {
	return p.configDir
}

// ConfigFile returns the main configuration file path
func (p *Paths) ConfigFile() string {
	return filepath.Join(p.configDir, ConfigFileName)
}

// EnvironmentsFile returns the environments list file path
func (p *Paths) EnvironmentsFile() string {
	return filepath.Join(p.configDir, EnvironmentsFileName)
}

// EnvironmentsDir returns the directory for downloaded environments
func (p *Paths) EnvironmentsDir() string {
	return filepath.Join(p.configDir, EnvironmentsDirName)
}

// EnvironmentDir returns the directory for a specific environment
func (p *Paths) EnvironmentDir(envPath string) string {
	return filepath.Join(p.EnvironmentsDir(), envPath)
}

// EnsureConfigDir ensures the configuration directory exists
func (p *Paths) EnsureConfigDir() error {
	return os.MkdirAll(p.configDir, 0755)
}

// EnsureEnvironmentsDir ensures the environments directory exists
func (p *Paths) EnsureEnvironmentsDir() error {
	return os.MkdirAll(p.EnvironmentsDir(), 0755)
}

// EnsureEnvironmentDir ensures a specific environment directory exists
func (p *Paths) EnsureEnvironmentDir(envPath string) error {
	return os.MkdirAll(p.EnvironmentDir(envPath), 0755)
}

// Exists checks if a path exists
func (p *Paths) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ConfigExists checks if the configuration file exists
func (p *Paths) ConfigExists() bool {
	return p.Exists(p.ConfigFile())
}

// EnvironmentsFileExists checks if the environments file exists
func (p *Paths) EnvironmentsFileExists() bool {
	return p.Exists(p.EnvironmentsFile())
}

// EnvironmentExists checks if a specific environment is downloaded
func (p *Paths) EnvironmentExists(envPath string) bool {
	return p.Exists(p.EnvironmentDir(envPath))
}
