package alibabacloud

import (
	"fmt"
	"net/http"
	"os"
	
	"github.com/aliyun/alibaba-cloud-ops-mcp-server-go/internal/config"
)

// Credentials represents AlibabaCloud credentials
type Credentials struct {
	AccessKeyID     string
	AccessKeySecret string
	SecurityToken   string
}

// CredentialsProvider provides AlibabaCloud credentials
type CredentialsProvider struct {
	cfg      *config.Config
	settings *config.Settings
}

// NewCredentialsProvider creates a new credentials provider
func NewCredentialsProvider(cfg *config.Config) *CredentialsProvider {
	return &CredentialsProvider{
		cfg:      cfg,
		settings: config.GetSettings(),
	}
}

// GetCredentials gets credentials from various sources
func (p *CredentialsProvider) GetCredentials(r *http.Request) (*Credentials, error) {
	// Try to get credentials from HTTP headers first
	if r != nil {
		if creds := p.getCredentialsFromHeaders(r); creds != nil {
			return creds, nil
		}
	}
	
	// If HeadersCredentialOnly is true and no headers found, return error
	if p.settings.HeadersCredentialOnly {
		return nil, fmt.Errorf("headers credential only mode enabled but no credentials in headers")
	}
	
	// Fall back to environment variables or config
	return p.getCredentialsFromConfig()
}

// GetCredentialsSimple gets credentials without HTTP request context
func (p *CredentialsProvider) GetCredentialsSimple() (*Credentials, error) {
	return p.GetCredentials(nil)
}

// getCredentialsFromHeaders extracts credentials from HTTP headers
func (p *CredentialsProvider) getCredentialsFromHeaders(r *http.Request) *Credentials {
	accessKeyID := r.Header.Get("X-ACS-AccessKey-Id")
	if accessKeyID == "" {
		accessKeyID = r.Header.Get("x-acs-accesskey-id")
	}
	
	if accessKeyID == "" {
		return nil
	}
	
	accessKeySecret := r.Header.Get("X-ACS-AccessKey-Secret")
	if accessKeySecret == "" {
		accessKeySecret = r.Header.Get("x-acs-accesskey-secret")
	}
	
	securityToken := r.Header.Get("X-ACS-Security-Token")
	if securityToken == "" {
		securityToken = r.Header.Get("x-acs-security-token")
	}
	
	return &Credentials{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		SecurityToken:   securityToken,
	}
}

// getCredentialsFromConfig gets credentials from config or environment
func (p *CredentialsProvider) getCredentialsFromConfig() (*Credentials, error) {
	accessKeyID := p.cfg.AccessKeyID
	accessKeySecret := p.cfg.AccessKeySecret
	
	// Try environment variables if not in config
	if accessKeyID == "" {
		accessKeyID = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	}
	if accessKeySecret == "" {
		accessKeySecret = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	}
	
	if accessKeyID == "" || accessKeySecret == "" {
		return nil, fmt.Errorf("credentials not found in config or environment")
	}
	
	securityToken := os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN")
	
	return &Credentials{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		SecurityToken:   securityToken,
	}, nil
}

// IsValid checks if credentials are valid
func (c *Credentials) IsValid() bool {
	return c.AccessKeyID != "" && c.AccessKeySecret != ""
}

