package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"exchange-travel-planner/backend/internal/store"
)

func setup() (*Server, http.Handler) {
	s := NewServer(store.New())
	return s, s.Routes()
}

func TestHealthOK(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Fatalf("expected ok, got %s", body["status"])
	}
}

func TestHealthMethodNotAllowed(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 405 {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestCORSPreflight(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodOptions, "/health", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing CORS header")
	}
}

func TestTravelWindows(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/travel-windows", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]json.RawMessage
	json.NewDecoder(w.Body).Decode(&body)
	if _, ok := body["windows"]; !ok {
		t.Fatal("missing windows key")
	}
}

func TestTravelWindows_MethodNotAllowed(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodPost, "/api/travel-windows", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 405 {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestCalendarImport(t *testing.T) {
	_, h := setup()
	body := `{"events":[{"type":"exam","title":"Test","start":"2026-05-01","end":"2026-05-01","priority":5}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/calendar/import", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCalendarImport_BadJSON(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodPost, "/api/calendar/import", bytes.NewBufferString("{bad"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestTripOptimize(t *testing.T) {
	_, h := setup()
	body := `{"budgetCap":300,"maxTravelHours":6,"partySize":1,"style":"culture","departureCity":"Berlin"}`
	req := httptest.NewRequest(http.MethodPost, "/api/trips/optimize", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]json.RawMessage
	json.NewDecoder(w.Body).Decode(&resp)
	if _, ok := resp["options"]; !ok {
		t.Fatal("missing options")
	}
}

func TestGetTrip(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/trips/trip-1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetTrip_NotFound(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/trips/nonexistent", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestShareTrip(t *testing.T) {
	_, h := setup()
	body := `{"memberIds":["alice","bob"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/trips/trip-1/share", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestShareTrip_NotFound(t *testing.T) {
	_, h := setup()
	body := `{"memberIds":["alice"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/trips/nonexistent/share", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestBudgetEntries_Get(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/budget/entries?userId=demo-user", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestBudgetEntries_Post(t *testing.T) {
	_, h := setup()
	body := `{"category":"travel","amount":99,"date":"2026-03-01","note":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/budget/entries", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestBudgetEntries_PostMissingFields(t *testing.T) {
	_, h := setup()
	body := `{"category":"","amount":0,"date":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/budget/entries", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestBudgetForecast(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/budget/forecast?userId=demo-user&tripId=trip-1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSearchTransport(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/search/transport?from=Berlin&to=Prague", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSearchTransport_MissingTo(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/search/transport?from=Berlin", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSearchStays(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/search/stays?city=Prague", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSearchStays_MissingCity(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/search/stays", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestConflictsEvaluate(t *testing.T) {
	_, h := setup()
	body := `{"windowId":"w-1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/conflicts/evaluate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestConflictsEvaluate_MissingWindowID(t *testing.T) {
	_, h := setup()
	body := `{"windowId":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/conflicts/evaluate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestTripRoutes_ShortPath(t *testing.T) {
	_, h := setup()
	req := httptest.NewRequest(http.MethodGet, "/api/trips/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
