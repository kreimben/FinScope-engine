package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type ValuationIndicator struct {
	Name  string    `json:"name"`
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

// Add UnmarshalJSON for ValuationIndicator
func (vi *ValuationIndicator) UnmarshalJSON(b []byte) error {
	type Alias ValuationIndicator
	aux := &struct {
		Date string `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(vi),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	// Attempt to parse the date string in RFC3339 format
	parsedTime, err := time.Parse(time.RFC3339, aux.Date)
	if err != nil {
		// If RFC3339 fails, try "2006-01-02" format
		parsedTime, err = time.Parse("2006-01-02", aux.Date)
		if err != nil {
			return fmt.Errorf("error parsing date: %v", err)
		}
	}

	vi.Date = parsedTime
	return nil
}
