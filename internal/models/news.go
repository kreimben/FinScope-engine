package models

import (
	"time"
)

type FinanceNews struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	PublishedDate time.Time `json:"published_date"`
	OriginURL     string    `json:"origin_url"`
}
