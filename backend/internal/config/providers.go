package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	EnvTransportProviderBaseURL = "TRANSPORT_PROVIDER_BASE_URL"
	EnvTransportProviderAPIKey  = "TRANSPORT_PROVIDER_API_KEY"
)

// ProviderConfig contains typed runtime configuration for a transport provider.
type ProviderConfig struct {
	Transport TransportProviderConfig
}

// TransportProviderConfig stores credentials and endpoint for transport search provider.
type TransportProviderConfig struct {
	BaseURL string
	APIKey  string
}

func LoadProviderConfigFromEnv() (ProviderConfig, error) {
	baseURL := strings.TrimSpace(os.Getenv(EnvTransportProviderBaseURL))
	apiKey := strings.TrimSpace(os.Getenv(EnvTransportProviderAPIKey))

	var missing []string
	if baseURL == "" {
		missing = append(missing, EnvTransportProviderBaseURL)
	}
	if apiKey == "" {
		missing = append(missing, EnvTransportProviderAPIKey)
	}
	if len(missing) > 0 {
		return ProviderConfig{}, fmt.Errorf("provider config invalid: %w", MissingEnvVarsError{Vars: missing})
	}

	return ProviderConfig{
		Transport: TransportProviderConfig{
			BaseURL: baseURL,
			APIKey:  apiKey,
		},
	}, nil
}

// MissingEnvVarsError reports one or more missing env vars.
type MissingEnvVarsError struct {
	Vars []string
}

func (e MissingEnvVarsError) Error() string {
	if len(e.Vars) == 0 {
		return "missing required env vars"
	}
	return fmt.Sprintf("missing required env vars: %s", strings.Join(e.Vars, ", "))
}

func (e MissingEnvVarsError) Is(target error) bool {
	_, ok := target.(MissingEnvVarsError)
	return ok
}

func IsMissingEnvVars(err error) bool {
	return errors.Is(err, MissingEnvVarsError{})
}
