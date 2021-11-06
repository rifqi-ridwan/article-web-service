package entity

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	ErrNotFound   = gorm.ErrRecordNotFound
	CacheNotExist = redis.Nil
)
