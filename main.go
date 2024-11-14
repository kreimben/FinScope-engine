package main

import (
	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/crawler"
)

func main() {
	cfg := config.LoadConfig()
	crawler.StartFinanceYahooCrawler(cfg)
}
