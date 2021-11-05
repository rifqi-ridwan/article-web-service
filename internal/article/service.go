package article

import (
	"article-web-service/internal/entity"
	"context"
	"fmt"
)

type ArticleService interface {
	Search(ctx context.Context, article *entity.Article) ([]entity.Article, error)
	Create(ctx context.Context, article *entity.Article) error
}

type service struct {
	repo ArticleRepository
}

func NewService(r ArticleRepository) ArticleService {
	return &service{r}
}

func (s *service) Search(ctx context.Context, article *entity.Article) ([]entity.Article, error) {
	var fields []string
	var values []interface{}

	if article.Title == "" && article.Content == "" {
		return s.getAll(ctx)
	}

	if article.Title != "" {
		fields = append(fields, "title LIKE ?")
		title := fmt.Sprintf("%%%s%%", article.Title)
		values = append(values, title)
	}

	if article.Content != "" {
		fields = append(fields, "content LIKE ?")
		content := fmt.Sprintf("%%%s%%", article.Content)
		values = append(values, content)
	}

	return s.repo.FindByParams(ctx, fields, values)
}

func (s *service) Create(ctx context.Context, article *entity.Article) error {
	return s.repo.Insert(ctx, article)
}

func (s *service) getAll(ctx context.Context) ([]entity.Article, error) {
	return s.repo.FindAll(ctx)
}
