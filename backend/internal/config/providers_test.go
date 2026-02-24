package config

import (
	"errors"
	"testing"
)

func TestLoadProviderConfigFromEnv_Success(t *testing.T) {
	t.Setenv(EnvTransportProviderBaseURL, "https://provider.example.com")
	t.Setenv(EnvTransportProviderAPIKey, "secret-token")

	cfg, err := LoadProviderConfigFromEnv()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if cfg.Transport.BaseURL != "https://provider.example.com" {
		t.Fatalf("unexpected base url: %s", cfg.Transport.BaseURL)
	}
	if cfg.Transport.APIKey != "secret-token" {
		t.Fatalf("unexpected api key: %s", cfg.Transport.APIKey)
	}
}

func TestLoadProviderConfigFromEnv_MissingVars(t *testing.T) {
	t.Setenv(EnvTransportProviderBaseURL, "")
	t.Setenv(EnvTransportProviderAPIKey, "")

	_, err := LoadProviderConfigFromEnv()
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsMissingEnvVars(err) {
		t.Fatalf("expected missing env vars error, got %T: %v", err, err)
	}
	if !errors.Is(err, MissingEnvVarsError{}) {
		t.Fatalf("expected errors.Is for MissingEnvVarsError to be true, got %v", err)
	}
	want := "provider config invalid: missing required env vars: TRANSPORT_PROVIDER_BASE_URL, TRANSPORT_PROVIDER_API_KEY"
	if err.Error() != want {
		t.Fatalf("unexpected error message\nwant: %s\ngot:  %s", want, err.Error())
	}
}

func TestLoadProviderConfigFromEnv_MissingSingleVar(t *testing.T) {
	t.Setenv(EnvTransportProviderBaseURL, "https://provider.example.com")
	t.Setenv(EnvTransportProviderAPIKey, "")

	_, err := LoadProviderConfigFromEnv()
	if err == nil {
		t.Fatal("expected error")
	}
	want := "provider config invalid: missing required env vars: TRANSPORT_PROVIDER_API_KEY"
	if err.Error() != want {
		t.Fatalf("unexpected error message\nwant: %s\ngot:  %s", want, err.Error())
	}
}
