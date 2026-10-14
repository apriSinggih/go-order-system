package service

import (
	"context"

	"go-order-system/internal/entity"
	"go-order-system/internal/repository"
)

type ProductService interface {
	GetAllProducts(ctx context.Context) ([]entity.Product, error)
	GetProductByID(ctx context.Context, id int64) (*entity.Product, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	return s.productRepo.GetAll(ctx)
}

func (s *productService) GetProductByID(ctx context.Context, id int64) (*entity.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}