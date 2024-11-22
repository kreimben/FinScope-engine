package economic_indicators

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/database"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

const GDP_NAME = "GDP"
const GDP_OBSERVATION_START_DATE = "1948-01-01"

// GatherGDP gets the GDP data from FRED and save it to the database
func GatherGDP(cfg *config.Config) error {
	logging.Logger.Debug("Getting GDP")

	url := getURLQuery(GDP_NAME, GDP_OBSERVATION_START_DATE, "q", "")
	logging.Logger.WithField("url", url).Debug("FRED URL")

	// get the response
	response, err := http.Get(url)
	if err != nil {
		logging.Logger.WithError(err).Error("Error getting GDP")
		return err
	}
	defer response.Body.Close()

	// unmarshal the response
	var gdp models.GDP
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Logger.WithError(err).Error("Error reading GDP response body")
		return err
	}
	err = json.Unmarshal(body, &gdp)
	if err != nil {
		logging.Logger.WithError(err).Error("Error unmarshalling GDP")
		return err
	}

	// delete all the data in the table
	database.DeleteAllEconomicIndicators(cfg, GDP_NAME)

	// save the data to the database
	return database.SaveGDP(gdp, cfg)
}
