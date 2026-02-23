package provider

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestParseDurationHours(t *testing.T) {
	cases := []struct {
		in       string
		want     float64
		wantOkay bool
	}{
		{in: "02:30:00", want: 2.5, wantOkay: true},
		{in: "1d03:15:00", want: 27.25, wantOkay: true},
		{in: "", wantOkay: false},
		{in: "bad", wantOkay: false},
	}

	for _, tc := range cases {
		got, ok := parseDurationHours(tc.in)
		if ok != tc.wantOkay {
			t.Fatalf("parseDurationHours(%q) ok=%v, want %v", tc.in, ok, tc.wantOkay)
		}
		if ok && got != tc.want {
			t.Fatalf("parseDurationHours(%q)=%v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestSearchTransport_MapsProviderPayload(t *testing.T) {
	p := &OpenTransportProvider{
		enabled: true,
		baseURL: "https://provider.example.test",
		client: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path != "/connections" {
					t.Fatalf("unexpected path %s", r.URL.Path)
				}
				body := `{"connections":[{"duration":"02:10:00","products":["IC"],"transfers":1},{"duration":"03:00:00","products":["BUS"],"transfers":0}]}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}),
		},
	}

	opts, err := p.SearchTransport(context.Background(), "Berlin", "Prague")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts) != 2 {
		t.Fatalf("expected 2 options, got %d", len(opts))
	}
	if opts[0].Provider != "OpenTransportData" {
		t.Fatalf("unexpected provider %s", opts[0].Provider)
	}
	if opts[0].Mode != "train" {
		t.Fatalf("expected train mode, got %s", opts[0].Mode)
	}
	if opts[1].Mode != "bus" {
		t.Fatalf("expected bus mode, got %s", opts[1].Mode)
	}
}

func TestSearchTransport_Disabled(t *testing.T) {
	p := &OpenTransportProvider{enabled: false}
	if _, err := p.SearchTransport(context.Background(), "Berlin", "Prague"); err == nil {
		t.Fatal("expected disabled error")
	}
}
