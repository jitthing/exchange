package domain

type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityHighRisk Severity = "high-risk"
)

type AcademicEventType string

const (
	AcademicClass    AcademicEventType = "class"
	AcademicDeadline AcademicEventType = "deadline"
	AcademicExam     AcademicEventType = "exam"
	AcademicHoliday  AcademicEventType = "holiday"
)

type TravelWindow struct {
	ID        string   `json:"id"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	Score     int      `json:"score"`
	Conflicts []string `json:"conflicts"`
}

type AcademicEvent struct {
	ID       string            `json:"id"`
	Type     AcademicEventType `json:"type"`
	Title    string            `json:"title"`
	Start    string            `json:"start"`
	End      string            `json:"end"`
	Priority int               `json:"priority"`
}

type TripConstraint struct {
	BudgetCap      float64 `json:"budgetCap"`
	MaxTravelHours float64 `json:"maxTravelHours"`
	PartySize      int     `json:"partySize"`
	Style          string  `json:"style"`
	WindowID       string  `json:"windowId"`
	DepartureCity  string  `json:"departureCity"`
}

type TransportOption struct {
	Provider      string  `json:"provider"`
	Mode          string  `json:"mode"`
	DurationHours float64 `json:"durationHours"`
	Price         float64 `json:"price"`
	Deeplink      string  `json:"deeplink"`
}

type StayOption struct {
	Provider     string  `json:"provider"`
	Kind         string  `json:"kind"`
	NightlyPrice float64 `json:"nightlyPrice"`
	Rating       float64 `json:"rating"`
	Deeplink     string  `json:"deeplink"`
}

type ConflictAlert struct {
	Severity       Severity `json:"severity"`
	Reason         string   `json:"reason"`
	RelatedEventID string   `json:"relatedEventId"`
}

type TripOption struct {
	ID                 string            `json:"id"`
	Destination        string            `json:"destination"`
	ReasonTags         []string          `json:"reasonTags"`
	TotalEstimatedCost float64           `json:"totalEstimatedCost"`
	TransportOptions   []TransportOption `json:"transportOptions"`
	StayOptions        []StayOption      `json:"stayOptions"`
	RiskLevel          Severity          `json:"riskLevel"`
}

type Trip struct {
	ID            string   `json:"id"`
	OwnerID       string   `json:"ownerId"`
	Destination   string   `json:"destination"`
	WindowID      string   `json:"windowId"`
	Members       []string `json:"members"`
	Itinerary     []string `json:"itinerary"`
	EstimatedCost float64  `json:"estimatedCost"`
}

type BudgetEntry struct {
	ID       string  `json:"id"`
	UserID   string  `json:"userId"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Date     string  `json:"date"`
	TripID   string  `json:"tripId,omitempty"`
	Note     string  `json:"note,omitempty"`
}

type ForecastResult struct {
	ProjectedMonthlySpend float64 `json:"projectedMonthlySpend"`
	RemainingBudget       float64 `json:"remainingBudget"`
	Affordability         string  `json:"affordability"`
}
