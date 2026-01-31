package api

import "github.com/vulhub/vulhub-cli/pkg/types"

// Response is the standard API response wrapper
type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// EnvironmentResponse represents an environment in API responses
type EnvironmentResponse struct {
	Path       string   `json:"path"`
	Name       string   `json:"name"`
	CVE        []string `json:"cve,omitempty"`
	App        string   `json:"app"`
	Tags       []string `json:"tags,omitempty"`
	Downloaded bool     `json:"downloaded"`
	Running    bool     `json:"running"`
}

// EnvironmentStatusResponse represents the status of an environment
type EnvironmentStatusResponse struct {
	Environment EnvironmentResponse       `json:"environment"`
	Containers  []ContainerStatusResponse `json:"containers,omitempty"`
	Running     bool                      `json:"running"`
	LocalPath   string                    `json:"local_path,omitempty"`
}

// ContainerStatusResponse represents a container status
type ContainerStatusResponse struct {
	ID        string                `json:"id"`
	Name      string                `json:"name"`
	Image     string                `json:"image"`
	Status    string                `json:"status"`
	State     string                `json:"state"`
	Ports     []PortMappingResponse `json:"ports,omitempty"`
	CreatedAt string                `json:"created_at,omitempty"`
	StartedAt string                `json:"started_at,omitempty"`
}

// PortMappingResponse represents a port mapping
type PortMappingResponse struct {
	HostIP        string `json:"host_ip,omitempty"`
	HostPort      string `json:"host_port"`
	ContainerPort string `json:"container_port"`
	Protocol      string `json:"protocol"`
}

// EnvironmentInfoResponse represents detailed environment information
type EnvironmentInfoResponse struct {
	Environment EnvironmentResponse `json:"environment"`
	Readme      string              `json:"readme,omitempty"`
	ComposeFile string              `json:"compose_file,omitempty"`
	Downloaded  bool                `json:"downloaded"`
	LocalPath   string              `json:"local_path,omitempty"`
}

// EnvironmentListResponse represents a list of environments
type EnvironmentListResponse struct {
	Environments []EnvironmentResponse `json:"environments"`
	Total        int                   `json:"total"`
}

// StatusListResponse represents a list of environment statuses
type StatusListResponse struct {
	Environments []EnvironmentStatusResponse `json:"environments"`
	Total        int                         `json:"total"`
}

// StartRequest represents a request to start an environment
type StartRequest struct {
	Pull          bool `json:"pull"`
	Build         bool `json:"build"`
	ForceRecreate bool `json:"force_recreate"`
	SkipPortCheck bool `json:"skip_port_check"`
}

// CleanRequest represents a request to clean an environment
type CleanRequest struct {
	RemoveVolumes bool `json:"remove_volumes"`
	RemoveImages  bool `json:"remove_images"`
	RemoveFiles   bool `json:"remove_files"`
}

// SystemStatusResponse represents system status
type SystemStatusResponse struct {
	Initialized  bool   `json:"initialized"`
	LastSyncTime string `json:"last_sync_time,omitempty"`
	NeedSync     bool   `json:"need_sync"`
	Version      string `json:"version"`
}

// SyncupResponse represents the result of a syncup operation
type SyncupResponse struct {
	Success      bool   `json:"success"`
	LastSyncTime string `json:"last_sync_time"`
	Total        int    `json:"total"`
}

// Helper functions to convert types

func toEnvironmentResponse(env types.Environment, downloaded, running bool) EnvironmentResponse {
	return EnvironmentResponse{
		Path:       env.Path,
		Name:       env.Name,
		CVE:        env.CVE,
		App:        env.App,
		Tags:       env.Tags,
		Downloaded: downloaded,
		Running:    running,
	}
}

func toContainerStatusResponse(c types.ContainerStatus) ContainerStatusResponse {
	ports := make([]PortMappingResponse, len(c.Ports))
	for i, p := range c.Ports {
		ports[i] = PortMappingResponse{
			HostIP:        p.HostIP,
			HostPort:      p.HostPort,
			ContainerPort: p.ContainerPort,
			Protocol:      p.Protocol,
		}
	}
	return ContainerStatusResponse{
		ID:        c.ID,
		Name:      c.Name,
		Image:     c.Image,
		Status:    c.Status,
		State:     c.State,
		Ports:     ports,
		CreatedAt: c.CreatedAt,
		StartedAt: c.StartedAt,
	}
}

func toEnvironmentStatusResponse(s types.EnvironmentStatus) EnvironmentStatusResponse {
	containers := make([]ContainerStatusResponse, len(s.Containers))
	for i, c := range s.Containers {
		containers[i] = toContainerStatusResponse(c)
	}
	return EnvironmentStatusResponse{
		Environment: toEnvironmentResponse(s.Environment, true, s.Running),
		Containers:  containers,
		Running:     s.Running,
		LocalPath:   s.LocalPath,
	}
}

func toEnvironmentInfoResponse(info *types.EnvironmentInfo) EnvironmentInfoResponse {
	return EnvironmentInfoResponse{
		Environment: toEnvironmentResponse(info.Environment, info.Downloaded, false),
		Readme:      info.Readme,
		ComposeFile: info.ComposeFile,
		Downloaded:  info.Downloaded,
		LocalPath:   info.LocalPath,
	}
}
