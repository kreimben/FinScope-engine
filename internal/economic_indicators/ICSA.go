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

const ICSA_NAME = "ICSA"
const ICSA_OBSERVATION_START_DATE = "1967-01-01"

func GatherICSA(cfg *config.Config) error {
	logging.Logger.Debug("Getting ICSA")

	url := getURLQuery(ICSA_NAME, ICSA_OBSERVATION_START_DATE, "w", "")
	logging.Logger.WithField("url", url).Debug("FRED URL")

	response, err := http.Get(url)
	if err != nil {
		logging.Logger.WithError(err).Error("Error getting ICSA")
		return err
	}
	defer response.Body.Close()

	var icsa models.ICSA
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Logger.WithError(err).Error("Error reading ICSA response body")
		return err
	}
	err = json.Unmarshal(body, &icsa)
	if err != nil {
		logging.Logger.WithError(err).Error("Error unmarshalling ICSA")
		return err
	}

	database.DeleteAllEconomicIndicators(cfg, ICSA_NAME)

	return database.SaveICSA(icsa, cfg)
}
