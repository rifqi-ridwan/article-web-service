package article

import (
	"article-web-service/internal/entity"
	"context"
	"strings"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	Insert(ctx context.Context, article *entity.Article) error
	FindAll(ctx context.Context) ([]entity.Article, error)
	FindByParams(ctx context.Context, fields []string, values []interface{}) ([]entity.Article, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) ArticleRepository {
	return &repository{db}
}

func (r *repository) Insert(ctx context.Context, article *entity.Article) error {
	result := r.db.Create(article)
	return result.Error
}

func (r *repository) FindAll(ctx context.Context) ([]entity.Article, error) {
	var articles []entity.Article
	result := r.db.Order("created_at desc").Find(&articles)
	return articles, result.Error
}

func (r *repository) FindByParams(ctx context.Context, fields []string, values []interface{}) ([]entity.Article, error) {
	if len(fields) == 0 && len(values) == 0 {
		return r.FindAll(ctx)
	}

	var articles []entity.Article
	result := r.db.Order("created_at desc").Where(strings.Join(fields, " AND "), values...).Find(&articles)
	return articles, result.Error
}
