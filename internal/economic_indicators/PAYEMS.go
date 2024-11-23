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

const PAYEMS_NAME = "PAYEMS"
const PAYEMS_OBSERVATION_START_DATE = "1939-01-01"

func GatherPAYEMS(cfg *config.Config) error {
	logging.Logger.Debug("Getting PAYEMS")

	url := getFREDQuery(PAYEMS_NAME, PAYEMS_OBSERVATION_START_DATE, "m", "")
	logging.Logger.WithField("url", url).Debug("FRED URL")

	response, err := http.Get(url)
	if err != nil {
		logging.Logger.WithError(err).Error("Error getting PAYEMS")
		return err
	}
	defer response.Body.Close()

	var payems models.PAYEMS
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Logger.WithError(err).Error("Error reading PAYEMS response body")
		return err
	}
	err = json.Unmarshal(body, &payems)
	if err != nil {
		logging.Logger.WithError(err).Error("Error unmarshalling PAYEMS")
		return err
	}

	database.DeleteAllEconomicIndicators(cfg, PAYEMS_NAME)

	return database.SavePAYEMS(payems, cfg)
}
