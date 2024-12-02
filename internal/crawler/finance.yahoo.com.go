package crawler

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/kreimben/FinScope-engine/internal/config"
	"github.com/kreimben/FinScope-engine/internal/database"
	"github.com/kreimben/FinScope-engine/internal/models"
	"github.com/kreimben/FinScope-engine/pkg/logging"
	"github.com/kreimben/FinScope-engine/pkg/utils"
	"github.com/sirupsen/logrus"
)

func StartFinanceYahooCrawler(cfg *config.Config) {
	log := logging.Logger

	c := colly.NewCollector(
		colly.AllowedDomains("finance.yahoo.com"),
		colly.Async(true),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if utils.HasPathAfterNews(link) {
			exists, err := database.CheckURLExists(cfg, link)
			if err != nil {
				log.WithError(err).Error("Error checking URL in database")
				return
			}
			if !exists {
				e.Request.Visit(link)
			} else {
				log.WithField("link", link).Debug("URL already visited")
			}
		}
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		link := e.Request.URL.String()
		if link == "https://finance.yahoo.com/topic/latest-news/" {
			return
		}

		// Remove ads links.
		e.DOM.Find("strong").Remove()

		title := e.ChildText("h1.cover-title")
		content := e.ChildText(".article .body")
		publishedDateStr := e.ChildAttr("time", "datetime")
		publishedDate, err := time.Parse(time.RFC3339, publishedDateStr)
		if err != nil {
			log.WithError(err).Error("Error parsing published date")
			return
		}

		// Generate embedding for the content
		embedding, err := utils.GenerateEmbedding(cfg.HuggingFaceAPIKey, content)
		if err != nil {
			log.WithError(err).Error("Error generating embedding")
			return
		}

		data := models.FinanceNews{
			Title:         title,
			Content:       content,
			PublishedDate: publishedDate,
			OriginURL:     link,
			ContentVector: embedding,
		}
		log.WithFields(logrus.Fields{
			"title":          title,
			"published_date": publishedDate,
			"origin_url":     link,
		}).Debug("Inserting news into database")

		err = database.InsertNews(cfg, data)
		if err != nil {
			log.WithError(err).Error("Error inserting into database")
		} else {
			log.WithField("title", title).Debug("Inserted news into database")
		}
	})

	log.Info("[Finance Yahoo] Starting collector")
	err := c.Visit("https://finance.yahoo.com/topic/latest-news/")
	if err != nil {
		log.WithError(err).Error("Error starting collector")
	}

	c.Wait()

	log.Info("[Finance Yahoo] Crawling finished")
}
