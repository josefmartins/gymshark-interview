package service

import (
	"context"
	"errors"
	"gymshark-interview/internal/model"
	"gymshark-interview/internal/storage"
	"testing"
)

func TestAddPackageFailsOnStorageConstraint(t *testing.T) {
	mockStorage := &mockPackageStorage{wantErr: storage.ErrConstraintViolation}
	service := NewPackageService(mockStorage)

	_, err := service.AddPackageSize(context.TODO(), "ABC", 100)
	if err == nil || !errors.Is(err, ErrConstraintViolation) {
		t.Fail()
	}
}

func TestAddPackageOK(t *testing.T) {
	wantProduct := &model.Product{ID: "123", Name: "ABC", PackageSizes: []int{1, 2, 3, 100}}
	mockStorage := &mockPackageStorage{wantRes: wantProduct}
	service := NewPackageService(mockStorage)

	product, err := service.AddPackageSize(context.TODO(), "ABC", 100)
	if err != nil {
		t.Fail()
	}
	if product != wantProduct {
		t.Fail()
	}
}

func TestRemovePackageFailsWithError(t *testing.T) {
	mockStorage := &mockPackageStorage{wantErr: errors.New("db is unhealthy")}
	service := NewPackageService(mockStorage)

	_, err := service.AddPackageSize(context.TODO(), "ABC", 100)
	if err == nil {
		t.Fail()
	}
}

func TestRemovePackageOK(t *testing.T) {
	wantProduct := &model.Product{ID: "123", Name: "ABC", PackageSizes: []int{1, 2}}
	mockStorage := &mockPackageStorage{wantRes: wantProduct}
	service := NewPackageService(mockStorage)

	product, err := service.AddPackageSize(context.TODO(), "ABC", 3)
	if err != nil {
		t.Fail()
	}
	if product != wantProduct {
		t.Fail()
	}
}
