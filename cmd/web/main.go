package main

import (
	"article-web-service/internal/article"
	"article-web-service/internal/cache"
	"article-web-service/internal/entity"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "article-web-service/internal/pkg/loadenv"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DATABASE"), os.Getenv("POSTGRES_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed connect to database: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	db.AutoMigrate(&entity.Article{})

	cache := cache.NewCache(redisClient)
	articleRepo := article.NewRepository(db)
	articleService := article.NewService(articleRepo, cache)
	articleAPI := article.NewAPIHandler(articleService)

	e.GET("/healthz", healthzHandler)

	articlesGroup := e.Group("/articles")
	articlesGroup.GET("", articleAPI.Search)
	articlesGroup.GET("/:id", articleAPI.FindByID)
	articlesGroup.POST("", articleAPI.Store)

	serverAddr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Fatal(e.Start(serverAddr))
}

func healthzHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "ok"})
}
