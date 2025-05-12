package service

import (
	"context"
	"errors"
	"gymshark-interview/internal/model"
	"gymshark-interview/internal/storage"
	"slices"
	"testing"
)

func TestCreateProductsOK(t *testing.T) {
	wantRes := &model.Product{ID: "ABC", Name: "one", PackageSizes: []int{5, 10}}
	mockStorage := &mockProductStorage{wantRes: wantRes}
	service := NewProductService(mockStorage)

	product, err := service.Create(context.TODO(), model.Product{
		ID:           wantRes.ID,
		Name:         wantRes.Name,
		PackageSizes: wantRes.PackageSizes,
	})
	if err != nil {
		t.Fail()
	}
	if wantRes != product {
		t.Fail()
	}
}

func TestCreateProductConstraintStorageError(t *testing.T) {
	mockStorage := &mockProductStorage{
		wantErr: storage.ErrConstraintViolation,
	}
	service := NewProductService(mockStorage)

	_, err := service.Create(context.TODO(), model.Product{ID: "ABC", Name: "one", PackageSizes: []int{5, 10}})
	if err == nil || !errors.Is(err, ErrConstraintViolation) {
		t.Fail()
	}
}

func TestListProductsOK(t *testing.T) {
	wantRes := []model.Product{
		{ID: "ABC", Name: "one", PackageSizes: []int{5, 10}},
		{ID: "DEF", Name: "two"},
	}
	mockStorage := &mockProductStorage{wantRes: wantRes}
	service := NewProductService(mockStorage)

	products, err := service.List(context.TODO())
	if err != nil {
		t.Fail()
	}
	if !slices.EqualFunc(wantRes, products, func(a, b model.Product) bool {
		return a.ID == b.ID && a.Name == b.Name && slices.Equal(a.PackageSizes, b.PackageSizes)
	}) {
		t.Fail()
	}
}

func TestListProductsStorageError(t *testing.T) {
	mockStorage := &mockProductStorage{
		wantErr: errors.New("storage failed"),
	}
	service := NewProductService(mockStorage)

	_, err := service.List(context.TODO())
	if err == nil {
		t.Fail()
	}
}

func TestDeleteProductOK(t *testing.T) {
	mockStorage := &mockProductStorage{}
	service := NewProductService(mockStorage)

	err := service.DeleteByID(context.TODO(), "ABC")
	if err != nil {
		t.Fail()
	}
}

func TestDeleteProductStorageError(t *testing.T) {
	mockStorage := &mockProductStorage{
		wantErr: errors.New("storage failed"),
	}
	service := NewProductService(mockStorage)

	err := service.DeleteByID(context.TODO(), "ABC")
	if err == nil {
		t.Fail()
	}
}
