package domain

// DataStore defines the interface for all data operations.
// Both the in-memory store and the PostgreSQL-backed store implement this.
type DataStore interface {
	ImportAcademicEvents(events []AcademicEvent) []AcademicEvent
	ListTravelWindows(from, to string) []TravelWindow
	OptimizeTrips(c TripConstraint) []TripOption
	GetTrip(id string) *Trip
	ShareTrip(tripID string, memberIDs []string) *Trip
	AddBudgetEntry(entry BudgetEntry) BudgetEntry
	ListBudgetEntries(userID string) []BudgetEntry
	Forecast(userID, tripID string) ForecastResult
	EvaluateConflicts(windowID string) []ConflictAlert
	SearchTransport(from, to string) []TransportOption
	SearchStays(city string) []StayOption
	Close() error
}
