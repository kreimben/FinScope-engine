package models

import "time"

type ValuationIndicator struct {
	Name  string    `json:"name"`
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}
