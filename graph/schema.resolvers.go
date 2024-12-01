package graph

import (
	"context"

	"github.com/kreimben/FinScope-engine/graph/generated"
	"github.com/kreimben/FinScope-engine/graph/model"
)

// StockNews is the resolver for the stockNews field.
func (r *queryResolver) StockNews(ctx context.Context, name *string, ticker *string) ([]*model.News, error) {
	var query string
	if name != nil {
		query = *name
	} else if ticker != nil {
		query = *ticker
	}

	// Get news from database
	dbNews, err := r.DB.GetNewsByQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	// Convert database news to GraphQL news
	var news []*model.News
	for _, n := range dbNews {
		news = append(news, &model.News{
			Title:         n.Title,
			Content:       n.Content,
			PublishedDate: n.PublishedDate.String(),
			URL:           n.OriginURL,
		})
	}

	return news, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
