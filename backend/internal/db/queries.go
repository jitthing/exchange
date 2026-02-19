package db

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"exchange-travel-planner/backend/internal/domain"
)

// PgStore implements domain.DataStore backed by PostgreSQL.
type PgStore struct {
	pool *pgxpool.Pool
}

func NewPgStore(pool *pgxpool.Pool) *PgStore {
	return &PgStore{pool: pool}
}

func (s *PgStore) Close() error {
	s.pool.Close()
	return nil
}

func makeID(prefix string) string {
	return fmt.Sprintf("%s-%06d", prefix, rand.Intn(999999))
}

func (s *PgStore) ImportAcademicEvents(events []domain.AcademicEvent) []domain.AcademicEvent {
	ctx := context.Background()
	for _, e := range events {
		if e.ID == "" {
			e.ID = makeID("ev")
		}
		_, _ = s.pool.Exec(ctx,
			`INSERT INTO academic_events (id, type, title, start_date, end_date, priority)
			 VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (id) DO NOTHING`,
			e.ID, string(e.Type), e.Title, e.Start, e.End, e.Priority)
	}
	rows, err := s.pool.Query(ctx, `SELECT id, type, title, start_date, end_date, priority FROM academic_events ORDER BY start_date`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []domain.AcademicEvent
	for rows.Next() {
		var ev domain.AcademicEvent
		var t string
		_ = rows.Scan(&ev.ID, &t, &ev.Title, &ev.Start, &ev.End, &ev.Priority)
		ev.Type = domain.AcademicEventType(t)
		result = append(result, ev)
	}
	return result
}

func (s *PgStore) ListTravelWindows(from, to string) []domain.TravelWindow {
	ctx := context.Background()
	query := `SELECT id, start_date, end_date, score, conflicts FROM travel_windows WHERE 1=1`
	args := []any{}
	n := 0
	if from != "" {
		n++
		query += fmt.Sprintf(` AND end_date >= $%d`, n)
		args = append(args, from)
	}
	if to != "" {
		n++
		query += fmt.Sprintf(` AND start_date <= $%d`, n)
		args = append(args, to)
	}
	query += ` ORDER BY start_date`
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []domain.TravelWindow
	for rows.Next() {
		var w domain.TravelWindow
		var conflictsJSON []byte
		_ = rows.Scan(&w.ID, &w.StartDate, &w.EndDate, &w.Score, &conflictsJSON)
		_ = json.Unmarshal(conflictsJSON, &w.Conflicts)
		if w.Conflicts == nil {
			w.Conflicts = []string{}
		}
		result = append(result, w)
	}
	if result == nil {
		result = []domain.TravelWindow{}
	}
	return result
}

type destRow struct {
	City           string
	BaseTravelHrs  float64
	TransportBase  float64
	HostelNightEUR float64
	Tags           []string
}

func (s *PgStore) loadDestinations() []destRow {
	ctx := context.Background()
	rows, err := s.pool.Query(ctx, `SELECT city, base_travel_hrs, transport_base, hostel_night_eur, tags FROM destinations`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []destRow
	for rows.Next() {
		var d destRow
		var tagsJSON []byte
		_ = rows.Scan(&d.City, &d.BaseTravelHrs, &d.TransportBase, &d.HostelNightEUR, &tagsJSON)
		_ = json.Unmarshal(tagsJSON, &d.Tags)
		out = append(out, d)
	}
	return out
}

func roundF(value float64, precision int) float64 {
	factor := math.Pow10(precision)
	return math.Round(value*factor) / factor
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
	ctx := context.Background()
	var t domain.Trip
	var membersJSON, itineraryJSON []byte
	err := s.pool.QueryRow(ctx,
		`SELECT id, owner_id, destination, window_id, members, itinerary, estimated_cost FROM trips WHERE id=$1`, id).
		Scan(&t.ID, &t.OwnerID, &t.Destination, &t.WindowID, &membersJSON, &itineraryJSON, &t.EstimatedCost)
	if err != nil {
		return nil
	}
	_ = json.Unmarshal(membersJSON, &t.Members)
	_ = json.Unmarshal(itineraryJSON, &t.Itinerary)
	return &t
}

func (s *PgStore) ShareTrip(tripID string, memberIDs []string) *domain.Trip {
	ctx := context.Background()
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
	membersJSON, _ := json.Marshal(merged)
	_, _ = s.pool.Exec(ctx, `UPDATE trips SET members=$1 WHERE id=$2`, membersJSON, tripID)
	t.Members = merged
	return t
}

func (s *PgStore) AddBudgetEntry(entry domain.BudgetEntry) domain.BudgetEntry {
	ctx := context.Background()
	entry.ID = makeID("b")
	_, _ = s.pool.Exec(ctx,
		`INSERT INTO budget_entries (id, user_id, category, amount, currency, date, trip_id, note)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		entry.ID, entry.UserID, entry.Category, entry.Amount, entry.Currency, entry.Date, entry.TripID, entry.Note)
	return entry
}

func (s *PgStore) ListBudgetEntries(userID string) []domain.BudgetEntry {
	ctx := context.Background()
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, category, amount, currency, date, COALESCE(trip_id,''), COALESCE(note,'') FROM budget_entries WHERE user_id=$1 ORDER BY date`, userID)
	if err != nil {
		return []domain.BudgetEntry{}
	}
	defer rows.Close()
	var result []domain.BudgetEntry
	for rows.Next() {
		var e domain.BudgetEntry
		_ = rows.Scan(&e.ID, &e.UserID, &e.Category, &e.Amount, &e.Currency, &e.Date, &e.TripID, &e.Note)
		result = append(result, e)
	}
	if result == nil {
		result = []domain.BudgetEntry{}
	}
	return result
}

func (s *PgStore) Forecast(userID, tripID string) domain.ForecastResult {
	entries := s.ListBudgetEntries(userID)
	spend := 0.0
	for _, e := range entries {
		spend += e.Amount
	}

	ctx := context.Background()
	var budget float64
	err := s.pool.QueryRow(ctx, `SELECT budget FROM monthly_budgets WHERE user_id=$1`, userID).Scan(&budget)
	if err != nil || budget == 0 {
		budget = 900
	}

	tripCost := 0.0
	if tripID != "" {
		t := s.GetTrip(tripID)
		if t != nil {
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
	ctx := context.Background()
	var startDate, endDate string
	err := s.pool.QueryRow(ctx, `SELECT start_date, end_date FROM travel_windows WHERE id=$1`, windowID).Scan(&startDate, &endDate)
	if err != nil {
		return []domain.ConflictAlert{}
	}
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	rows, err := s.pool.Query(ctx, `SELECT id, type, title, start_date FROM academic_events`)
	if err != nil {
		return []domain.ConflictAlert{}
	}
	defer rows.Close()

	var alerts []domain.ConflictAlert
	for rows.Next() {
		var id, typ, title, evStart string
		_ = rows.Scan(&id, &typ, &title, &evStart)
		eventDate, err := time.Parse("2006-01-02", evStart)
		if err != nil || eventDate.Before(start) || eventDate.After(end) {
			continue
		}
		sev := domain.SeverityInfo
		if typ == "exam" {
			sev = domain.SeverityHighRisk
		} else if typ == "deadline" {
			sev = domain.SeverityWarning
		}
		alerts = append(alerts, domain.ConflictAlert{
			Severity:       sev,
			Reason:         typ + " overlap: " + title,
			RelatedEventID: id,
		})
	}
	if alerts == nil {
		alerts = []domain.ConflictAlert{}
	}
	return alerts
}

func (s *PgStore) SearchTransport(_from, to string) []domain.TransportOption {
	dests := s.loadDestinations()
	for _, entry := range dests {
		if strings.EqualFold(entry.City, to) {
			return []domain.TransportOption{
				{Provider: "EuroRail Connect", Mode: "train", DurationHours: entry.BaseTravelHrs, Price: math.Round(entry.TransportBase * 0.9), Deeplink: "https://example.com/train"},
				{Provider: "SkySaver", Mode: "flight", DurationHours: roundF(entry.BaseTravelHrs*0.65, 1), Price: math.Round(entry.TransportBase * 1.18), Deeplink: "https://example.com/flight"},
			}
		}
	}
	return []domain.TransportOption{}
}

func (s *PgStore) SearchStays(city string) []domain.StayOption {
	dests := s.loadDestinations()
	for _, entry := range dests {
		if strings.EqualFold(entry.City, city) {
			return []domain.StayOption{
				{Provider: "HostelGraph", Kind: "hostel", NightlyPrice: entry.HostelNightEUR, Rating: 4.2, Deeplink: "https://example.com/hostel"},
				{Provider: "StudentStay", Kind: "budget-hotel", NightlyPrice: entry.HostelNightEUR + 14, Rating: 4.0, Deeplink: "https://example.com/hotel"},
			}
		}
	}
	return []domain.StayOption{}
}
