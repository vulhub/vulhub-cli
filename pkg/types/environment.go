package types

// Environment represents a vulnerability environment from vulhub
type Environment struct {
	// Path is the relative path in the vulhub repository (e.g., "log4j/CVE-2021-44228")
	Path string `toml:"path"`

	// Name is the human-readable name of the environment
	Name string `toml:"name"`

	// CVE is the list of CVE identifiers (e.g., ["CVE-2021-44228"])
	CVE []string `toml:"cve,omitempty"`

	// App is the application name (e.g., "log4j", "apache", "tomcat")
	App string `toml:"app"`

	// Tags are searchable tags for the environment
	Tags []string `toml:"tags,omitempty"`
}

// GetCVE returns the first CVE identifier or empty string
func (e *Environment) GetCVE() string {
	if len(e.CVE) > 0 {
		return e.CVE[0]
	}
	return ""
}

// HasCVE checks if the environment has a specific CVE
func (e *Environment) HasCVE(cve string) bool {
	for _, c := range e.CVE {
		if c == cve {
			return true
		}
	}
	return false
}

// EnvironmentList represents the list of environments from environments.toml
type EnvironmentList struct {
	// Tags is the global list of available tags
	Tags []string `toml:"tags,omitempty"`

	// Environment is the list of all available environments (note: singular in TOML)
	Environment []Environment `toml:"environment"`
}

// ContainerStatus represents the status of a Docker container
type ContainerStatus struct {
	// ID is the container ID
	ID string

	// Name is the container name
	Name string

	// Image is the container image
	Image string

	// Status is the container status (e.g., "running", "exited")
	Status string

	// State is the container state (e.g., "running", "exited")
	State string

	// Ports is the list of exposed ports
	Ports []PortMapping

	// CreatedAt is the creation time
	CreatedAt string

	// StartedAt is the start time
	StartedAt string
}

// PortMapping represents a port mapping from container to host
type PortMapping struct {
	// HostIP is the host IP address
	HostIP string

	// HostPort is the host port number
	HostPort string

	// ContainerPort is the container port number
	ContainerPort string

	// Protocol is the protocol (tcp/udp)
	Protocol string
}

// EnvironmentStatus represents the status of an environment
type EnvironmentStatus struct {
	// Environment is the environment definition
	Environment Environment

	// Containers is the list of container statuses
	Containers []ContainerStatus

	// Running indicates if any container is running
	Running bool

	// LocalPath is the local path where the environment is downloaded
	LocalPath string
}

// EnvironmentInfo represents detailed information about an environment
type EnvironmentInfo struct {
	// Environment is the base environment definition
	Environment Environment

	// Readme is the content of the README file
	Readme string

	// ComposeFile is the content of the docker-compose.yml file
	ComposeFile string

	// Downloaded indicates if the environment is downloaded locally
	Downloaded bool

	// LocalPath is the local path where the environment is downloaded
	LocalPath string
}
