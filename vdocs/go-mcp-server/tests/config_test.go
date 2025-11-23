package tests

import (
	"os"
	"testing"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	
	if cfg.Transport != "stdio" {
		t.Errorf("Expected default transport stdio, got %s", cfg.Transport)
	}
	
	if cfg.Port != 8000 {
		t.Errorf("Expected default port 8000, got %d", cfg.Port)
	}
	
	if cfg.Environment != "domestic" {
		t.Errorf("Expected default environment domestic, got %s", cfg.Environment)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("MCP_TRANSPORT", "http")
	os.Setenv("MCP_PORT", "9000")
	os.Setenv("ALIBABA_CLOUD_ACCESS_KEY_ID", "test-id")
	os.Setenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET", "test-secret")
	defer func() {
		os.Unsetenv("MCP_TRANSPORT")
		os.Unsetenv("MCP_PORT")
		os.Unsetenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
		os.Unsetenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	}()
	
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	
	if cfg.Transport != "http" {
		t.Errorf("Expected transport http, got %s", cfg.Transport)
	}
	
	if cfg.AccessKeyID != "test-id" {
		t.Errorf("Expected access key ID test-id, got %s", cfg.AccessKeyID)
	}
	
	if cfg.AccessKeySecret != "test-secret" {
		t.Errorf("Expected access key secret test-secret, got %s", cfg.AccessKeySecret)
	}
}

func TestValidateConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AccessKeyID = "test"
	cfg.AccessKeySecret = "test"
	
	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
	
	// Test invalid transport
	cfg.Transport = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for invalid transport")
	}
	
	// Test invalid environment
	cfg.Transport = "stdio"
	cfg.Environment = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for invalid environment")
	}
}

func TestSettings(t *testing.T) {
	settings := config.GetSettings()
	
	if settings.IsDomestic() != true {
		t.Error("Expected default environment to be domestic")
	}
	
	err := settings.SetEnvironment("international")
	if err != nil {
		t.Errorf("Failed to set environment: %v", err)
	}
	
	if settings.IsInternational() != true {
		t.Error("Expected environment to be international")
	}
	
	settings.SetHeadersCredentialOnly(true)
	if settings.HeadersCredentialOnly != true {
		t.Error("Expected HeadersCredentialOnly to be true")
	}
}

