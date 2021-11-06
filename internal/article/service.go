package article

import (
	"article-web-service/internal/cache"
	"article-web-service/internal/entity"
	"context"
	"log"
	"strconv"
	"time"
)

const defaultDuration = 60 * time.Hour

type ArticleService interface {
	Search(ctx context.Context, keyword string, author string) ([]entity.Article, error)
	FindByID(ctx context.Context, id int) (entity.Article, error)
	Store(ctx context.Context, article *entity.Article) error
}

type service struct {
	repo  ArticleRepository
	cache cache.CacheRepository
}

func NewService(r ArticleRepository, c cache.CacheRepository) ArticleService {
	return &service{r, c}
}

func (s *service) Search(ctx context.Context, keyword string, author string) ([]entity.Article, error) {
	var fields []string
	var values []interface{}

	if keyword == "" && author == "" {
		return s.findAll(ctx)
	}

	if keyword != "" {
		s.repo.QueryBuilder(ctx, &fields, "title", "LIKE", "AND")
		values = append(values, keyword)
		s.repo.QueryBuilder(ctx, &fields, "content", "LIKE", "AND")
		values = append(values, keyword)
	}

	if author != "" {
		s.repo.QueryBuilder(ctx, &fields, "author", "=", "AND")
		values = append(values, author)
	}

	return s.repo.FindByParams(ctx, fields, values)
}

func (s *service) FindByID(ctx context.Context, id int) (entity.Article, error) {
	stringID := strconv.Itoa(id)
	article, err := s.findCacheByID(ctx, stringID)
	if err != nil && err != entity.CacheNotExist {
		return article, err
	}

	article, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return article, err
	}

	err = s.cache.WriteJSON(ctx, stringID, article, defaultDuration)
	if err != nil {
		log.Printf("failed to create cache: %v", err)
	}
	return article, nil
}

func (s *service) findCacheByID(ctx context.Context, id string) (entity.Article, error) {
	var article entity.Article
	err := s.cache.ReadJSON(ctx, id, &article)
	return article, err
}

func (s *service) Store(ctx context.Context, article *entity.Article) error {
	return s.repo.Store(ctx, article)
}

func (s *service) findAll(ctx context.Context) ([]entity.Article, error) {
	return s.repo.FindAll(ctx)
}
