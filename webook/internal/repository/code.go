package repository

import (
	"awesomeProject/webook/internal/repository/cache"
	"golang.org/x/net/context"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

//import "github.com/rogpeppe/go-internal/cache"

type CodeRepository struct {
	cache *cache.CodeRedisCache
}

func NewCodeRepository(c *cache.CodeRedisCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}
func (repo *CodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}
func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}
