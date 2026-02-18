package store

import (
	"testing"

	"exchange-travel-planner/backend/internal/domain"
)

func TestListTravelWindows_All(t *testing.T) {
	s := New()
	windows := s.ListTravelWindows("", "")
	if len(windows) != 3 {
		t.Fatalf("expected 3 windows, got %d", len(windows))
	}
}

func TestListTravelWindows_Filtered(t *testing.T) {
	s := New()
	windows := s.ListTravelWindows("2026-03-01", "2026-03-10")
	if len(windows) != 1 {
		t.Fatalf("expected 1 window in range, got %d", len(windows))
	}
	if windows[0].ID != "w-1" {
		t.Fatalf("expected w-1, got %s", windows[0].ID)
	}
}

func TestListTravelWindows_FromOnly(t *testing.T) {
	s := New()
	windows := s.ListTravelWindows("2026-04-01", "")
	if len(windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(windows))
	}
}

func TestOptimizeTrips(t *testing.T) {
	s := New()
	opts := s.OptimizeTrips(domain.TripConstraint{
		BudgetCap:      300,
		MaxTravelHours: 6,
		PartySize:      1,
		Style:          "culture",
		DepartureCity:  "Berlin",
	})
	if len(opts) != 4 {
		t.Fatalf("expected 4 options, got %d", len(opts))
	}
	// Culture style should boost Prague and Krakow
	for _, o := range opts {
		if o.Destination == "" {
			t.Fatal("empty destination")
		}
		if o.TotalEstimatedCost <= 0 {
			t.Fatalf("cost should be positive for %s", o.Destination)
		}
		if len(o.TransportOptions) != 2 {
			t.Fatalf("expected 2 transport options for %s", o.Destination)
		}
		if len(o.StayOptions) != 2 {
			t.Fatalf("expected 2 stay options for %s", o.Destination)
		}
	}
}

func TestOptimizeTrips_OverBudget(t *testing.T) {
	s := New()
	opts := s.OptimizeTrips(domain.TripConstraint{
		BudgetCap:      10,
		MaxTravelHours: 6,
		PartySize:      1,
		Style:          "city",
	})
	for _, o := range opts {
		if o.TotalEstimatedCost <= 10 {
			continue
		}
		if o.RiskLevel != domain.SeverityWarning {
			t.Fatalf("expected warning risk for over-budget %s, got %s", o.Destination, o.RiskLevel)
		}
	}
}

func TestGetTrip_Found(t *testing.T) {
	s := New()
	trip := s.GetTrip("trip-1")
	if trip == nil {
		t.Fatal("expected trip-1")
	}
	if trip.Destination != "Prague" {
		t.Fatalf("expected Prague, got %s", trip.Destination)
	}
}

func TestGetTrip_NotFound(t *testing.T) {
	s := New()
	trip := s.GetTrip("nonexistent")
	if trip != nil {
		t.Fatal("expected nil for nonexistent trip")
	}
}

func TestShareTrip(t *testing.T) {
	s := New()
	trip := s.ShareTrip("trip-1", []string{"alice", "bob"})
	if trip == nil {
		t.Fatal("expected trip")
	}
	if len(trip.Members) != 3 { // demo-user + alice + bob
		t.Fatalf("expected 3 members, got %d", len(trip.Members))
	}
}

func TestShareTrip_NotFound(t *testing.T) {
	s := New()
	trip := s.ShareTrip("nonexistent", []string{"alice"})
	if trip != nil {
		t.Fatal("expected nil")
	}
}

func TestShareTrip_EmptyMember(t *testing.T) {
	s := New()
	trip := s.ShareTrip("trip-1", []string{""})
	if trip == nil {
		t.Fatal("expected trip")
	}
	// Empty string should be filtered out
	if len(trip.Members) != 1 {
		t.Fatalf("expected 1 member, got %d: %v", len(trip.Members), trip.Members)
	}
}

func TestShareTrip_Dedup(t *testing.T) {
	s := New()
	trip := s.ShareTrip("trip-1", []string{"demo-user", "demo-user"})
	if trip == nil {
		t.Fatal("expected trip")
	}
	if len(trip.Members) != 1 {
		t.Fatalf("expected 1 member after dedup, got %d", len(trip.Members))
	}
}

func TestAddBudgetEntry(t *testing.T) {
	s := New()
	entry := s.AddBudgetEntry(domain.BudgetEntry{
		UserID:   "demo-user",
		Category: "travel",
		Amount:   50,
		Currency: "EUR",
		Date:     "2026-03-01",
	})
	if entry.ID == "" {
		t.Fatal("expected generated ID")
	}
	if entry.Amount != 50 {
		t.Fatalf("expected 50, got %f", entry.Amount)
	}
}

func TestListBudgetEntries(t *testing.T) {
	s := New()
	entries := s.ListBudgetEntries("demo-user")
	if len(entries) != 2 {
		t.Fatalf("expected 2 seed entries, got %d", len(entries))
	}
}

func TestListBudgetEntries_Empty(t *testing.T) {
	s := New()
	entries := s.ListBudgetEntries("unknown-user")
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestForecast_WithoutTrip(t *testing.T) {
	s := New()
	f := s.Forecast("demo-user", "")
	if f.ProjectedMonthlySpend != 480 {
		t.Fatalf("expected 480, got %f", f.ProjectedMonthlySpend)
	}
	if f.RemainingBudget != 420 {
		t.Fatalf("expected 420, got %f", f.RemainingBudget)
	}
	if f.Affordability != "green" {
		t.Fatalf("expected green, got %s", f.Affordability)
	}
}

func TestForecast_WithTrip(t *testing.T) {
	s := New()
	f := s.Forecast("demo-user", "trip-1")
	// 480 + 220 = 700, remaining = 200, which is not < 200 so green
	if f.ProjectedMonthlySpend != 700 {
		t.Fatalf("expected 700, got %f", f.ProjectedMonthlySpend)
	}
	if f.Affordability != "green" {
		t.Fatalf("expected green, got %s", f.Affordability)
	}
}

func TestForecast_Amber(t *testing.T) {
	s := New()
	s.AddBudgetEntry(domain.BudgetEntry{
		UserID: "demo-user", Category: "living", Amount: 1, Currency: "EUR", Date: "2026-02-10",
	})
	f := s.Forecast("demo-user", "trip-1")
	// 481 + 220 = 701, remaining = 199 < 200
	if f.Affordability != "amber" {
		t.Fatalf("expected amber, got %s", f.Affordability)
	}
}

func TestForecast_Red(t *testing.T) {
	s := New()
	// Add a big entry to push over budget
	s.AddBudgetEntry(domain.BudgetEntry{
		UserID: "demo-user", Category: "living", Amount: 500, Currency: "EUR", Date: "2026-02-10",
	})
	f := s.Forecast("demo-user", "trip-1")
	if f.Affordability != "red" {
		t.Fatalf("expected red, got %s", f.Affordability)
	}
}

func TestForecast_UnknownUser(t *testing.T) {
	s := New()
	f := s.Forecast("nobody", "")
	if f.ProjectedMonthlySpend != 0 {
		t.Fatalf("expected 0, got %f", f.ProjectedMonthlySpend)
	}
	if f.Affordability != "green" {
		t.Fatalf("expected green, got %s", f.Affordability)
	}
}

func TestEvaluateConflicts_NoConflict(t *testing.T) {
	s := New()
	alerts := s.EvaluateConflicts("w-1") // Mar 6-8, no events overlap
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestEvaluateConflicts_WithConflict(t *testing.T) {
	s := New()
	// w-2 is Mar 20-22, deadline is Mar 24 - no overlap
	// Let's add an event that overlaps with w-2
	s.ImportAcademicEvents([]domain.AcademicEvent{
		{Type: domain.AcademicExam, Title: "Test Exam", Start: "2026-03-21", End: "2026-03-21", Priority: 5},
	})
	alerts := s.EvaluateConflicts("w-2")
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Severity != domain.SeverityHighRisk {
		t.Fatalf("expected high-risk, got %s", alerts[0].Severity)
	}
}

func TestEvaluateConflicts_UnknownWindow(t *testing.T) {
	s := New()
	alerts := s.EvaluateConflicts("nonexistent")
	if len(alerts) != 0 {
		t.Fatalf("expected 0, got %d", len(alerts))
	}
}

func TestSearchTransport_Known(t *testing.T) {
	s := New()
	opts := s.SearchTransport("Berlin", "Prague")
	if len(opts) != 2 {
		t.Fatalf("expected 2, got %d", len(opts))
	}
}

func TestSearchTransport_Unknown(t *testing.T) {
	s := New()
	opts := s.SearchTransport("Berlin", "Tokyo")
	if len(opts) != 0 {
		t.Fatalf("expected 0, got %d", len(opts))
	}
}

func TestSearchStays_Known(t *testing.T) {
	s := New()
	opts := s.SearchStays("Budapest")
	if len(opts) != 2 {
		t.Fatalf("expected 2, got %d", len(opts))
	}
}

func TestSearchStays_Unknown(t *testing.T) {
	s := New()
	opts := s.SearchStays("Tokyo")
	if len(opts) != 0 {
		t.Fatalf("expected 0, got %d", len(opts))
	}
}

func TestImportAcademicEvents(t *testing.T) {
	s := New()
	events := s.ImportAcademicEvents([]domain.AcademicEvent{
		{Type: domain.AcademicClass, Title: "Math", Start: "2026-03-01", End: "2026-03-01", Priority: 2},
	})
	if len(events) != 4 { // 3 seed + 1 new
		t.Fatalf("expected 4, got %d", len(events))
	}
	// The new one should have an auto-generated ID
	last := events[len(events)-1]
	if last.ID == "" {
		t.Fatal("expected auto-generated ID")
	}
}

func TestImportAcademicEvents_WithID(t *testing.T) {
	s := New()
	events := s.ImportAcademicEvents([]domain.AcademicEvent{
		{ID: "custom-id", Type: domain.AcademicExam, Title: "Physics", Start: "2026-04-01", End: "2026-04-01", Priority: 5},
	})
	last := events[len(events)-1]
	if last.ID != "custom-id" {
		t.Fatalf("expected custom-id, got %s", last.ID)
	}
}
