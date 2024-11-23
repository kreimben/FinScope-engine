package database

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

func SaveBuffettIndicators(indicators []models.ValuationIndicator, cfg *config.Config) error {
	query := NewSupabaseURLQuery(cfg, "valuation_indicators")
	requestURL := query.Build()
	logging.Logger.WithField("url", requestURL).Debug("SAVE BUFFETT INDICATORS URL")

	jsonData, err := json.Marshal(indicators)
	if err != nil {
		return err
	}
	logging.Logger.WithField("data", string(jsonData)).Debug("SAVE BUFFETT INDICATORS DATA")

	resp, err := POST(requestURL, cfg, jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		logging.Logger.WithField("status", resp.Status).Error("Failed to save Buffett Indicators")
		return errors.New("failed to save Buffett Indicators")
	}

	logging.Logger.WithField("status", resp.Status).Debug("SAVE BUFFETT INDICATORS STATUS")
	return nil
}

func DeleteAllBuffettIndicators(cfg *config.Config) error {
	query := NewSupabaseURLQuery(cfg, "valuation_indicators")
	query.Add("name", "eq.Buffett_Indicator")
	requestURL := query.Build()
	logging.Logger.WithField("url", requestURL).Debug("DELETE BUFFETT INDICATORS URL")

	resp, err := DELETE(requestURL, cfg)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		logging.Logger.WithField("status", resp.Status).Error("Failed to delete Buffett Indicators")
		return errors.New("failed to delete Buffett Indicators")
	}

	return nil
}
