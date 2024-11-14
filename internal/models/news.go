package models

type FinanceNews struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	PublishedDate string `json:"published_date"`
	OriginURL     string `json:"origin_url"`
}
