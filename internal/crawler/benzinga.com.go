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
	"golang.org/x/exp/rand"
)

var newsUrlFormat = []string{
	"https://www.benzinga.com/trading-ideas/*",
	"https://www.benzinga.com/news/*",
	"https://www.benzinga.com/[0-9][0-9]/[0-9][0-9]/*",
	"https://www.benzinga.com/markets/*",
}

func StartBenzingaCrawler(cfg *config.Config) {
	log := logging.Logger

	c := colly.NewCollector(
		colly.AllowedDomains("benzinga.com"),
		colly.Async(true),
		colly.AllowedDomains("benzinga.com", "www.benzinga.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*benzinga.com*",
		Parallelism: 1,
	})

	// On request
	uaList := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", uaList[rand.Intn(len(uaList))])
	})

	// On error
	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 429 {
			log.Println("Rate limited. Retrying after delay...")
			time.Sleep(3 * time.Second)
			r.Request.Retry()
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if !utils.ContainsURLLink(newsUrlFormat, link) {
			return
		}
		if utils.HasPathAfterNews(link) {
			exists, err := database.CheckURLExists(cfg, link)
			if err != nil {
				log.WithError(err).Error("Error checking URL in database")
				return
			}
			if !exists {
				e.Request.Visit(link)
			}
		}
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		link := e.Request.URL.String()
		if !utils.ContainsURLLink(newsUrlFormat, link) {
			return
		}
		// log.WithField("link", link).Debug("Visited URL")

		title := e.ChildText("h1.layout-title")
		if title == "" {
			log.WithField("link", link).Error("No title found")
			return
		}
		content := e.ChildText("div.article-content-body-only")
		if content == "" {
			log.WithField("link", link).Error("No content found")
			return
		}
		publishedDateStr := e.ChildText("span.date")
		if publishedDateStr == "" {
			log.WithField("link", link).Error("No published date found")
			return
		}
		publishedDate, err := time.Parse("January 2, 2006 3:04 PM", publishedDateStr)
		if err != nil {
			log.WithError(err).Error("Error parsing published date")
			return
		}

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

		err = database.InsertNews(cfg, data)
		if err != nil {
			log.WithError(err).Error("Error inserting into database")
		} else {
			log.WithField("title", title).Debug("Inserted news into database")
		}
	})

	log.Info("[Benzinga] Starting collector")
	err := c.Visit("https://www.benzinga.com/recent")
	if err != nil {
		log.WithError(err).Error("Error starting collector")
	}

	c.Wait()

	log.Info("[Benzinga] Crawling finished")
}
