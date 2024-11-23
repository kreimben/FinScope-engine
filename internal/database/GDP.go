package database

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

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
			Unit:        "USD Billion",
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

	// Optionally read and log the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logging.Logger.WithField("response_body", string(body)).Debug("SAVE GDP RESPONSE BODY")

	// check status code
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		logging.Logger.WithField("status", resp.Status).Error("Failed to save GDP")
		return errors.New("failed to save GDP")
	}

	logging.Logger.WithField("status", resp.Status).Debug("SAVE STATUS")

	return nil
}

// GetUSGDPByDate retrieves the GDP data for the US by a specific date
func GetUSGDPByDate(cfg *config.Config, date time.Time) (models.GDP, error) {
	query := NewSupabaseURLQuery(cfg, "economic_indicators")
	query.Add("name", "eq.GDP").And()
	query.Add("country", "eq.US").And()
	query.Add("select", "name,country,release_date,actual_value,unit").And()
	query.Add("release_date", "lte."+date.Format("2006-01-02T15:04:05Z")).And()
	query.Add("order", "release_date.desc").And()
	query.Add("limit", "1")
	requestURL := query.Build()
	logging.Logger.WithField("url", requestURL).Debug("GET GDP URL")

	resp, err := GET(requestURL, cfg)
	if err != nil {
		return models.GDP{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logging.Logger.WithField("status", resp.Status).Error("Failed to get GDP")
		return models.GDP{}, errors.New("failed to get GDP")
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.GDP{}, err
	}
	logging.Logger.WithField("body", string(body)).Debug("GET GDP RESPONSE BODY")

	var economicIndicators []models.EconomicIndicator
	// Unmarshal the read body
	if err := json.Unmarshal(body, &economicIndicators); err != nil {
		return models.GDP{}, err
	}

	gdp := models.GDP{
		Observations: []models.Observation{},
	}

	if len(economicIndicators) > 0 {
		gdp.Observations = append(gdp.Observations, models.Observation{
			Date:  models.CustomDate{Time: economicIndicators[0].ReleaseDate},
			Value: economicIndicators[0].ActualValue,
		})
	} else {
		logging.Logger.Error("No GDP data found for the given date")
		return models.GDP{}, errors.New("no GDP data found")
	}

	logging.Logger.WithField("gdp", gdp).Debug("GET GDP")
	return gdp, nil
}
