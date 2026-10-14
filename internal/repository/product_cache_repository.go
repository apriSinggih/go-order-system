package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-order-system/internal/entity"

	"github.com/redis/go-redis/v9"
)

type productCacheRepository struct {
	repo  ProductRepository
	redis *redis.Client
}

func NewProductCacheRepository(repo ProductRepository, redis *redis.Client) ProductRepository {
	return &productCacheRepository{repo: repo, redis: redis}
}

func (r *productCacheRepository) GetAll(ctx context.Context) ([]entity.Product, error) {
	cacheKey := "products:all"
	val, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var products []entity.Product
		if err := json.Unmarshal([]byte(val), &products); err == nil {
			return products, nil
		}
	}

	products, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	bytes, _ := json.Marshal(products)
	r.redis.Set(ctx, cacheKey, bytes, 5*time.Minute)
	return products, nil
}

func (r *productCacheRepository) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	cacheKey := fmt.Sprintf("products:%d", id)
	val, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var product entity.Product
		if err := json.Unmarshal([]byte(val), &product); err == nil {
			return &product, nil
		}
	}

	product, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	bytes, _ := json.Marshal(product)
	r.redis.Set(ctx, cacheKey, bytes, 5*time.Minute)
	return product, nil
}

func (r *productCacheRepository) UpdateStock(ctx context.Context, id int64, qty int) error {
	err := r.repo.UpdateStock(ctx, id, qty)
	if err != nil {
		return err
	}
	r.redis.Del(ctx, "products:all", fmt.Sprintf("products:%d", id))
	return nil
}