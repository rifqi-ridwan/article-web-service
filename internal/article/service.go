package article

import (
	"article-web-service/internal/entity"
	"context"
	"fmt"
)

type ArticleService interface {
	Search(ctx context.Context, article *entity.Article) ([]entity.Article, error)
	FindByID(ctx context.Context, id int) (entity.Article, error)
	Store(ctx context.Context, article *entity.Article) error
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

	if article.Title == "" && article.Content == "" && article.Author == "" {
		return s.findAll(ctx)
	}

	if article.Author != "" {
		fields = append(fields, "author = ?")
		values = append(values, article.Author)
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

func (s *service) FindByID(ctx context.Context, id int) (entity.Article, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Store(ctx context.Context, article *entity.Article) error {
	return s.repo.Store(ctx, article)
}

func (s *service) findAll(ctx context.Context) ([]entity.Article, error) {
	return s.repo.FindAll(ctx)
}
