package main

import (
	"context"
	"encoding/json"
	"os"
	"slices"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/crawler"
	"github.com/kreimben/FinScope-engine/internal/economic_indicators"
	"github.com/kreimben/FinScope-engine/internal/valuation_indicators"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

var cfg *config.Config

func init() {
	// Load config.
	cfg = config.LoadConfig()
	logging.Logger = logging.NewLogger()
}

type FinScopeEngineEvent struct {
	Execute []string `json:"execute"`
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	// unmarshal the event
	var finScopeEngineEvent FinScopeEngineEvent
	err := json.Unmarshal(event, &finScopeEngineEvent)
	if err != nil {
		return err
	}

	// Economic Indicators
	if slices.Contains(finScopeEngineEvent.Execute, "GDP") {
		err := economic_indicators.GatherGDP(cfg)
		if err != nil {
			return err
		}
	}

	if slices.Contains(finScopeEngineEvent.Execute, "CPI") {
		err := economic_indicators.GatherCPI(cfg)
		if err != nil {
			return err
		}
	}

	if slices.Contains(finScopeEngineEvent.Execute, "UNRATE") {
		err := economic_indicators.GatherUNRATE(cfg)
		if err != nil {
			return err
		}
	}

	if slices.Contains(finScopeEngineEvent.Execute, "PCEPI") {
		err := economic_indicators.GatherPCEPI(cfg)
		if err != nil {
			return err
		}
	}

	if slices.Contains(finScopeEngineEvent.Execute, "DFEDTARU") {
		err := economic_indicators.GatherDFEDTARU(cfg)
		if err != nil {
			return err
		}
	}

	if slices.Contains(finScopeEngineEvent.Execute, "WM2NS") {
		err := economic_indicators.GatherWM2NS(cfg)
		if err != nil {
			return err
		}
	}

	if slices.Contains(finScopeEngineEvent.Execute, "PAYEMS") {
		err := economic_indicators.GatherPAYEMS(cfg)
		if err != nil {
			return err
		}
	}

	if slices.Contains(finScopeEngineEvent.Execute, "PPIFIS") {
		err := economic_indicators.GatherPPIFIS(cfg)
		if err != nil {
			return err
		}
	}

	if slices.Contains(finScopeEngineEvent.Execute, "ICSA") {
		err := economic_indicators.GatherICSA(cfg)
		if err != nil {
			return err
		}
	}

	// FRED Release Schedules
	if slices.Contains(finScopeEngineEvent.Execute, "FRED_Gather_Release_Schedules") {
		err := economic_indicators.GatherReleaseSchedules(cfg)
		if err != nil {
			return err
		}
	}

	if slices.Contains(finScopeEngineEvent.Execute, "FRED_Gather_Today_Release_Indicator_And_Mark_As_Done") {
		err := economic_indicators.GatherTodayReleaseIndicatorAndMarkAsDone(cfg)
		if err != nil {
			return err
		}
	}

	// Valuation Indicators
	if slices.Contains(finScopeEngineEvent.Execute, "Buffett_Indicator") {
		err := valuation_indicators.CalculateAndSaveHistoricalBuffettIndicator(
			cfg,
			time.Now().Add(-1*time.Hour*24),
			time.Now(),
		)
		if err != nil {
			return err
		}
	}

	// News
	if slices.Contains(finScopeEngineEvent.Execute, "yahoo_finance") {
		crawler.StartFinanceYahooCrawler(cfg)
	}

	if slices.Contains(finScopeEngineEvent.Execute, "benzinga") {
		crawler.StartBenzingaCrawler(cfg)
	}

	return nil
}

func main() {
	if os.Getenv("DEBUG") == "true" {
		handleRequest(context.Background(), json.RawMessage(`{"execute": ["DFEDTARU"]}`)) //  FRED_Gather_Today_Release_Indicator_And_Mark_As_Done
	} else {
		lambda.Start(handleRequest)
	}
}
