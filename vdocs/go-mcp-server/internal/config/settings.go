package config

import "fmt"

// Settings represents runtime settings
type Settings struct {
	HeadersCredentialOnly bool
	Environment           string
}

var globalSettings = &Settings{
	HeadersCredentialOnly: false,
	Environment:           "domestic",
}

// GetSettings returns the global settings instance
func GetSettings() *Settings {
	return globalSettings
}

// SetHeadersCredentialOnly sets the headers credential only flag
func (s *Settings) SetHeadersCredentialOnly(value bool) {
	s.HeadersCredentialOnly = value
}

// SetEnvironment sets the environment
func (s *Settings) SetEnvironment(env string) error {
	if env != "domestic" && env != "international" {
		return fmt.Errorf("invalid environment: %s", env)
	}
	s.Environment = env
	return nil
}

// IsInternational returns true if the environment is international
func (s *Settings) IsInternational() bool {
	return s.Environment == "international"
}

// IsDomestic returns true if the environment is domestic
func (s *Settings) IsDomestic() bool {
	return s.Environment == "domestic"
}

