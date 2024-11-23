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

const PPIFIS_NAME = "PPIFIS"
const PPIFIS_OBSERVATION_START_DATE = "2009-12-01"

func GatherPPIFIS(cfg *config.Config) error {
	logging.Logger.Debug("Getting PPIFIS")

	url := getFREDQuery(PPIFIS_NAME, PPIFIS_OBSERVATION_START_DATE, "m", "pch")
	logging.Logger.WithField("url", url).Debug("FRED URL")

	response, err := http.Get(url)
	if err != nil {
		logging.Logger.WithError(err).Error("Error getting PPIFIS")
		return err
	}
	defer response.Body.Close()

	var ppifis models.PPIFIS
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Logger.WithError(err).Error("Error reading PPIFIS response body")
		return err
	}
	err = json.Unmarshal(body, &ppifis)
	if err != nil {
		logging.Logger.WithError(err).Error("Error unmarshalling PPIFIS")
		return err
	}

	database.DeleteAllEconomicIndicators(cfg, PPIFIS_NAME)

	return database.SavePPIFIS(ppifis, cfg)
}
