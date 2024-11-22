package database

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

func SavePAYEMS(payems models.PAYEMS, cfg *config.Config) error {
	query := NewSupabaseURLQuery(cfg, "economic_indicators")
	requestURL := query.Build()
	logging.Logger.WithField("url", requestURL).Debug("SAVE URL")

	economicIndicators := []models.EconomicIndicator{}

	for _, observation := range payems.Observations {
		economicIndicators = append(economicIndicators, models.EconomicIndicator{
			Name:        "PAYEMS",
			Country:     "US",
			ReleaseDate: observation.Date.Time,
			ActualValue: observation.Value,
			Unit:        "Thousands of Persons", // Update the unit here
		})
	}

	payemsJsonData, err := json.Marshal(economicIndicators)
	if err != nil {
		return err
	}

	resp, err := POST(requestURL, cfg, payemsJsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != http.StatusCreated {
		logging.Logger.WithField("status", resp.Status).Error("Failed to save PAYEMS")
		return errors.New("failed to save PAYEMS")
	}

	logging.Logger.WithField("status", resp.Status).Debug("SAVE STATUS")

	return nil
}
