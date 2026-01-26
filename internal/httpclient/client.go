package httpclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// Client wraps http.Client and provides proxy configuration
type Client struct {
	*http.Client
	proxyURL string
}

// DefaultTimeout is the default HTTP request timeout
const DefaultTimeout = 30 * time.Second

// ProxyURL returns the configured proxy URL
func (c *Client) ProxyURL() string {
	return c.proxyURL
}

// HasProxy returns true if a proxy is configured
func (c *Client) HasProxy() bool {
	return c.proxyURL != ""
}

// StandardClient returns the underlying *http.Client
func (c *Client) StandardClient() *http.Client {
	return c.Client
}

// createProxyTransport creates an http.Transport with proxy configuration
func createProxyTransport(proxyURL string) (*http.Transport, error) {
	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
	}

	scheme := strings.ToLower(parsed.Scheme)
	switch scheme {
	case "http", "https":
		return &http.Transport{
			Proxy: http.ProxyURL(parsed),
		}, nil
	case "socks5", "socks5h":
		return createSOCKS5Transport(parsed)
	default:
		return nil, fmt.Errorf("unsupported proxy scheme: %s (supported: http, https, socks5)", scheme)
	}
}

// createSOCKS5Transport creates an http.Transport with SOCKS5 proxy
func createSOCKS5Transport(proxyURL *url.URL) (*http.Transport, error) {
	host := proxyURL.Host

	var auth *proxy.Auth
	if proxyURL.User != nil {
		auth = &proxy.Auth{
			User: proxyURL.User.Username(),
		}
		if password, ok := proxyURL.User.Password(); ok {
			auth.Password = password
		}
	}

	dialer, err := proxy.SOCKS5("tcp", host, auth, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

	return &http.Transport{
		Dial: dialer.Dial,
	}, nil
}

// ValidateProxyURL validates a proxy URL string
func ValidateProxyURL(proxyURL string) error {
	if proxyURL == "" {
		return nil
	}

	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}

	scheme := strings.ToLower(parsed.Scheme)
	switch scheme {
	case "http", "https", "socks5", "socks5h":
		// Valid schemes
	default:
		return fmt.Errorf("unsupported proxy scheme: %s (supported: http, https, socks5)", scheme)
	}

	if parsed.Host == "" {
		return fmt.Errorf("proxy URL must have a host")
	}

	return nil
}
