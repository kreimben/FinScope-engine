package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type FinanceNews struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	PublishedDate time.Time `json:"published_date"`
	OriginURL     string    `json:"origin_url"`
	ContentVector []float32 `json:"content_vector"`
}

// Add UnmarshalJSON for FinanceNews
func (fn *FinanceNews) UnmarshalJSON(b []byte) error {
	type Alias FinanceNews
	aux := &struct {
		PublishedDate string `json:"published_date"`
		*Alias
	}{
		Alias: (*Alias)(fn),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	// Parse the published_date string
	parsedTime, err := time.Parse(time.RFC3339, aux.PublishedDate)
	if err != nil {
		return fmt.Errorf("error parsing published_date: %v", err)
	}

	fn.PublishedDate = parsedTime
	return nil
}
