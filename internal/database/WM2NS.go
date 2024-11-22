package database

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

func SaveWM2NS(wm2ns models.WM2NS, cfg *config.Config) error {
	query := NewSupabaseURLQuery(cfg, "economic_indicators")
	requestURL := query.Build()
	logging.Logger.WithField("url", requestURL).Debug("SAVE URL")

	economicIndicators := []models.EconomicIndicator{}

	for _, observation := range wm2ns.Observations {
		economicIndicators = append(economicIndicators, models.EconomicIndicator{
			Name:        "WM2NS",
			Country:     "US",
			ReleaseDate: observation.Date.Time,
			ActualValue: observation.Value,
			Unit:        "USD Billion", // Update the unit here
		})
	}

	wm2nsJsonData, err := json.Marshal(economicIndicators)
	if err != nil {
		return err
	}

	resp, err := POST(requestURL, cfg, wm2nsJsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != http.StatusCreated {
		logging.Logger.WithField("status", resp.Status).Error("Failed to save WM2NS")
		return errors.New("failed to save WM2NS")
	}

	logging.Logger.WithField("status", resp.Status).Debug("SAVE STATUS")

	return nil
}
