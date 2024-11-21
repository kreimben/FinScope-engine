package database

import (
	"errors"
	"net/http"

	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

// delete all the data in the table
func DeleteAllEconomicIndicators(cfg *config.Config, name string) error {
	query := NewSupabaseURLQuery(cfg, "economic_indicators")
	query.And("name", "eq."+name)
	requestURL := query.Build()
	logging.Logger.WithField("url", requestURL).Debug("DELETE URL")

	resp, err := DELETE(requestURL, cfg)
	if err != nil {
		logging.Logger.WithError(err).Error("Error occurred while sending DELETE request")
		return err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logging.Logger.WithError(cerr).Error("Error occurred while closing response body")
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		logging.Logger.WithField("status", resp.Status).Error("Failed to delete GDP data")
		return errors.New("failed to delete GDP data")
	}

	return nil
}
