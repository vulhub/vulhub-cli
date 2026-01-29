package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server represents the API server
type Server struct {
	host     string
	port     int
	engine   *gin.Engine
	handlers *Handlers
	server   *http.Server
}

// NewServer creates a new API server
func NewServer(
	host string,
	port int,
	handlers *Handlers,
) *Server {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	// Create Gin engine
	engine := gin.New()

	// Add middleware
	engine.Use(gin.Recovery())
	engine.Use(requestLogger())
	engine.Use(corsMiddleware())

	s := &Server{
		host:     host,
		port:     port,
		engine:   engine,
		handlers: handlers,
	}

	// Setup routes
	s.setupRoutes()

	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := s.engine.Group("/api/v1")
	{
		// System
		v1.GET("/status", s.handlers.GetSystemStatus)
		v1.POST("/syncup", s.handlers.Syncup)

		// Environments
		envs := v1.Group("/environments")
		{
			envs.GET("", s.handlers.ListAllEnvironments)
			envs.GET("/downloaded", s.handlers.ListDownloadedEnvironments)
			envs.GET("/running", s.handlers.ListRunningEnvironments)

			// Single environment operations (using wildcard for paths like log4j/CVE-2021-44228)
			// Note: We use separate route groups because Gin doesn't allow multiple wildcards
			// The path format is: /api/v1/environments/info/log4j/CVE-2021-44228
			envs.GET("/info/*path", s.handlers.GetEnvironmentInfo)
			envs.GET("/status/*path", s.handlers.GetEnvironmentStatus)
			envs.POST("/start/*path", s.handlers.StartEnvironment)
			envs.POST("/stop/*path", s.handlers.StopEnvironment)
			envs.POST("/restart/*path", s.handlers.RestartEnvironment)
			envs.DELETE("/clean/*path", s.handlers.CleanEnvironment)
		}
	}
}

// Start starts the API server
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	slog.Info("starting API server", "address", addr)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("API server error", "error", err)
		}
	}()

	return nil
}

// Stop gracefully stops the API server
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	slog.Info("stopping API server")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.server.Shutdown(shutdownCtx)
}

// Middleware

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log after request is processed
		duration := time.Since(start)
		slog.Debug("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration,
			"client_ip", c.ClientIP(),
		)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
