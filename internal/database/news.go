package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/kreimben/FinScope-engine/internal/config"
)

type DB struct {
	cfg *config.Config
}

type FinanceNews struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	PublishedDate time.Time `json:"published_date"`
	OriginURL     string    `json:"origin_url"`
}

func New(cfg *config.Config) *DB {
	return &DB{cfg: cfg}
}

func (db *DB) GetNewsByQuery(ctx context.Context, query string) ([]FinanceNews, error) {
	// Create Supabase URL query builder
	queryBuilder := NewSupabaseURLQuery(db.cfg, "finance_news")
	queryBuilder.Add("select", "title,content,published_date,origin_url")
	queryBuilder.Add("textSearch", fmt.Sprintf("title,content.%s", query))
	queryBuilder.Add("order", "published_date.desc")
	queryBuilder.Add("limit", "50")
	queryURL := queryBuilder.Build()

	// Make GET request to Supabase
	resp, err := GET(queryURL, db.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get news: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse JSON response
	var news []FinanceNews
	if err := json.Unmarshal(body, &news); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return news, nil
}
