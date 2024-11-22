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

const WM2NS_NAME = "WM2NS"
const WM2NS_OBSERVATION_START_DATE = "1980-12-01"

func GatherWM2NS(cfg *config.Config) error {
	logging.Logger.Debug("Getting WM2NS")

	url := getURLQuery(WM2NS_NAME, WM2NS_OBSERVATION_START_DATE, "m", "")
	logging.Logger.WithField("url", url).Debug("FRED URL")

	response, err := http.Get(url)
	if err != nil {
		logging.Logger.WithError(err).Error("Error getting WM2NS")
		return err
	}
	defer response.Body.Close()

	var wm2ns models.WM2NS
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Logger.WithError(err).Error("Error reading WM2NS response body")
		return err
	}
	err = json.Unmarshal(body, &wm2ns)
	if err != nil {
		logging.Logger.WithError(err).Error("Error unmarshalling WM2NS")
		return err
	}

	database.DeleteAllEconomicIndicators(cfg, WM2NS_NAME)

	return database.SaveWM2NS(wm2ns, cfg)
}
