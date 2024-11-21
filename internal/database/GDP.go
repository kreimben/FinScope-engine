package database

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

// save the data to the database
func SaveGDP(gdp models.GDP, cfg *config.Config) error {
	query := NewSupabaseURLQuery(cfg, "economic_indicators")
	requestURL := query.Build()
	logging.Logger.WithField("url", requestURL).Debug("SAVE URL")
	// GDP to EconomicIndicator
	economicIndicators := []models.EconomicIndicator{}

	for _, observation := range gdp.Observations {
		economicIndicators = append(economicIndicators, models.EconomicIndicator{
			Name:        "GDP",
			Country:     "US",
			ReleaseDate: observation.Date.Time,
			ActualValue: observation.Value,
			Unit:        "USD",
		})
	}

	gdpJsonData, err := json.Marshal(economicIndicators)
	if err != nil {
		return err
	}

	resp, err := POST(requestURL, cfg, gdpJsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	logging.Logger.WithField("status", resp.Status).Debug("SAVE STATUS")

	return nil
}

// delete all the data in the table
func DeleteAllEconomicIndicators(cfg *config.Config, name string) error {
	query := NewSupabaseURLQuery(cfg, "economic_indicators")
	query.And("name", "eq."+name)
	requestURL := query.Build()
	logging.Logger.WithField("url", requestURL).Debug("DELETE URL")

	resp, err := DELETE(requestURL, cfg)
	if err != nil {
		logging.Logger.WithError(err).Error("Error occurred while sending DELETE request")
		return err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logging.Logger.WithError(cerr).Error("Error occurred while closing response body")
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		logging.Logger.WithField("status", resp.Status).Error("Failed to delete GDP data")
		return errors.New("failed to delete GDP data")
	}

	return nil
}
