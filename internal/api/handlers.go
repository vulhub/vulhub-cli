package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/vulhub/vulhub-cli/internal/config"
	"github.com/vulhub/vulhub-cli/internal/environment"
	"github.com/vulhub/vulhub-cli/internal/github"
	"github.com/vulhub/vulhub-cli/internal/resolver"
	"github.com/vulhub/vulhub-cli/pkg/types"
)

// Version is set at build time
var Version = "dev"

// Handlers contains all API handlers
type Handlers struct {
	config      config.Manager
	environment environment.Manager
	resolver    resolver.Resolver
	downloader  *github.Downloader
}

// NewHandlers creates a new Handlers instance
func NewHandlers(
	cfg config.Manager,
	env environment.Manager,
	res resolver.Resolver,
	downloader *github.Downloader,
) *Handlers {
	return &Handlers{
		config:      cfg,
		environment: env,
		resolver:    res,
		downloader:  downloader,
	}
}

// success returns a successful response
func success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// errorResponse returns an error response
func errorResponse(c *gin.Context, status int, code, message string) {
	c.JSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// GetSystemStatus returns system status
func (h *Handlers) GetSystemStatus(c *gin.Context) {
	resp := SystemStatusResponse{
		Initialized: h.config.IsInitialized(),
		NeedSync:    h.config.NeedSync(),
		Version:     Version,
	}

	if !h.config.GetLastSyncTime().IsZero() {
		resp.LastSyncTime = h.config.GetLastSyncTime().Format("2006-01-02 15:04:05")
	}

	success(c, resp)
}

// ListAllEnvironments returns all available environments
func (h *Handlers) ListAllEnvironments(c *gin.Context) {
	envList, err := h.config.LoadEnvironments(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "LOAD_ERROR", err.Error())
		return
	}

	// Get downloaded and running environments for status
	downloaded := h.getDownloadedEnvPaths()
	running := h.getRunningEnvPaths(c)

	envs := make([]EnvironmentResponse, len(envList.Environment))
	for i, env := range envList.Environment {
		envs[i] = toEnvironmentResponse(env, downloaded[env.Path], running[env.Path])
	}

	success(c, EnvironmentListResponse{
		Environments: envs,
		Total:        len(envs),
	})
}

// ListDownloadedEnvironments returns all downloaded environments with status
func (h *Handlers) ListDownloadedEnvironments(c *gin.Context) {
	statuses, err := h.environment.ListDownloaded(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "LIST_ERROR", err.Error())
		return
	}

	envStatuses := make([]EnvironmentStatusResponse, len(statuses))
	for i, s := range statuses {
		envStatuses[i] = toEnvironmentStatusResponse(s)
	}

	success(c, StatusListResponse{
		Environments: envStatuses,
		Total:        len(envStatuses),
	})
}

// ListRunningEnvironments returns all running environments
func (h *Handlers) ListRunningEnvironments(c *gin.Context) {
	statuses, err := h.environment.ListRunning(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "LIST_ERROR", err.Error())
		return
	}

	envStatuses := make([]EnvironmentStatusResponse, len(statuses))
	for i, s := range statuses {
		envStatuses[i] = toEnvironmentStatusResponse(s)
	}

	success(c, StatusListResponse{
		Environments: envStatuses,
		Total:        len(envStatuses),
	})
}

// GetEnvironmentInfo returns detailed information about an environment
func (h *Handlers) GetEnvironmentInfo(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		errorResponse(c, http.StatusBadRequest, "MISSING_PATH", "Environment path is required")
		return
	}

	// URL decode and clean the path
	path = strings.TrimPrefix(path, "/")

	env, err := h.resolveEnvironment(c, path)
	if err != nil {
		return // Error already sent
	}

	info, err := h.environment.GetInfo(c.Request.Context(), *env)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "INFO_ERROR", err.Error())
		return
	}

	success(c, toEnvironmentInfoResponse(info))
}

// GetEnvironmentStatus returns the status of an environment
func (h *Handlers) GetEnvironmentStatus(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		errorResponse(c, http.StatusBadRequest, "MISSING_PATH", "Environment path is required")
		return
	}

	path = strings.TrimPrefix(path, "/")

	env, err := h.resolveEnvironment(c, path)
	if err != nil {
		return
	}

	status, err := h.environment.Status(c.Request.Context(), *env)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "STATUS_ERROR", err.Error())
		return
	}

	success(c, toEnvironmentStatusResponse(*status))
}

// StartEnvironment starts an environment
func (h *Handlers) StartEnvironment(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		errorResponse(c, http.StatusBadRequest, "MISSING_PATH", "Environment path is required")
		return
	}

	path = strings.TrimPrefix(path, "/")

	env, err := h.resolveEnvironment(c, path)
	if err != nil {
		return
	}

	var req StartRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	opts := environment.StartOptions{
		Pull:          req.Pull,
		Build:         req.Build,
		ForceRecreate: req.ForceRecreate,
	}

	if err := h.environment.Start(c.Request.Context(), *env, opts); err != nil {
		errorResponse(c, http.StatusInternalServerError, "START_ERROR", err.Error())
		return
	}

	// Get status after starting
	status, err := h.environment.Status(c.Request.Context(), *env)
	if err != nil {
		// Started but couldn't get status
		success(c, gin.H{"message": "Environment started successfully"})
		return
	}

	success(c, toEnvironmentStatusResponse(*status))
}

// StopEnvironment stops an environment
func (h *Handlers) StopEnvironment(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		errorResponse(c, http.StatusBadRequest, "MISSING_PATH", "Environment path is required")
		return
	}

	path = strings.TrimPrefix(path, "/")

	env, err := h.resolveEnvironment(c, path)
	if err != nil {
		return
	}

	if !h.environment.IsDownloaded(*env) {
		errorResponse(c, http.StatusBadRequest, "NOT_DOWNLOADED", "Environment is not downloaded")
		return
	}

	if err := h.environment.Stop(c.Request.Context(), *env); err != nil {
		errorResponse(c, http.StatusInternalServerError, "STOP_ERROR", err.Error())
		return
	}

	success(c, gin.H{"message": "Environment stopped successfully"})
}

// RestartEnvironment restarts an environment
func (h *Handlers) RestartEnvironment(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		errorResponse(c, http.StatusBadRequest, "MISSING_PATH", "Environment path is required")
		return
	}

	path = strings.TrimPrefix(path, "/")

	env, err := h.resolveEnvironment(c, path)
	if err != nil {
		return
	}

	if !h.environment.IsDownloaded(*env) {
		errorResponse(c, http.StatusBadRequest, "NOT_DOWNLOADED", "Environment is not downloaded")
		return
	}

	if err := h.environment.Restart(c.Request.Context(), *env); err != nil {
		errorResponse(c, http.StatusInternalServerError, "RESTART_ERROR", err.Error())
		return
	}

	// Get status after restarting
	status, err := h.environment.Status(c.Request.Context(), *env)
	if err != nil {
		success(c, gin.H{"message": "Environment restarted successfully"})
		return
	}

	success(c, toEnvironmentStatusResponse(*status))
}

// CleanEnvironment cleans up an environment
func (h *Handlers) CleanEnvironment(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		errorResponse(c, http.StatusBadRequest, "MISSING_PATH", "Environment path is required")
		return
	}

	path = strings.TrimPrefix(path, "/")

	env, err := h.resolveEnvironment(c, path)
	if err != nil {
		return
	}

	var req CleanRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		errorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Default to removing files if no options specified
	if !req.RemoveVolumes && !req.RemoveImages && !req.RemoveFiles {
		req.RemoveFiles = true
		req.RemoveVolumes = true
	}

	opts := environment.CleanOptions{
		RemoveVolumes: req.RemoveVolumes,
		RemoveImages:  req.RemoveImages,
		RemoveFiles:   req.RemoveFiles,
	}

	if err := h.environment.Clean(c.Request.Context(), *env, opts); err != nil {
		errorResponse(c, http.StatusInternalServerError, "CLEAN_ERROR", err.Error())
		return
	}

	success(c, gin.H{"message": "Environment cleaned successfully"})
}

// Syncup synchronizes the environment list from GitHub
func (h *Handlers) Syncup(c *gin.Context) {
	ctx := c.Request.Context()

	// Download environment list
	envData, err := h.downloader.DownloadEnvironmentsList(ctx)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "SYNC_ERROR", err.Error())
		return
	}

	// Parse the TOML data
	var envList types.EnvironmentList
	if _, err := toml.Decode(string(envData), &envList); err != nil {
		errorResponse(c, http.StatusInternalServerError, "PARSE_ERROR", err.Error())
		return
	}

	// Save environment list
	if err := h.config.SaveEnvironments(ctx, &envList); err != nil {
		errorResponse(c, http.StatusInternalServerError, "SAVE_ERROR", err.Error())
		return
	}

	// Update last sync time
	if err := h.config.UpdateLastSyncTime(ctx); err != nil {
		errorResponse(c, http.StatusInternalServerError, "UPDATE_TIME_ERROR", err.Error())
		return
	}

	success(c, SyncupResponse{
		Success:      true,
		LastSyncTime: h.config.GetLastSyncTime().Format("2006-01-02 15:04:05"),
		Total:        len(envList.Environment),
	})
}

// Helper methods

func (h *Handlers) resolveEnvironment(c *gin.Context, path string) (*types.Environment, error) {
	result, err := h.resolver.Resolve(c.Request.Context(), path)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "RESOLVE_ERROR", err.Error())
		return nil, err
	}

	if result.HasNoMatches() {
		errorResponse(c, http.StatusNotFound, "NOT_FOUND", "Environment not found: "+path)
		return nil, errNotFound
	}

	if result.HasMultipleMatches() {
		// Return the matches for the client to choose
		matches := result.GetMatchedEnvironments()
		paths := lo.Map(matches, func(e types.Environment, _ int) string { return e.Path })
		errorResponse(c, http.StatusConflict, "MULTIPLE_MATCHES",
			"Multiple environments match. Please specify: "+strings.Join(paths, ", "))
		return nil, errMultipleMatches
	}

	return result.Environment, nil
}

func (h *Handlers) getDownloadedEnvPaths() map[string]bool {
	downloaded := make(map[string]bool)
	envList, err := h.config.LoadEnvironments(context.Background())
	if err != nil {
		return downloaded
	}
	for _, env := range envList.Environment {
		if h.environment.IsDownloaded(env) {
			downloaded[env.Path] = true
		}
	}
	return downloaded
}

func (h *Handlers) getRunningEnvPaths(c *gin.Context) map[string]bool {
	running := make(map[string]bool)
	statuses, err := h.environment.ListRunning(c.Request.Context())
	if err != nil {
		return running
	}
	for _, s := range statuses {
		running[s.Environment.Path] = true
	}
	return running
}

// Sentinel errors
var (
	errNotFound        = &apiError{message: "not found"}
	errMultipleMatches = &apiError{message: "multiple matches"}
)

type apiError struct {
	message string
}

func (e *apiError) Error() string {
	return e.message
}
