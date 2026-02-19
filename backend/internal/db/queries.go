package db

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"gorm.io/gorm"

	"exchange-travel-planner/backend/internal/domain"
)

// PgStore implements domain.DataStore backed by PostgreSQL via GORM.
type PgStore struct {
	db *gorm.DB
}

func NewPgStore(db *gorm.DB) *PgStore {
	return &PgStore{db: db}
}

func (s *PgStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func makeID(prefix string) string {
	return fmt.Sprintf("%s-%06d", prefix, rand.Intn(999999))
}

func roundF(value float64, precision int) float64 {
	factor := math.Pow10(precision)
	return math.Round(value*factor) / factor
}

// --- Interface implementations ---

func (s *PgStore) ImportAcademicEvents(events []domain.AcademicEvent) []domain.AcademicEvent {
	for _, e := range events {
		if e.ID == "" {
			e.ID = makeID("ev")
		}
		m := AcademicEventModel{
			ID: e.ID, Type: string(e.Type), Title: e.Title,
			StartDate: e.Start, EndDate: e.End, Priority: e.Priority,
		}
		s.db.Where("id = ?", m.ID).FirstOrCreate(&m)
	}
	var models []AcademicEventModel
	s.db.Order("start_date").Find(&models)
	result := make([]domain.AcademicEvent, len(models))
	for i, m := range models {
		result[i] = domain.AcademicEvent{
			ID: m.ID, Type: domain.AcademicEventType(m.Type), Title: m.Title,
			Start: m.StartDate, End: m.EndDate, Priority: m.Priority,
		}
	}
	return result
}

func (s *PgStore) ListTravelWindows(from, to string) []domain.TravelWindow {
	q := s.db.Model(&TravelWindowModel{})
	if from != "" {
		q = q.Where("end_date >= ?", from)
	}
	if to != "" {
		q = q.Where("start_date <= ?", to)
	}
	var models []TravelWindowModel
	q.Order("start_date").Find(&models)
	result := make([]domain.TravelWindow, len(models))
	for i, m := range models {
		conflicts := []string(m.Conflicts)
		if conflicts == nil {
			conflicts = []string{}
		}
		result[i] = domain.TravelWindow{
			ID: m.ID, StartDate: m.StartDate, EndDate: m.EndDate,
			Score: m.Score, Conflicts: conflicts,
		}
	}
	return result
}

func (s *PgStore) loadDestinations() []DestinationModel {
	var dests []DestinationModel
	s.db.Find(&dests)
	return dests
}

func (s *PgStore) OptimizeTrips(c domain.TripConstraint) []domain.TripOption {
	dests := s.loadDestinations()

	type scored struct {
		domain.TripOption
		Score float64
	}
	items := make([]scored, 0, len(dests))
	for _, entry := range dests {
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
			{Provider: "EuroRail Connect", Mode: "train", DurationHours: roundF(entry.BaseTravelHrs, 1), Price: math.Round(transportPrice * 0.92), Deeplink: "https://example.com/train"},
			{Provider: "BudgetBus Europe", Mode: "bus", DurationHours: roundF(entry.BaseTravelHrs*1.3, 1), Price: math.Round(transportPrice * 0.76), Deeplink: "https://example.com/bus"},
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

func (s *PgStore) GetTrip(id string) *domain.Trip {
	var m TripModel
	if err := s.db.First(&m, "id = ?", id).Error; err != nil {
		return nil
	}
	return &domain.Trip{
		ID: m.ID, OwnerID: m.OwnerID, Destination: m.Destination,
		WindowID: m.WindowID, Members: []string(m.Members),
		Itinerary: []string(m.Itinerary), EstimatedCost: m.EstimatedCost,
	}
}

func (s *PgStore) ShareTrip(tripID string, memberIDs []string) *domain.Trip {
	t := s.GetTrip(tripID)
	if t == nil {
		return nil
	}
	existing := map[string]struct{}{}
	for _, m := range t.Members {
		existing[m] = struct{}{}
	}
	for _, m := range memberIDs {
		if m != "" {
			existing[m] = struct{}{}
		}
	}
	merged := make([]string, 0, len(existing))
	for m := range existing {
		merged = append(merged, m)
	}
	sort.Strings(merged)
	s.db.Model(&TripModel{}).Where("id = ?", tripID).Update("members", JSONStringSlice(merged))
	t.Members = merged
	return t
}

func (s *PgStore) AddBudgetEntry(entry domain.BudgetEntry) domain.BudgetEntry {
	entry.ID = makeID("b")
	m := BudgetEntryModel{
		ID: entry.ID, UserID: entry.UserID, Category: entry.Category,
		Amount: entry.Amount, Currency: entry.Currency, Date: entry.Date,
		TripID: entry.TripID, Note: entry.Note,
	}
	s.db.Create(&m)
	return entry
}

func (s *PgStore) ListBudgetEntries(userID string) []domain.BudgetEntry {
	var models []BudgetEntryModel
	s.db.Where("user_id = ?", userID).Order("date").Find(&models)
	result := make([]domain.BudgetEntry, len(models))
	for i, m := range models {
		result[i] = domain.BudgetEntry{
			ID: m.ID, UserID: m.UserID, Category: m.Category,
			Amount: m.Amount, Currency: m.Currency, Date: m.Date,
			TripID: m.TripID, Note: m.Note,
		}
	}
	return result
}

func (s *PgStore) Forecast(userID, tripID string) domain.ForecastResult {
	entries := s.ListBudgetEntries(userID)
	spend := 0.0
	for _, e := range entries {
		spend += e.Amount
	}

	var mb MonthlyBudgetModel
	budget := 900.0
	if err := s.db.First(&mb, "user_id = ?", userID).Error; err == nil && mb.Budget > 0 {
		budget = mb.Budget
	}

	tripCost := 0.0
	if tripID != "" {
		if t := s.GetTrip(tripID); t != nil {
			tripCost = t.EstimatedCost
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
		ProjectedMonthlySpend: roundF(projected, 2),
		RemainingBudget:       roundF(remaining, 2),
		Affordability:         affordability,
	}
}

func (s *PgStore) EvaluateConflicts(windowID string) []domain.ConflictAlert {
	var w TravelWindowModel
	if err := s.db.First(&w, "id = ?", windowID).Error; err != nil {
		return []domain.ConflictAlert{}
	}
	start, _ := time.Parse("2006-01-02", w.StartDate)
	end, _ := time.Parse("2006-01-02", w.EndDate)

	var events []AcademicEventModel
	s.db.Find(&events)

	var alerts []domain.ConflictAlert
	for _, ev := range events {
		eventDate, err := time.Parse("2006-01-02", ev.StartDate)
		if err != nil || eventDate.Before(start) || eventDate.After(end) {
			continue
		}
		sev := domain.SeverityInfo
		if ev.Type == "exam" {
			sev = domain.SeverityHighRisk
		} else if ev.Type == "deadline" {
			sev = domain.SeverityWarning
		}
		alerts = append(alerts, domain.ConflictAlert{
			Severity:       sev,
			Reason:         ev.Type + " overlap: " + ev.Title,
			RelatedEventID: ev.ID,
		})
	}
	if alerts == nil {
		alerts = []domain.ConflictAlert{}
	}
	return alerts
}

func (s *PgStore) SearchTransport(_from, to string) []domain.TransportOption {
	var dest DestinationModel
	if err := s.db.First(&dest, "LOWER(city) = LOWER(?)", to).Error; err != nil {
		return []domain.TransportOption{}
	}
	return []domain.TransportOption{
		{Provider: "EuroRail Connect", Mode: "train", DurationHours: dest.BaseTravelHrs, Price: math.Round(dest.TransportBase * 0.9), Deeplink: "https://example.com/train"},
		{Provider: "SkySaver", Mode: "flight", DurationHours: roundF(dest.BaseTravelHrs*0.65, 1), Price: math.Round(dest.TransportBase * 1.18), Deeplink: "https://example.com/flight"},
	}
}

func (s *PgStore) SearchStays(city string) []domain.StayOption {
	var dest DestinationModel
	if err := s.db.First(&dest, "LOWER(city) = LOWER(?)", city).Error; err != nil {
		return []domain.StayOption{}
	}
	return []domain.StayOption{
		{Provider: "HostelGraph", Kind: "hostel", NightlyPrice: dest.HostelNightEUR, Rating: 4.2, Deeplink: "https://example.com/hostel"},
		{Provider: "StudentStay", Kind: "budget-hotel", NightlyPrice: dest.HostelNightEUR + 14, Rating: 4.0, Deeplink: "https://example.com/hotel"},
	}
}
