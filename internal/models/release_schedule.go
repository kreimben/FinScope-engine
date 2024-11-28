package models

import (
	"fmt"
	"strings"
	"time"
)

type ReleaseDate struct {
	SeriesID    string    `json:"series_id"`
	ReleaseDate time.Time `json:"release_date"`
}

// Add UnmarshalJSON for ReleaseDate
func (rd *ReleaseDate) UnmarshalJSON(b []byte) error {
	// Remove quotes from the JSON string
	str := strings.Trim(string(b), `"`)

	// Parse the date string
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return fmt.Errorf("error parsing release_date: %v", err)
	}

	rd.ReleaseDate = t
	return nil
}
