package article_test

import (
	"article-web-service/internal/article"
	"article-web-service/internal/entity"
	"article-web-service/internal/test/mocks"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	fakeArticle   entity.Article
	fakeArticles  []entity.Article
	mockService   *mocks.ArticleService
	e             *echo.Echo
	expected      string
	fakeError     error
	expectedError string
)

func setupTestCase(t *testing.T) func(t *testing.T) {
	err := faker.FakeData(&fakeArticle)
	assert.NoError(t, err)
	fakeArticles = append(fakeArticles, fakeArticle)

	mockService = new(mocks.ArticleService)
	e = echo.New()
	return func(t *testing.T) {
		fakeArticle = entity.Article{}
		fakeArticles = []entity.Article{}
		mockService = &mocks.ArticleService{}
		e = &echo.Echo{}
		expected = ""
		fakeError = nil
		expectedError = ""
	}
}

func TestSearch(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)

	expected = createExpected(entity.Article{}, fakeArticles)
	cases := []struct {
		name     string
		keyword  string
		author   string
		expected string
		fakedata []entity.Article
	}{
		{"searchwithoutquery", "", "", expected, fakeArticles},
		{"searchwithquery", "testkeyword", "testauthor", expected, fakeArticles},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockService.On("Search", mock.Anything, tc.keyword, tc.author).Return(tc.fakedata, nil)

			req := httptest.NewRequest(http.MethodGet, "/articles", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.QueryParams().Add("keyword", tc.keyword)
			c.QueryParams().Add("author", tc.author)

			handler := article.NewAPIHandler(mockService)

			if assert.NoError(t, handler.Search(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, expected, rec.Body.String())
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestSearchError(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)

	fakeError = errors.New("error")
	resultError := echo.Map{"message": fakeError.Error()}
	expectedByte, err := json.Marshal(&resultError)
	assert.NoError(t, err)
	expectedError = string(expectedByte) + "\n"
	mockService.On("Search", mock.Anything, "", "").Return(nil, fakeError)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := article.NewAPIHandler(mockService)

	if assert.NoError(t, handler.Search(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, expectedError, rec.Body.String())
		mockService.AssertExpectations(t)
	}
}

func TestFindByID(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)

	expected = createExpected(fakeArticle, []entity.Article{})
	mockService.On("FindByID", mock.Anything, fakeArticle.ID).Return(fakeArticle, nil)

	req := httptest.NewRequest(http.MethodGet, "/articles/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	stringID := fmt.Sprintf("%d", fakeArticle.ID)
	c.SetParamValues(stringID)

	handler := article.NewAPIHandler(mockService)

	if assert.NoError(t, handler.FindByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
		mockService.AssertExpectations(t)
	}
}

func TestFindByIDError(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)

	expectedError := errors.New("failed to get data")
	expected = createExpected(fakeArticle, []entity.Article{})
	cases := []struct {
		name       string
		id         int
		err        error
		sendID     bool
		statusCode int
	}{
		{"emptyID", 0, nil, true, http.StatusBadRequest},
		{"notSendID", 0, nil, false, http.StatusBadRequest},
		{"NotFound", fakeArticle.ID, entity.ErrNotFound, true, http.StatusNotFound},
		{"InternalError", fakeArticle.ID, expectedError, true, http.StatusInternalServerError},
	}

	mockService = &mocks.ArticleService{}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/articles/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tc.sendID {
				c.SetParamNames("id")
				stringID := fmt.Sprintf("%d", tc.id)
				c.SetParamValues(stringID)
			}

			if tc.err != nil {
				mockService = new(mocks.ArticleService)
				mockService.On("FindByID", mock.Anything, fakeArticle.ID).Return(entity.Article{}, tc.err)
			}

			handler := article.NewAPIHandler(mockService)

			if assert.NoError(t, handler.FindByID(c)) {
				assert.Equal(t, tc.statusCode, rec.Code)
				if tc.err != nil {
					mockService.AssertExpectations(t)
				}
			}
		})
	}
}

func TestStore(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)
	newFakeArticle := fakeArticle
	newFakeArticle.ID = 0
	newFakeArticle.CreatedAt = time.Time{}
	newFakeArticle.UpdatedAt = time.Time{}

	expected = createExpected(newFakeArticle, []entity.Article{})
	mockService.On("Store", mock.Anything, &newFakeArticle).Return(nil)

	httpReqBody := fmt.Sprintf(`{"title":"%s", "author":"%s", "content":"%s"}`, fakeArticle.Title, fakeArticle.Author, fakeArticle.Content)
	req := httptest.NewRequest(http.MethodPost, "/articles", strings.NewReader(httpReqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := article.NewAPIHandler(mockService)

	if assert.NoError(t, handler.Store(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
		mockService.AssertExpectations(t)
	}
}

func TestStoreError(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)
	newFakeArticle := fakeArticle
	newFakeArticle.ID = 0
	newFakeArticle.CreatedAt = time.Time{}
	newFakeArticle.UpdatedAt = time.Time{}

	errExpected := errors.New("failed to store data")

	cases := []struct {
		name       string
		title      string
		statusCode int
		failedBind bool
		err        error
	}{
		{"failedBind", newFakeArticle.Title, http.StatusUnsupportedMediaType, true, nil},
		{"failedValidateArticle", "", http.StatusBadRequest, false, nil},
		{"failedStore", newFakeArticle.Title, http.StatusInternalServerError, false, errExpected},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err != nil {
				mockService.On("Store", mock.Anything, &newFakeArticle).Return(tc.err)
			}

			httpReqBody := fmt.Sprintf(`{"title":"%s", "author":"%s", "content":"%s"}`, tc.title, newFakeArticle.Author, newFakeArticle.Content)
			req := httptest.NewRequest(http.MethodPost, "/articles", strings.NewReader(httpReqBody))
			if !tc.failedBind {
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := article.NewAPIHandler(mockService)

			if assert.NoError(t, handler.Store(c)) {
				assert.Equal(t, tc.statusCode, rec.Code)
				if tc.err != nil {
					mockService.AssertExpectations(t)
				}
			}
		})
	}
}

func createExpected(article entity.Article, articles []entity.Article) string {
	var resultData echo.Map
	if len(articles) != 0 {
		resultData = echo.Map{"data": articles}
	} else {
		resultData = echo.Map{"data": article}
	}

	expectedByte, _ := json.Marshal(&resultData)
	return string(expectedByte) + "\n"
}
