package main

import (
	"context"
	"encoding/json"
	"os"
	"slices"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/crawler"
	"github.com/kreimben/FinScope-engine/internal/economic_indicators"
	"github.com/kreimben/FinScope-engine/pkg/logging"
)

var cfg *config.Config

func init() {
	// Load config.
	cfg = config.LoadConfig()
	logging.Logger = logging.NewLogger()
}

type FinScopeEngineEvent struct {
	Execute []string `json:"execute"` //
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	// unmarshal the event
	var finScopeEngineEvent FinScopeEngineEvent
	err := json.Unmarshal(event, &finScopeEngineEvent)
	if err != nil {
		return err
	}

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
		handleRequest(context.Background(), json.RawMessage(`{"execute": []}`))
	} else {
		lambda.Start(handleRequest)
	}
}
