package config

import (
	"fmt"
	"os"
	"strings"
)

// Config represents the application configuration
type Config struct {
	// Transport configuration
	Transport string
	Host      string
	Port      int
	
	// AlibabaCloud credentials
	AccessKeyID     string
	AccessKeySecret string
	
	// Server configuration
	Services                []string
	Environment             string // domestic or international
	HeadersCredentialOnly   bool
	
	// Server info
	ServerName    string
	ServerVersion string
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Transport:             "stdio",
		Host:                  "127.0.0.1",
		Port:                  8000,
		Environment:           "domestic",
		HeadersCredentialOnly: false,
		ServerName:            "alibaba-cloud-ops-mcp-server-go",
		ServerVersion:         "1.0.0",
	}
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	if v := os.Getenv("MCP_TRANSPORT"); v != "" {
		c.Transport = v
	}
	if v := os.Getenv("MCP_HOST"); v != "" {
		c.Host = v
	}
	if v := os.Getenv("MCP_PORT"); v != "" {
		// Parse port from string
		var port int
		if _, err := fmt.Sscanf(v, "%d", &port); err == nil {
			c.Port = port
		}
	}
	
	// AlibabaCloud credentials
	if v := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"); v != "" {
		c.AccessKeyID = v
	}
	if v := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"); v != "" {
		c.AccessKeySecret = v
	}
	
	// Services
	if v := os.Getenv("MCP_SERVICES"); v != "" {
		c.Services = parseServices(v)
	}
	
	// Environment
	if v := os.Getenv("MCP_ENV"); v != "" {
		c.Environment = v
	}
	
	// Headers credential only
	if v := os.Getenv("MCP_HEADERS_CREDENTIAL_ONLY"); v == "true" {
		c.HeadersCredentialOnly = true
	}
}

// parseServices parses comma-separated services string
func parseServices(s string) []string {
	parts := strings.Split(s, ",")
	services := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			services = append(services, strings.ToLower(trimmed))
		}
	}
	return services
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate transport
	validTransports := map[string]bool{
		"stdio": true,
		"http":  true,
		"sse":   true,
	}
	if !validTransports[c.Transport] {
		return fmt.Errorf("invalid transport: %s, must be one of: stdio, http, sse", c.Transport)
	}
	
	// Validate environment
	if c.Environment != "domestic" && c.Environment != "international" {
		return fmt.Errorf("invalid environment: %s, must be domestic or international", c.Environment)
	}
	
	// Validate credentials (unless HeadersCredentialOnly is true)
	if !c.HeadersCredentialOnly {
		if c.AccessKeyID == "" || c.AccessKeySecret == "" {
			return fmt.Errorf("ALIBABA_CLOUD_ACCESS_KEY_ID and ALIBABA_CLOUD_ACCESS_KEY_SECRET must be set")
		}
	}
	
	return nil
}

