package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONStringSlice is a custom type for JSONB text arrays.
type JSONStringSlice []string

func (j *JSONStringSlice) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONStringSlice) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

// --- GORM Models ---

type AcademicEventModel struct {
	ID        string `gorm:"column:id;primaryKey"`
	Type      string `gorm:"column:type"`
	Title     string `gorm:"column:title"`
	StartDate string `gorm:"column:start_date"`
	EndDate   string `gorm:"column:end_date"`
	Priority  int    `gorm:"column:priority"`
}

func (AcademicEventModel) TableName() string { return "academic_events" }

type TravelWindowModel struct {
	ID        string          `gorm:"column:id;primaryKey"`
	StartDate string          `gorm:"column:start_date"`
	EndDate   string          `gorm:"column:end_date"`
	Score     int             `gorm:"column:score"`
	Conflicts JSONStringSlice `gorm:"column:conflicts;type:jsonb"`
}

func (TravelWindowModel) TableName() string { return "travel_windows" }

type TripModel struct {
	ID            string          `gorm:"column:id;primaryKey"`
	OwnerID       string          `gorm:"column:owner_id"`
	Destination   string          `gorm:"column:destination"`
	WindowID      string          `gorm:"column:window_id"`
	Members       JSONStringSlice `gorm:"column:members;type:jsonb"`
	Itinerary     JSONStringSlice `gorm:"column:itinerary;type:jsonb"`
	EstimatedCost float64         `gorm:"column:estimated_cost"`
}

func (TripModel) TableName() string { return "trips" }

type BudgetEntryModel struct {
	ID       string  `gorm:"column:id;primaryKey"`
	UserID   string  `gorm:"column:user_id"`
	Category string  `gorm:"column:category"`
	Amount   float64 `gorm:"column:amount"`
	Currency string  `gorm:"column:currency"`
	Date     string  `gorm:"column:date"`
	TripID   string  `gorm:"column:trip_id"`
	Note     string  `gorm:"column:note"`
}

func (BudgetEntryModel) TableName() string { return "budget_entries" }

type MonthlyBudgetModel struct {
	UserID string  `gorm:"column:user_id;primaryKey"`
	Budget float64 `gorm:"column:budget"`
}

func (MonthlyBudgetModel) TableName() string { return "monthly_budgets" }

type DestinationModel struct {
	City           string          `gorm:"column:city;primaryKey"`
	BaseTravelHrs  float64         `gorm:"column:base_travel_hrs"`
	TransportBase  float64         `gorm:"column:transport_base"`
	HostelNightEUR float64         `gorm:"column:hostel_night_eur"`
	Tags           JSONStringSlice `gorm:"column:tags;type:jsonb"`
}

func (DestinationModel) TableName() string { return "destinations" }
