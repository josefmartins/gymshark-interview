package service

import (
	"context"
	"gymshark-interview/internal/model"
)

type mockPackageStorage struct {
	wantRes interface{}
	wantErr error
}

func (m *mockPackageStorage) GetProductWithPackageSizes(ctx context.Context, id string) (*model.Product, error) {
	if m.wantErr != nil {
		return nil, m.wantErr
	}
	return (m.wantRes).(*model.Product), nil
}
func (m *mockPackageStorage) AddPackageSize(ctx context.Context, productId string, size int) error {
	return m.wantErr
}
func (m *mockPackageStorage) RemovePackageSize(ctx context.Context, productId string, size int) error {
	return m.wantErr
}

type mockProductStorage struct {
	wantRes interface{}
	wantErr error
}

func (m *mockProductStorage) ListProducts(ctx context.Context) ([]model.Product, error) {
	if m.wantErr != nil {
		return nil, m.wantErr
	}
	return (m.wantRes).([]model.Product), nil
}
func (m *mockProductStorage) GetProductWithPackageSizes(ctx context.Context, id string) (*model.Product, error) {
	if m.wantErr != nil {
		return nil, m.wantErr
	}
	return (m.wantRes).(*model.Product), nil
}
func (m *mockProductStorage) CreateProduct(ctx context.Context, product model.Product) (*model.Product, error) {
	if m.wantErr != nil {
		return nil, m.wantErr
	}
	return (m.wantRes).(*model.Product), nil
}
func (m *mockProductStorage) DeleteProduct(ctx context.Context, id string) error {
	return m.wantErr
}
