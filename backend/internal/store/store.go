package store

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"exchange-travel-planner/backend/internal/domain"
)

type destinationSeed struct {
	City           string
	BaseTravelHrs  float64
	TransportBase  float64
	HostelNightEUR float64
	Tags           []string
}

type Store struct {
	mu sync.RWMutex

	academicEvents []domain.AcademicEvent
	travelWindows  []domain.TravelWindow
	trips          []domain.Trip
	budgetEntries  []domain.BudgetEntry
	monthlyBudget  map[string]float64
	destinations   []destinationSeed
}

func New() *Store {
	return &Store{
		academicEvents: []domain.AcademicEvent{
			{ID: "ev-1", Type: domain.AcademicExam, Title: "Economics Midterm", Start: "2026-03-18", End: "2026-03-18", Priority: 5},
			{ID: "ev-2", Type: domain.AcademicDeadline, Title: "Group Project Deadline", Start: "2026-03-24", End: "2026-03-24", Priority: 4},
			{ID: "ev-3", Type: domain.AcademicHoliday, Title: "Public Holiday", Start: "2026-04-03", End: "2026-04-05", Priority: 1},
		},
		travelWindows: []domain.TravelWindow{
			{ID: "w-1", StartDate: "2026-03-06", EndDate: "2026-03-08", Score: 88, Conflicts: []string{}},
			{ID: "w-2", StartDate: "2026-03-20", EndDate: "2026-03-22", Score: 52, Conflicts: []string{"Near major deadline"}},
			{ID: "w-3", StartDate: "2026-04-03", EndDate: "2026-04-06", Score: 95, Conflicts: []string{}},
		},
		trips: []domain.Trip{
			{ID: "trip-1", OwnerID: "demo-user", Destination: "Prague", WindowID: "w-1", Members: []string{"demo-user"}, Itinerary: []string{"Old Town walk", "Charles Bridge sunrise"}, EstimatedCost: 220},
		},
		budgetEntries: []domain.BudgetEntry{
			{ID: "b-1", UserID: "demo-user", Category: "living", Amount: 420, Currency: "EUR", Date: "2026-02-05", Note: "Rent split"},
			{ID: "b-2", UserID: "demo-user", Category: "travel", Amount: 60, Currency: "EUR", Date: "2026-02-08", Note: "Train to Vienna"},
		},
		monthlyBudget: map[string]float64{"demo-user": 900},
		destinations: []destinationSeed{
			{City: "Prague", BaseTravelHrs: 3.8, TransportBase: 55, HostelNightEUR: 28, Tags: []string{"culture", "city"}},
			{City: "Budapest", BaseTravelHrs: 4.7, TransportBase: 47, HostelNightEUR: 24, Tags: []string{"nightlife", "city"}},
			{City: "Ljubljana", BaseTravelHrs: 5.2, TransportBase: 41, HostelNightEUR: 30, Tags: []string{"nature", "city"}},
			{City: "Krakow", BaseTravelHrs: 2.9, TransportBase: 50, HostelNightEUR: 22, Tags: []string{"culture", "city"}},
		},
	}
}

func makeID(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s-%06d", prefix, rand.Intn(999999))
}

func (s *Store) ImportAcademicEvents(events []domain.AcademicEvent) []domain.AcademicEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, event := range events {
		if event.ID == "" {
			event.ID = makeID("ev")
		}
		s.academicEvents = append(s.academicEvents, event)
	}
	return append([]domain.AcademicEvent(nil), s.academicEvents...)
}

func (s *Store) ListTravelWindows(from, to string) []domain.TravelWindow {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if from == "" && to == "" {
		return append([]domain.TravelWindow(nil), s.travelWindows...)
	}
	res := make([]domain.TravelWindow, 0)
	for _, window := range s.travelWindows {
		afterFrom := from == "" || window.EndDate >= from
		beforeTo := to == "" || window.StartDate <= to
		if afterFrom && beforeTo {
			res = append(res, window)
		}
	}
	return res
}

func (s *Store) OptimizeTrips(c domain.TripConstraint) []domain.TripOption {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type scored struct {
		domain.TripOption
		Score float64
	}
	items := make([]scored, 0, len(s.destinations))
	for _, entry := range s.destinations {
		styleBoost := 0.0
		for _, tag := range entry.Tags {
			if tag == c.Style {
				styleBoost = 14
				break
			}
		}
		durationPenalty := math.Max(0, entry.BaseTravelHrs-c.MaxTravelHours) * 18
		transportPrice := entry.TransportBase + float64(c.PartySize*7)
		stayPrice := entry.HostelNightEUR * 2 * float64(c.PartySize)
		total := math.Round(transportPrice + stayPrice)
		budgetPenalty := 0.0
		if total > c.BudgetCap {
			budgetPenalty = (total - c.BudgetCap) / 3
		}
		score := 100 + styleBoost - durationPenalty - budgetPenalty

		reasons := make([]string, 0, 3)
		if entry.BaseTravelHrs <= c.MaxTravelHours {
			reasons = append(reasons, "short-transit")
		}
		if total <= c.BudgetCap {
			reasons = append(reasons, "within-budget")
		}
		if styleBoost > 0 {
			reasons = append(reasons, "style-match")
		}
		if len(reasons) == 0 {
			reasons = []string{"stretch-choice"}
		}

		transport := []domain.TransportOption{
			{Provider: "EuroRail Connect", Mode: "train", DurationHours: round(entry.BaseTravelHrs, 1), Price: math.Round(transportPrice * 0.92), Deeplink: "https://example.com/train"},
			{Provider: "BudgetBus Europe", Mode: "bus", DurationHours: round(entry.BaseTravelHrs*1.3, 1), Price: math.Round(transportPrice * 0.76), Deeplink: "https://example.com/bus"},
		}
		stays := []domain.StayOption{
			{Provider: "HostelGraph", Kind: "hostel", NightlyPrice: entry.HostelNightEUR, Rating: 4.3, Deeplink: "https://example.com/hostel"},
			{Provider: "StudentStay", Kind: "budget-hotel", NightlyPrice: entry.HostelNightEUR + 16, Rating: 4.0, Deeplink: "https://example.com/hotel"},
		}

		risk := domain.SeverityInfo
		if total > c.BudgetCap {
			risk = domain.SeverityWarning
		}

		items = append(items, scored{
			TripOption: domain.TripOption{
				ID:                 makeID("opt"),
				Destination:        entry.City,
				ReasonTags:         reasons,
				TotalEstimatedCost: total,
				TransportOptions:   transport,
				StayOptions:        stays,
				RiskLevel:          risk,
			},
			Score: score,
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Score > items[j].Score })
	out := make([]domain.TripOption, 0, len(items))
	for _, item := range items {
		out = append(out, item.TripOption)
	}
	return out
}

func (s *Store) GetTrip(id string) *domain.Trip {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, trip := range s.trips {
		if trip.ID == id {
			cp := trip
			return &cp
		}
	}
	return nil
}

func (s *Store) ShareTrip(tripID string, memberIDs []string) *domain.Trip {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.trips {
		if s.trips[i].ID == tripID {
			existing := map[string]struct{}{}
			for _, member := range s.trips[i].Members {
				existing[member] = struct{}{}
			}
			for _, member := range memberIDs {
				if member == "" {
					continue
				}
				existing[member] = struct{}{}
			}
			merged := make([]string, 0, len(existing))
			for member := range existing {
				merged = append(merged, member)
			}
			sort.Strings(merged)
			s.trips[i].Members = merged
			cp := s.trips[i]
			return &cp
		}
	}
	return nil
}

func (s *Store) AddBudgetEntry(entry domain.BudgetEntry) domain.BudgetEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry.ID = makeID("b")
	s.budgetEntries = append(s.budgetEntries, entry)
	return entry
}

func (s *Store) ListBudgetEntries(userID string) []domain.BudgetEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]domain.BudgetEntry, 0)
	for _, entry := range s.budgetEntries {
		if entry.UserID == userID {
			res = append(res, entry)
		}
	}
	return res
}

func (s *Store) Forecast(userID, tripID string) domain.ForecastResult {
	entries := s.ListBudgetEntries(userID)
	spend := 0.0
	for _, entry := range entries {
		spend += entry.Amount
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	budget := s.monthlyBudget[userID]
	if budget == 0 {
		budget = 900
	}
	tripCost := 0.0
	if tripID != "" {
		for _, trip := range s.trips {
			if trip.ID == tripID {
				tripCost = trip.EstimatedCost
				break
			}
		}
	}
	projected := spend + tripCost
	remaining := budget - projected
	affordability := "green"
	if remaining < 0 {
		affordability = "red"
	} else if remaining < 200 {
		affordability = "amber"
	}
	return domain.ForecastResult{
		ProjectedMonthlySpend: round(projected, 2),
		RemainingBudget:       round(remaining, 2),
		Affordability:         affordability,
	}
}

func (s *Store) EvaluateConflicts(windowID string) []domain.ConflictAlert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var target *domain.TravelWindow
	for _, window := range s.travelWindows {
		if window.ID == windowID {
			cp := window
			target = &cp
			break
		}
	}
	if target == nil {
		return []domain.ConflictAlert{}
	}
	start, _ := time.Parse("2006-01-02", target.StartDate)
	end, _ := time.Parse("2006-01-02", target.EndDate)
	alerts := make([]domain.ConflictAlert, 0)
	for _, event := range s.academicEvents {
		eventDate, err := time.Parse("2006-01-02", event.Start)
		if err != nil {
			continue
		}
		if eventDate.Before(start) || eventDate.After(end) {
			continue
		}
		sev := domain.SeverityInfo
		if event.Type == domain.AcademicExam {
			sev = domain.SeverityHighRisk
		} else if event.Type == domain.AcademicDeadline {
			sev = domain.SeverityWarning
		}
		alerts = append(alerts, domain.ConflictAlert{
			Severity:       sev,
			Reason:         string(event.Type) + " overlap: " + event.Title,
			RelatedEventID: event.ID,
		})
	}
	return alerts
}

func (s *Store) SearchTransport(_from, to string) []domain.TransportOption {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, entry := range s.destinations {
		if strings.EqualFold(entry.City, to) {
			return []domain.TransportOption{
				{Provider: "EuroRail Connect", Mode: "train", DurationHours: entry.BaseTravelHrs, Price: math.Round(entry.TransportBase * 0.9), Deeplink: "https://example.com/train"},
				{Provider: "SkySaver", Mode: "flight", DurationHours: round(entry.BaseTravelHrs*0.65, 1), Price: math.Round(entry.TransportBase * 1.18), Deeplink: "https://example.com/flight"},
			}
		}
	}
	return []domain.TransportOption{}
}

func (s *Store) SearchStays(city string) []domain.StayOption {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, entry := range s.destinations {
		if strings.EqualFold(entry.City, city) {
			return []domain.StayOption{
				{Provider: "HostelGraph", Kind: "hostel", NightlyPrice: entry.HostelNightEUR, Rating: 4.2, Deeplink: "https://example.com/hostel"},
				{Provider: "StudentStay", Kind: "budget-hotel", NightlyPrice: entry.HostelNightEUR + 14, Rating: 4.0, Deeplink: "https://example.com/hotel"},
			}
		}
	}
	return []domain.StayOption{}
}

func round(value float64, precision int) float64 {
	factor := math.Pow10(precision)
	return math.Round(value*factor) / factor
}
