package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/models"
)

func CheckURLExists(cfg *config.Config, url string) (bool, error) {
	// tokenString, err := auth.GenerateJWT(cfg)
	// if err != nil {
	// 	return false, err
	// }

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/v1/finance_news?origin_url=eq.%s", cfg.SupabaseURL, url), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("apikey", cfg.SupabaseAnonKey)
	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseAnonKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// // Read the response body for logging
	// bodyBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return false, fmt.Errorf("failed to read response body: %v", err)
	// }

	// // Log the response body
	// logging.Logger.WithField("response", string(bodyBytes)).Info("Response body")

	// Create a new reader with the body bytes for decoding
	var result []models.FinanceNews
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response: %v", err)
	}

	return len(result) > 0, nil
}

func InsertNews(cfg *config.Config, data models.FinanceNews) error {
	// tokenString, err := auth.GenerateJWT(cfg)
	// if err != nil {
	// 	return err
	// }

	jsonData, err := json.Marshal(map[string]interface{}{
		"title":          data.Title,
		"content":        data.Content,
		"published_date": data.PublishedDate.Format(time.RFC3339),
		"origin_url":     data.OriginURL,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/v1/finance_news", cfg.SupabaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("apikey", cfg.SupabaseAnonKey)
	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseAnonKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to insert news: %s", resp.Status)
	}

	return nil
}
