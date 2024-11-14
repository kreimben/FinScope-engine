package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/crawler"
)

func handleRequest(ctx context.Context, event json.RawMessage) error {
	// Load config.
	cfg := config.LoadConfig()

	// Start crawler.
	crawler.StartFinanceYahooCrawler(cfg)

	if r := recover(); r != nil {
		return fmt.Errorf("panic: %v", r)
	} else {
		return nil
	}
}

func main() {
	lambda.Start(handleRequest)
}
