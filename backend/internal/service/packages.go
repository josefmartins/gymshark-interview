package service

import (
	"context"
	"errors"
	"gymshark-interview/internal/model"
	"gymshark-interview/internal/storage"
)

func NewPackageService(storage PackagesStorage) *Packages {
	return &Packages{
		storage: storage,
	}
}

type Packages struct {
	storage PackagesStorage
}

var ErrProductNotFound = errors.New("product not found")

type PackagesStorage interface {
	GetProductWithPackageSizes(ctx context.Context, id string) (*model.Product, error)
	AddPackageSize(ctx context.Context, productId string, size int) error
	RemovePackageSize(ctx context.Context, productId string, size int) error
}

func (s *Packages) AddPackageSize(ctx context.Context, productID string, size int) (*model.Product, error) {
	err := s.storage.AddPackageSize(ctx, productID, size)
	if err != nil {
		if errors.Is(err, storage.ErrConstraintViolation) {
			return nil, ErrConstraintViolation
		}
		return nil, err
	}
	product, err := s.storage.GetProductWithPackageSizes(ctx, productID)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (s *Packages) RemovePackageSize(ctx context.Context, productID string, size int) (*model.Product, error) {
	err := s.storage.RemovePackageSize(ctx, productID, size)
	if err != nil {
		if errors.Is(err, storage.ErrConstraintViolation) {
			return nil, ErrConstraintViolation
		}
		return nil, err
	}
	product, err := s.storage.GetProductWithPackageSizes(ctx, productID)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}
