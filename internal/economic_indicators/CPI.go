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

const CPI_NAME = "CPIAUCSL"
const CPI_OBSERVATION_START_DATE = "1947-01-01"

func GatherCPI(cfg *config.Config) error {
	logging.Logger.Debug("Getting CPI")

	url := getFREDQuery(CPI_NAME, CPI_OBSERVATION_START_DATE, "m", "")
	logging.Logger.WithField("url", url).Debug("FRED URL")

	response, err := http.Get(url)
	if err != nil {
		logging.Logger.WithError(err).Error("Error getting CPI")
		return err
	}
	defer response.Body.Close()

	var cpi models.CPI
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Logger.WithError(err).Error("Error reading CPI response body")
		return err
	}
	err = json.Unmarshal(body, &cpi)
	if err != nil {
		logging.Logger.WithError(err).Error("Error unmarshalling CPI")
		return err
	}

	database.DeleteAllEconomicIndicators(cfg, CPI_NAME)

	return database.SaveCPI(cpi, cfg)
}
