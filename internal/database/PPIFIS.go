package database

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

func SavePPIFIS(ppifis models.PPIFIS, cfg *config.Config) error {
	query := NewSupabaseURLQuery(cfg, "economic_indicators")
	requestURL := query.Build()
	logging.Logger.WithField("url", requestURL).Debug("SAVE URL")

	economicIndicators := []models.EconomicIndicator{}

	for _, observation := range ppifis.Observations {
		economicIndicators = append(economicIndicators, models.EconomicIndicator{
			Name:        "PPIFIS",
			Country:     "US",
			ReleaseDate: observation.Date.Time,
			ActualValue: observation.Value,
			Unit:        "Percent", // Update the unit here
		})
	}

	ppifisJsonData, err := json.Marshal(economicIndicators)
	if err != nil {
		return err
	}

	resp, err := POST(requestURL, cfg, ppifisJsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != http.StatusCreated {
		logging.Logger.WithField("status", resp.Status).Error("Failed to save PPIFIS")
		return errors.New("failed to save PPIFIS")
	}

	logging.Logger.WithField("status", resp.Status).Debug("SAVE STATUS")

	return nil
}
