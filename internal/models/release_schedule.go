// Start of Selection
package models

import (
	"encoding/json"
	"time"
)

/*
 * "[{\"series_id\":\"PPIFIS\",\"release_date\":\"2024-12-12T00:00:00+00:00\",\"done\":false}]"
 * "{..., \"release_dates\":[{\"release_id\":10,\"date\":\"2024-12-11\"}]}"
 */
type ReleaseDate struct {
	SeriesID    string    `json:"series_id,omitempty"`    // GDP
	ReleaseId   uint16    `json:"release_id,omitempty"`   // 53
	ReleaseDate time.Time `json:"release_date,omitempty"` // from Supabase
	Date        time.Time `json:"date,omitempty"`         // from Fred
}

// Add UnmarshalJSON for ReleaseDate
func (rd *ReleaseDate) UnmarshalJSON(b []byte) error {
	type Alias ReleaseDate
	aux := &struct {
		ReleaseDate string `json:"release_date,omitempty"`
		Date        string `json:"date,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(rd),
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	if aux.ReleaseDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, aux.ReleaseDate)
		if err != nil {
			return err
		}
		rd.ReleaseDate = parsedTime
	} else if aux.Date != "" {
		parsedTime, err := time.Parse("2006-01-02", aux.Date)
		if err != nil {
			return err
		}
		rd.Date = parsedTime
	}
	return nil
}
