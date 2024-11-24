package models

import "time"

type ReleaseDate struct {
	SeriesID    string    `json:"series_id"`
	ReleaseDate time.Time `json:"release_date"`
}
