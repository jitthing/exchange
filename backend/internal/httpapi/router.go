package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"exchange-travel-planner/backend/internal/auth"
	"exchange-travel-planner/backend/internal/domain"
	"exchange-travel-planner/backend/internal/provider"
)

type Server struct {
	store             domain.DataStore
	transportProvider provider.TransportProvider
}

func NewServer(s domain.DataStore) *Server {
	return &Server{
		store:             s,
		transportProvider: provider.NewOpenTransportProviderFromEnv(),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", s.handleHealth)

	// Protected API routes — wrapped with RequireAuth
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/calendar/import", s.handleCalendarImport)
	apiMux.HandleFunc("/api/travel-windows", s.handleTravelWindows)
	apiMux.HandleFunc("/api/trips/optimize", s.handleTripOptimize)
	apiMux.HandleFunc("/api/trips/", s.handleTripRoutes)
	apiMux.HandleFunc("/api/budget/entries", s.handleBudgetEntries)
	apiMux.HandleFunc("/api/budget/forecast", s.handleBudgetForecast)
	apiMux.HandleFunc("/api/search/transport", s.handleSearchTransport)
	apiMux.HandleFunc("/api/search/stays", s.handleSearchStays)
	apiMux.HandleFunc("/api/conflicts/evaluate", s.handleConflicts)

	mux.Handle("/api/", auth.RequireAuth(apiMux))

	return corsMiddleware(mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleCalendarImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		Events []domain.AcademicEvent `json:"events"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": s.store.ImportAcademicEvents(req.Events)})
}

func (s *Server) handleTravelWindows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	writeJSON(w, http.StatusOK, map[string]any{"windows": s.store.ListTravelWindows(from, to)})
}

func (s *Server) handleTripOptimize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req domain.TripConstraint
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"options": s.store.OptimizeTrips(req)})
}

func (s *Server) handleTripRoutes(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeErr(w, http.StatusNotFound, "not found")
		return
	}
	tripID := parts[2]

	if len(parts) == 3 && r.Method == http.MethodGet {
		trip := s.store.GetTrip(tripID)
		if trip == nil {
			writeErr(w, http.StatusNotFound, "trip not found")
			return
		}
		writeJSON(w, http.StatusOK, trip)
		return
	}

	if len(parts) == 4 && parts[3] == "share" && r.Method == http.MethodPost {
		var req struct {
			MemberIDs []string `json:"memberIds"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid json")
			return
		}
		trip := s.store.ShareTrip(tripID, req.MemberIDs)
		if trip == nil {
			writeErr(w, http.StatusNotFound, "trip not found")
			return
		}
		writeJSON(w, http.StatusOK, trip)
		return
	}

	writeErr(w, http.StatusNotFound, "not found")
}

func (s *Server) handleBudgetEntries(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	if r.Method == http.MethodGet {
		writeJSON(w, http.StatusOK, map[string]any{"entries": s.store.ListBudgetEntries(userID)})
		return
	}

	if r.Method == http.MethodPost {
		var req domain.BudgetEntry
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid json")
			return
		}
		req.UserID = userID
		if req.Currency == "" {
			req.Currency = "EUR"
		}
		if req.Category == "" || req.Amount <= 0 || req.Date == "" {
			writeErr(w, http.StatusBadRequest, "missing required fields")
			return
		}
		entry := s.store.AddBudgetEntry(req)
		writeJSON(w, http.StatusCreated, entry)
		return
	}

	writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
}

func (s *Server) handleBudgetForecast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	userID := auth.UserIDFromContext(r.Context())
	tripID := r.URL.Query().Get("tripId")
	writeJSON(w, http.StatusOK, s.store.Forecast(userID, tripID))
}

func (s *Server) handleSearchTransport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if to == "" {
		writeErr(w, http.StatusBadRequest, "missing to")
		return
	}
	if s.transportProvider != nil {
		if options, err := s.transportProvider.SearchTransport(r.Context(), from, to); err == nil && len(options) > 0 {
			writeJSON(w, http.StatusOK, map[string]any{"options": options})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"options": s.store.SearchTransport(from, to)})
}

func (s *Server) handleSearchStays(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	city := r.URL.Query().Get("city")
	if city == "" {
		writeErr(w, http.StatusBadRequest, "missing city")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"options": s.store.SearchStays(city)})
}

func (s *Server) handleConflicts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		WindowID string `json:"windowId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.WindowID == "" {
		writeErr(w, http.StatusBadRequest, "missing windowId")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"alerts": s.store.EvaluateConflicts(req.WindowID)})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeErr(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
