package article

import (
	"article-web-service/internal/entity"

	"github.com/go-playground/validator/v10"
)

type articleCreateRequest struct {
	Author  string `json:"author" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (acr *articleCreateRequest) BuildArticleStruct() (entity.Article, error) {
	err := acr.validate()
	if err != nil {
		return entity.Article{}, err
	}

	article := entity.Article{
		Author:  acr.Author,
		Title:   acr.Title,
		Content: acr.Content,
	}

	return article, nil
}

func (acr *articleCreateRequest) validate() error {
	validate := validator.New()

	return validate.Struct(acr)
}
