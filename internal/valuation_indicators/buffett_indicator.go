package valuation_indicators

import (
	"time"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/database"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/sirupsen/logrus"
)

const FTW5000_SYMBOL = "^FTW5000"
const BUFFETT_INDICATOR_NAME = "Buffett_Indicator"

func CalculateBuffettIndicator(cfg *config.Config, date time.Time) (float64, error) {
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		logging.Logger.WithFields(logrus.Fields{
			"date": date,
		}).Warn("Skipping weekend date for Buffett Indicator calculation")
		return 0, nil
	}

	// Get the US GDP value for the specified date
	gdp, err := database.GetUSGDPByDate(cfg, date)
	if err != nil {
		logging.Logger.WithError(err).Error("Error getting US GDP")
		return 0, err
	}

	p := &chart.Params{
		Symbol:   FTW5000_SYMBOL,
		Start:    &datetime.Datetime{Month: int(date.Month()), Day: date.Day(), Year: date.Year()},
		End:      &datetime.Datetime{Month: int(date.Month()), Day: date.Day() + 10, Year: date.Year()},
		Interval: datetime.OneDay,
	}
	iter := chart.Get(p)

	// Iterate over results. Will exit upon any error.
	var totalUSStockMarketValue float64
	for iter.Next() {
		b := iter.Bar()
		logging.Logger.WithFields(logrus.Fields{
			"time":     time.Unix(int64(b.Timestamp), 0), // 미국 장 시작시간
			"adjClose": b.AdjClose.Round(3),
		}).Debug("BAR")

		value, _ := b.AdjClose.Round(3).Float64()
		totalUSStockMarketValue = value
		break
	}

	if err := iter.Err(); err != nil {
		logging.Logger.WithError(err).Error("Error iterating over chart results")
		return 0, err
	}

	// Calculate the Buffett Indicator
	buffettIndicator := totalUSStockMarketValue / gdp.Observations[0].Value

	return buffettIndicator, nil
}

func CalculateAndSaveHistoricalBuffettIndicator(cfg *config.Config, startDate, endDate time.Time) error {
	logging.Logger.Info("[Buffett_Indicator] Calculating and saving historical Buffett Indicator")
	// err := database.DeleteAllBuffettIndicators(cfg)
	// if err != nil {
	// 	logging.Logger.WithError(err).Error("Error deleting all Buffett Indicators")
	// 	return err
	// }

	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			continue
		}

		buffettIndicator, err := CalculateBuffettIndicator(cfg, date)
		if err != nil {
			logging.Logger.WithError(err).Error("Error calculating Buffett Indicator")
			return err
		}
		logging.Logger.WithFields(logrus.Fields{
			"date":  date,
			"value": buffettIndicator,
		}).Info(BUFFETT_INDICATOR_NAME)

		err = database.SaveBuffettIndicators([]models.ValuationIndicator{
			{
				Name:  BUFFETT_INDICATOR_NAME,
				Date:  date,
				Value: buffettIndicator,
			},
		}, cfg)
		if err != nil {
			logging.Logger.WithError(err).Error("Error saving Buffett Indicator")
			return err
		}

		time.Sleep(300 * time.Millisecond)
	}

	logging.Logger.Info("[Buffett_Indicator] Completed calculating and saving historical Buffett Indicator")
	return nil
}
