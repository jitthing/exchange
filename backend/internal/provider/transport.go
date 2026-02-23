package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"exchange-travel-planner/backend/internal/domain"
)

const (
	defaultProviderBaseURL = "https://transport.opendata.ch/v1"
	defaultProviderTimeout = 2500 * time.Millisecond
)

type TransportProvider interface {
	SearchTransport(ctx context.Context, from, to string) ([]domain.TransportOption, error)
}

type OpenTransportProvider struct {
	enabled bool
	baseURL string
	client  *http.Client
}

func NewOpenTransportProviderFromEnv() *OpenTransportProvider {
	enabled := strings.EqualFold(os.Getenv("REAL_PROVIDER_ENABLED"), "true")
	baseURL := strings.TrimSpace(os.Getenv("REAL_PROVIDER_BASE_URL"))
	if baseURL == "" {
		baseURL = defaultProviderBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	timeout := defaultProviderTimeout
	if raw := strings.TrimSpace(os.Getenv("REAL_PROVIDER_TIMEOUT_MS")); raw != "" {
		if ms, err := strconv.Atoi(raw); err == nil && ms > 0 {
			timeout = time.Duration(ms) * time.Millisecond
		}
	}

	return &OpenTransportProvider{
		enabled: enabled,
		baseURL: baseURL,
		client:  &http.Client{Timeout: timeout},
	}
}

func (p *OpenTransportProvider) SearchTransport(ctx context.Context, from, to string) ([]domain.TransportOption, error) {
	if !p.enabled {
		return nil, fmt.Errorf("provider disabled")
	}

	u, err := url.Parse(p.baseURL + "/connections")
	if err != nil {
		return nil, fmt.Errorf("parse provider url: %w", err)
	}

	q := u.Query()
	q.Set("to", to)
	if strings.TrimSpace(from) != "" {
		q.Set("from", from)
	}
	q.Set("limit", "3")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create provider request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("provider request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provider returned status %d", resp.StatusCode)
	}

	var payload struct {
		Connections []struct {
			Duration  string   `json:"duration"`
			Products  []string `json:"products"`
			Transfers int      `json:"transfers"`
		} `json:"connections"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode provider response: %w", err)
	}

	options := make([]domain.TransportOption, 0, len(payload.Connections))
	for _, c := range payload.Connections {
		hours, ok := parseDurationHours(c.Duration)
		if !ok {
			continue
		}

		mode := inferMode(c.Products)
		price := math.Round((18.0+hours*16.0+float64(c.Transfers*5.0))*100) / 100

		options = append(options, domain.TransportOption{
			Provider:      "OpenTransportData",
			Mode:          mode,
			DurationHours: math.Round(hours*10) / 10,
			Price:         price,
			Deeplink:      u.String(),
		})
	}

	if len(options) == 0 {
		return nil, fmt.Errorf("provider returned no transport options")
	}
	return options, nil
}

func inferMode(products []string) string {
	for _, product := range products {
		normalized := strings.ToUpper(strings.TrimSpace(product))
		if strings.Contains(normalized, "BUS") {
			return "bus"
		}
	}
	return "train"
}

func parseDurationHours(raw string) (float64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false
	}

	if parts := strings.Split(raw, "d"); len(parts) == 2 {
		days, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, false
		}
		h, ok := parseHMS(parts[1])
		if !ok {
			return 0, false
		}
		return float64(days*24) + h, true
	}

	h, ok := parseHMS(raw)
	if !ok {
		return 0, false
	}
	return h, true
}

func parseHMS(raw string) (float64, bool) {
	parts := strings.Split(raw, ":")
	if len(parts) < 2 {
		return 0, false
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, false
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, false
	}
	return float64(hours) + float64(minutes)/60, true
}
