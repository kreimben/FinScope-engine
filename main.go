package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/crawler"
)

var cfg *config.Config

func init() {
	// Load config.
	cfg = config.LoadConfig()
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	crawler.StartFinanceYahooCrawler(cfg)
	return nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		lambda.Start(handleRequest)
	} else {
		handleRequest(context.Background(), json.RawMessage{})
	}
}
