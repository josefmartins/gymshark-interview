package service

import (
	"context"
	"errors"
	"gymshark-interview/internal/model"
	"gymshark-interview/internal/storage"
)

func NewProductService(storage ProductsStorage) *Products {
	return &Products{
		storage: storage,
	}
}

type Products struct {
	storage ProductsStorage
}

type ProductsStorage interface {
	ListProducts(ctx context.Context) ([]model.Product, error)
	GetProductWithPackageSizes(ctx context.Context, id string) (*model.Product, error)
	CreateProduct(ctx context.Context, product model.Product) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

var (
	ErrProductNotFound        = errors.New("product not found")
	ErrConstraintViolation    = errors.New("constraint violation")
	ErrProductWithoutPackages = errors.New("product has no available package sizes")
)

func (s *Products) List(ctx context.Context) ([]model.Product, error) {
	return s.storage.ListProducts(ctx)
}

func (s *Products) GetByID(ctx context.Context, id string) (*model.Product, error) {
	product, err := s.storage.GetProductWithPackageSizes(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (s *Products) Create(ctx context.Context, product model.Product) (*model.Product, error) {
	res, err := s.storage.CreateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, storage.ErrConstraintViolation) {
			return nil, ErrConstraintViolation
		}
		return nil, err
	}
	return res, nil
}

func (s *Products) DeleteByID(ctx context.Context, id string) error {
	return s.storage.DeleteProduct(ctx, id)
}
