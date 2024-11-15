package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

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
	ticker := time.NewTicker(time.Minute * 1) // 1분마다 크롤링 수행
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			crawler.StartFinanceYahooCrawler(cfg)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		lambda.Start(handleRequest)
	} else {
		handleRequest(context.Background(), json.RawMessage{})
	}
}
