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

func StartCrawler(cfg *config.Config) {
	log := logging.NewLogger()

	c := colly.NewCollector(
		colly.AllowedDomains("finance.yahoo.com"),
		colly.Async(false),
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
				time.Sleep(500 * time.Millisecond)
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
		publishedDate := e.ChildText("time")

		data := models.FinanceNews{
			Title:         title,
			Content:       content,
			PublishedDate: publishedDate,
			OriginURL:     link,
		}
		log.WithFields(logrus.Fields{
			"title":          title,
			"published_date": publishedDate,
			"origin_url":     link,
		}).Debug("Inserting news into database")

		err := database.InsertNews(cfg, data)
		if err != nil {
			log.WithError(err).Error("Error inserting into database")
		} else {
			log.WithField("title", title).Info("Inserted news into database")
		}
	})

	log.Info("Starting collector")
	err := c.Visit("https://finance.yahoo.com/topic/latest-news/")
	if err != nil {
		log.WithError(err).Fatal("Error starting collector")
	}

	c.Wait()

	defer func() {
		if r := recover(); r != nil {
			log.WithField("panic", r).Error("Recovered from panic")
		}
		log.Info("Programme ended")
	}()
}
