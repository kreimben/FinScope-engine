package economic_indicators

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"encoding/json"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/database"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

const FRED_OBSERVATIONS_BASE_URL = "https://api.stlouisfed.org/fred/series/observations"

var releaseIDs = map[string]int{
	"GDP":      53,
	"CPIAUCSL": 10,
	"UNRATE":   50,
	"WM2NS":    21,
	"DFEDTARU": 101,
	"PCEPI":    54,
	"PAYEMS":   50,
	"PPIFIS":   46,
	"ICSA":     180,
}
var releaseIDFunctions = map[string]func(cfg *config.Config) error{
	"GDP":      GatherGDP,
	"CPIAUCSL": GatherCPI,
	"UNRATE":   GatherUNRATE,
	"WM2NS":    GatherWM2NS,
	"DFEDTARU": GatherDFEDTARU,
	"PCEPI":    GatherPCEPI,
	"PAYEMS":   GatherPAYEMS,
	"PPIFIS":   GatherPPIFIS,
	"ICSA":     GatherICSA,
}

func FRED_API_KEY() string {
	return os.Getenv("FRED_API_KEY")
}

func getFREDQuery(seriesID string, observationStartDate string, frequency string, units string) string {
	query := NewFSQuery(FRED_OBSERVATIONS_BASE_URL)
	query.Add("api_key", FRED_API_KEY()).And()
	query.Add("series_id", seriesID).And()
	query.Add("frequency", frequency).And()
	query.Add("observation_start", observationStartDate)

	if units != "" {
		query.And().Add("units", units)
	}
	return query.Build()
}

func GatherReleaseSchedules(cfg *config.Config) error {
	logging.Logger.Debug("Getting Release Schedules")

	for seriesID, releaseID := range releaseIDs {
		query := NewFSQuery("https://api.stlouisfed.org/fred/release/dates")
		query.Add("api_key", FRED_API_KEY()).And()
		query.Add("include_release_dates_with_no_data", "true").And()
		query.Add("realtime_start", time.Now().Format("2006-01-02")).And()
		query.Add("file_type", "json").And()
		query.Add("release_id", fmt.Sprintf("%d", releaseID))
		url := query.Build()
		logging.Logger.WithField("url", url).Debug("FRED Release Schedule URL")

		response, err := http.Get(url)
		if err != nil {
			logging.Logger.WithError(err).Error("Error getting release schedules")
			return err
		}
		defer response.Body.Close()

		var releaseResponse struct {
			ReleaseDates []models.ReleaseDate `json:"release_dates"`
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			logging.Logger.WithError(err).Error("Error reading release schedules response body")
			return err
		}
		err = json.Unmarshal(body, &releaseResponse)
		if err != nil {
			logging.Logger.WithError(err).Error("Error unmarshalling release schedules")
			return err
		}

		for _, rd := range releaseResponse.ReleaseDates {
			err = database.SaveReleaseDate(seriesID, rd.ReleaseDate, cfg)
			if err != nil {
				logging.Logger.WithError(err).Error("Error saving release date to database")
				continue
			}
		}
	}

	return nil
}

func GatherTodayReleaseIndicatorAndMarkAsDone(cfg *config.Config) error {
	today := time.Now().UTC()

	for seriesID := range releaseIDs {
		// Get the next release date
		releaseDate, err := database.GetNextReleaseDate(seriesID, cfg)
		if err != nil {
			logging.Logger.WithError(err).Error("Error getting next release date")
			continue
		}

		if releaseDate.ReleaseDate.Format("2006-01-02") == today.Format("2006-01-02") {
			// Gather the indicator
			err = releaseIDFunctions[seriesID](cfg)
			if err != nil {
				logging.Logger.WithError(err).Error("Error gathering indicator")
				continue
			}

			// check latest data whether today is the day.
			latestData, err := database.GetLatestValueDate(seriesID, cfg)
			if err != nil {
				logging.Logger.WithError(err).Error("Error getting latest data")
				continue
			}

			if latestData.Format("2006-01-02") == today.Format("2006-01-02") {
				// Mark the release date as done
				database.MarkReleaseDateAsDone(seriesID, today, cfg)
			}
		}
	}

	return nil
}
