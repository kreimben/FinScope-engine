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

const UNRATE_NAME = "UNRATE"
const UNRATE_OBSERVATION_START_DATE = "1948-01-01"

func GatherUNRATE(cfg *config.Config) {
	logging.Logger.Debug("Getting UNRATE")

	url := getURLQuery(UNRATE_NAME, UNRATE_OBSERVATION_START_DATE, "m")
	logging.Logger.WithField("url", url).Debug("FRED URL")

	response, err := http.Get(url)
	if err != nil {
		logging.Logger.WithError(err).Error("Error getting UNRATE")
		return
	}
	defer response.Body.Close()

	var unrate models.UNRATE
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Logger.WithError(err).Error("Error reading UNRATE response body")
		return
	}
	err = json.Unmarshal(body, &unrate)
	if err != nil {
		logging.Logger.WithError(err).Error("Error unmarshalling UNRATE")
		return
	}

	database.DeleteAllEconomicIndicators(cfg, UNRATE_NAME)

	database.SaveUNRATE(unrate, cfg)
}
