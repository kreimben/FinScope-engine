package database

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

func SaveICSA(icsa models.ICSA, cfg *config.Config) error {
	query := NewSupabaseURLQuery(cfg, "economic_indicators")
	requestURL := query.Build()
	logging.Logger.WithField("url", requestURL).Debug("SAVE URL")

	economicIndicators := []models.EconomicIndicator{}

	for _, observation := range icsa.Observations {
		economicIndicators = append(economicIndicators, models.EconomicIndicator{
			Name:        "ICSA",
			Country:     "US",
			ReleaseDate: observation.Date.Time,
			ActualValue: observation.Value,
			Unit:        "Number of Claims", // Update the unit here
		})
	}

	icsaJsonData, err := json.Marshal(economicIndicators)
	if err != nil {
		return err
	}

	resp, err := POST(requestURL, cfg, icsaJsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != http.StatusCreated {
		logging.Logger.WithField("status", resp.Status).Error("Failed to save ICSA")
		return errors.New("failed to save ICSA")
	}

	logging.Logger.WithField("status", resp.Status).Debug("SAVE STATUS")

	return nil
}
