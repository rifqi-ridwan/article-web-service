package article

import (
	"article-web-service/internal/entity"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type apiHandler struct {
	service ArticleService
}

func NewAPIHandler(s ArticleService) *apiHandler {
	return &apiHandler{s}
}

func (h *apiHandler) Search(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	author := c.QueryParam("author")

	articles, err := h.service.Search(c.Request().Context(), keyword, author)
	if err != nil && err != entity.ErrNotFound {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"data": articles})
}

func (h *apiHandler) FindByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		if err == nil {
			err = errors.New("id is empty")
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}

	article, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		if err == entity.ErrNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"data": ""})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"data": article})
}

func (h *apiHandler) Store(c echo.Context) error {
	var createRequest articleCreateRequest
	err := c.Bind(&createRequest)
	if err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, echo.Map{"message": err.Error()})
	}

	article, err := createRequest.BuildArticleStruct()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}

	err = h.service.Store(c.Request().Context(), &article)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"data": article})
}
