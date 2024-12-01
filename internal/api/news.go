package api

import (
	"context"

	"github.com/kreimben/FinScope-engine/internal/database"
)

// News represents the GraphQL type for news
type News struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	PublishedDate string `json:"publishedDate"`
	OriginURL     string `json:"url"`
}

// NewsQuery represents the root query type for news
type NewsQuery struct {
	db *database.DB
}

func NewNewsQuery(db *database.DB) *NewsQuery {
	return &NewsQuery{db: db}
}

// StockNews handles the stockNews query
func (r *NewsQuery) StockNews(ctx context.Context, name *string, ticker *string) ([]*News, error) {
	var query string
	if name != nil {
		query = *name
	} else if ticker != nil {
		query = *ticker
	}

	// Get news from database
	dbNews, err := r.db.GetNewsByQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	// Convert database news to GraphQL news
	var news []*News
	for _, n := range dbNews {
		news = append(news, &News{
			Title:         n.Title,
			Content:       n.Content,
			PublishedDate: n.PublishedDate.String(),
			OriginURL:     n.OriginURL,
		})
	}

	return news, nil
}
