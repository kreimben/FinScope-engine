package models

import "time"

type EconomicIndicator struct {
	Name          string    `json:"name"`    // GDP, CPI, etc.
	Country       string    `json:"country"` // US, EU, etc.
	ReleaseDate   time.Time `json:"release_date"`
	ActualValue   float64   `json:"actual_value"`
	ForecastValue float64   `json:"forecast_value,omitempty"`
	PreviousValue float64   `json:"previous_value,omitempty"`
	Unit          string    `json:"unit"` // USD, EUR etc...
}
