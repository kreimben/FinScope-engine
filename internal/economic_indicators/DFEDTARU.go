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

const DFEDTARU_NAME = "DFEDTARU"
const DFEDTARU_OBSERVATION_START_DATE = "2009-01-01"

func GatherDFEDTARU(cfg *config.Config) error {
	logging.Logger.Debug("Getting DFEDTARU")

	url := getURLQuery(DFEDTARU_NAME, DFEDTARU_OBSERVATION_START_DATE, "m", "")
	logging.Logger.WithField("url", url).Debug("FRED URL")

	response, err := http.Get(url)
	if err != nil {
		logging.Logger.WithError(err).Error("Error getting DFEDTARU")
		return err
	}
	defer response.Body.Close()

	var dfedtaru models.DFEDTARU
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Logger.WithError(err).Error("Error reading DFEDTARU response body")
		return err
	}
	err = json.Unmarshal(body, &dfedtaru)
	if err != nil {
		logging.Logger.WithError(err).Error("Error unmarshalling DFEDTARU")
		return err
	}

	database.DeleteAllEconomicIndicators(cfg, DFEDTARU_NAME)

	return database.SaveDFEDTARU(dfedtaru, cfg)
}
