package httpclient

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.uber.org/fx"

	"github.com/vulhub/vulhub-cli/internal/config"
)

// Module provides the httpclient module for fx
var Module = fx.Module("httpclient",
	fx.Provide(NewClientFromConfig),
)

// NewClientFromConfig creates an HTTP client from the config manager
func NewClientFromConfig(cfgMgr config.Manager) (*Client, error) {
	cfg := cfgMgr.Get()

	// Determine proxy URL with priority:
	// 1. VULHUB_PROXY environment variable (highest priority)
	// 2. HTTP_PROXY/HTTPS_PROXY environment variables
	// 3. Config file setting (lowest priority)
	proxyURL := cfg.Network.Proxy

	if envProxy := os.Getenv("VULHUB_PROXY"); envProxy != "" {
		proxyURL = envProxy
	} else if envProxy := os.Getenv("HTTPS_PROXY"); envProxy != "" {
		proxyURL = envProxy
	} else if envProxy := os.Getenv("HTTP_PROXY"); envProxy != "" {
		proxyURL = envProxy
	}

	// Determine timeout
	timeout := DefaultTimeout
	if cfg.Network.Timeout > 0 {
		timeout = time.Duration(cfg.Network.Timeout) * time.Second
	}

	slog.Debug("configuring HTTP client", "proxy", proxyURL, "timeout", timeout)
	httpClient := &http.Client{
		Timeout: timeout,
	}

	if proxyURL != "" {
		transport, err := createProxyTransport(proxyURL)
		if err != nil {
			return nil, err
		}
		httpClient.Transport = transport
	}

	return &Client{
		Client:   httpClient,
		proxyURL: proxyURL,
	}, nil
}

// SetProxyURL updates the proxy URL for the client.
// This is used when --proxy flag is provided via CLI.
func (c *Client) SetProxyURL(proxyURL string) error {
	if proxyURL == "" {
		// Remove proxy
		c.Client = &http.Client{
			Timeout: c.Client.Timeout,
		}
		c.proxyURL = ""
		return nil
	}

	transport, err := createProxyTransport(proxyURL)
	if err != nil {
		return err
	}

	c.Client = &http.Client{
		Transport: transport,
		Timeout:   c.Client.Timeout,
	}
	c.proxyURL = proxyURL
	return nil
}
