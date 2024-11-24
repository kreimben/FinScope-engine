package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

type SupabaseURLQuery struct {
	url string
}

func NewSupabaseURLQuery(cfg *config.Config, table string) *SupabaseURLQuery {
	return &SupabaseURLQuery{
		url: cfg.SupabaseURL + "/rest/v1/" + table + "?",
	}
}

func (s *SupabaseURLQuery) Add(key, value string) *SupabaseURLQuery {
	s.url += key + "=" + value
	return s
}

func (s *SupabaseURLQuery) And() *SupabaseURLQuery {
	s.url += "&"
	return s
}

func (s *SupabaseURLQuery) Build() string {
	return s.url
}

func GET(requestURL string, cfg *config.Config) (*http.Response, error) {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", cfg.SupabaseAnonKey)
	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseAnonKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}

func POST(requestURL string, cfg *config.Config, jsonData []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", cfg.SupabaseAnonKey)
	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseAnonKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}

func PATCH(requestURL string, cfg *config.Config, jsonData []byte) (*http.Response, error) {
	req, err := http.NewRequest("PATCH", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", cfg.SupabaseAnonKey)
	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseAnonKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}

func DELETE(requestURL string, cfg *config.Config) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", cfg.SupabaseAnonKey)
	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseAnonKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}

func CheckURLExists(cfg *config.Config, urlStr string) (bool, error) {
	// URL encode the origin_url parameter
	encodedURL := url.QueryEscape(urlStr)

	query := NewSupabaseURLQuery(cfg, "finance_news")
	query.Add("origin_url", "eq."+encodedURL)
	requestURL := query.Build()

	logging.Logger.WithField("requestURL", requestURL).Debug("Checking URL in database")

	resp, err := GET(requestURL, cfg)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	if len(body) == 0 {
		return false, nil // No data found, URL does not exist
	}

	var result []models.FinanceNews
	if err := json.Unmarshal(body, &result); err != nil {
		logging.Logger.WithField("body", string(body)).Error("Failed to decode response")
		return false, fmt.Errorf("failed to decode response: %v", err)
	}

	return len(result) > 0, nil
}

func InsertNews(cfg *config.Config, data models.FinanceNews) error {
	jsonData, err := json.Marshal(map[string]interface{}{
		"title":          data.Title,
		"content":        data.Content,
		"published_date": data.PublishedDate.Format(time.RFC3339),
		"origin_url":     data.OriginURL,
	})
	if err != nil {
		return err
	}

	query := NewSupabaseURLQuery(cfg, "finance_news")
	requestURL := query.Build()

	resp, err := POST(requestURL, cfg, jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to insert news: %s", resp.Status)
	}

	return nil
}

func SaveReleaseDate(seriesID string, releaseDate time.Time, cfg *config.Config) error {
	data := map[string]interface{}{
		"series_id":    seriesID,
		"release_date": releaseDate.Format(time.RFC3339),
		"done":         false,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	query := NewSupabaseURLQuery(cfg, "economic_indicator_release_schedules")
	requestURL := query.Build()

	resp, err := POST(requestURL, cfg, jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		logging.Logger.WithField("data", data).Debug("Release date already exists")
		return nil
	} else if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to save release date: %s", resp.Status)
	}

	return nil
}

func GetNextReleaseDate(seriesID string, cfg *config.Config) (models.ReleaseDate, error) {
	query := NewSupabaseURLQuery(cfg, "economic_indicator_release_schedules")
	query.Add("done", "eq.false").And()
	query.Add("series_id", "eq."+seriesID).And()
	query.Add("release_date", "gte."+time.Now().UTC().Format("2006-01-02")).And()
	query.Add("order", "release_date.asc").And()
	query.Add("limit", "1")
	requestURL := query.Build()

	resp, err := GET(requestURL, cfg)
	if err != nil {
		return models.ReleaseDate{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ReleaseDate{}, err
	}

	var releaseDate []models.ReleaseDate
	if err := json.Unmarshal(body, &releaseDate); err != nil {
		logging.Logger.WithField("body", string(body)).Error("Failed to decode response")
		return models.ReleaseDate{}, err
	}

	return releaseDate[0], nil
}

func MarkReleaseDateAsDone(seriesID string, releaseDate time.Time, cfg *config.Config) error {
	updateData := map[string]interface{}{
		"done": true,
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return err
	}

	query := NewSupabaseURLQuery(cfg, "economic_indicator_release_schedule")
	query.Add("series_id", "eq."+seriesID).And()
	query.Add("release_date", "eq."+releaseDate.Format("2006-01-02"))
	requestURL := query.Build()

	resp, err := PATCH(requestURL, cfg, jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to mark release date as done: %s", resp.Status)
	}

	return nil
}

func GetLatestValueDate(seriesID string, cfg *config.Config) (time.Time, error) {
	query := NewSupabaseURLQuery(cfg, "economic_indicators")
	query.Add("name", "eq."+seriesID).And()
	query.Add("order", "release_date.desc").And()
	query.Add("limit", "1")
	requestURL := query.Build()

	resp, err := GET(requestURL, cfg)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, err
	}
	logging.Logger.WithField("body", string(body)).Debug("Latest value date")

	var latestData []models.EconomicIndicator
	if err := json.Unmarshal(body, &latestData); err != nil {
		logging.Logger.WithField("body", string(body)).Error("Failed to decode response")
		return time.Time{}, err
	}

	return latestData[0].ReleaseDate, nil
}
