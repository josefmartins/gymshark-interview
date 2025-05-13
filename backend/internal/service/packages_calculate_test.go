package service

import (
	"context"
	"errors"
	"gymshark-interview/internal/model"
	"gymshark-interview/internal/storage"
	"slices"
	"testing"

	"github.com/google/uuid"
)

// given 99 & 100 sizes - test 98, 99, 100, 198, 199, 200, 201
func TestCalculatePackages(t *testing.T) {
	mockStorage := &mockPackageStorage{
		wantRes: &model.Product{
			ID:           uuid.NewString(),
			Name:         "ABC",
			PackageSizes: []int{99, 100},
		},
	}
	service := NewPackageService(mockStorage)

	tests := []struct {
		quantity    int
		expectedRes []model.PackageUnit
	}{
		{
			quantity: 98,
			expectedRes: []model.PackageUnit{{
				Size:   99,
				Amount: 1,
			}},
		},
		{
			quantity: 99,
			expectedRes: []model.PackageUnit{{
				Size:   99,
				Amount: 1,
			}},
		},
		{
			quantity: 100,
			expectedRes: []model.PackageUnit{{
				Size:   100,
				Amount: 1,
			}},
		},
		{
			quantity: 101,
			expectedRes: []model.PackageUnit{
				{
					Size:   99,
					Amount: 2,
				},
			},
		},
		{
			quantity: 198,
			expectedRes: []model.PackageUnit{{
				Size:   99,
				Amount: 2,
			}},
		},
		{
			quantity: 199,
			expectedRes: []model.PackageUnit{
				{
					Size:   99,
					Amount: 1,
				},
				{
					Size:   100,
					Amount: 1,
				},
			},
		},
		{
			quantity: 200,
			expectedRes: []model.PackageUnit{{
				Size:   100,
				Amount: 2,
			}},
		},
	}

	for _, testCase := range tests {
		res, err := service.CalculatePackages(context.TODO(), "ABC", testCase.quantity)
		if err != nil {
			t.Fail()
		}
		slices.SortFunc(res.PackageUnits, func(a model.PackageUnit, b model.PackageUnit) int {
			if a.Size < b.Size {
				return -1
			}
			return 1
		})

		if !slices.Equal(res.PackageUnits, testCase.expectedRes) {
			t.Fail()
		}
	}

}

// given 2 3 5 Items - test 31 59 61 89 91
func TestCalculatePackagesMoreScenarios(t *testing.T) {
	mockStorage := &mockPackageStorage{
		wantRes: &model.Product{
			ID:           uuid.NewString(),
			Name:         "ABC",
			PackageSizes: []int{2, 3, 5},
		},
	}
	service := NewPackageService(mockStorage)

	tests := []struct {
		quantity    int
		expectedRes []model.PackageUnit
	}{
		{
			quantity:    31,
			expectedRes: []model.PackageUnit{{Size: 3, Amount: 2}, {Size: 5, Amount: 5}},
		},
		{
			quantity:    59,
			expectedRes: []model.PackageUnit{{Size: 2, Amount: 2}, {Size: 5, Amount: 11}},
		},
		{
			quantity:    61,
			expectedRes: []model.PackageUnit{{Size: 3, Amount: 2}, {Size: 5, Amount: 11}},
		},
		{
			quantity:    89,
			expectedRes: []model.PackageUnit{{Size: 2, Amount: 2}, {Size: 5, Amount: 17}},
		},
		{
			quantity:    91,
			expectedRes: []model.PackageUnit{{Size: 3, Amount: 2}, {Size: 5, Amount: 17}},
		},
	}

	for _, testCase := range tests {
		res, err := service.CalculatePackages(context.TODO(), "ABC", testCase.quantity)
		if err != nil {
			t.Fail()
		}
		slices.SortFunc(res.PackageUnits, func(a model.PackageUnit, b model.PackageUnit) int {
			if a.Size < b.Size {
				return -1
			}
			return 1
		})
		if !slices.Equal(res.PackageUnits, testCase.expectedRes) {
			t.Fail()
		}
	}
}

// given 23 31 53 Items - test 500000
func TestCalculatePackagesMoreRequirementsExamples(t *testing.T) {
	mockStorage := &mockPackageStorage{
		wantRes: &model.Product{
			ID:           uuid.NewString(),
			Name:         "ABC",
			PackageSizes: []int{23, 31, 53},
		},
	}
	service := NewPackageService(mockStorage)

	tests := []struct {
		quantity    int
		expectedRes []model.PackageUnit
	}{
		{
			quantity:    500000,
			expectedRes: []model.PackageUnit{{Size: 23, Amount: 2}, {Size: 31, Amount: 7}, {Size: 53, Amount: 9429}},
		},
	}

	for _, testCase := range tests {
		res, err := service.CalculatePackages(context.TODO(), "ABC", testCase.quantity)
		if err != nil {
			t.Fail()
		}
		slices.SortFunc(res.PackageUnits, func(a model.PackageUnit, b model.PackageUnit) int {
			if a.Size < b.Size {
				return -1
			}
			return 1
		})
		if !slices.Equal(res.PackageUnits, testCase.expectedRes) {
			t.Fail()
		}
	}
}

// given 250 500 1000 2000 5000 Items - test 1, 250, 251, 501, 12001
func TestCalculatePackagesRequirementsExamples(t *testing.T) {
	mockStorage := &mockPackageStorage{
		wantRes: &model.Product{
			ID:           uuid.NewString(),
			Name:         "ABC",
			PackageSizes: []int{250, 500, 1000, 2000, 5000},
		},
	}
	service := NewPackageService(mockStorage)

	tests := []struct {
		quantity    int
		expectedRes []model.PackageUnit
	}{
		{
			quantity: 1,
			expectedRes: []model.PackageUnit{{
				Size:   250,
				Amount: 1,
			}},
		},
		{
			quantity: 250,
			expectedRes: []model.PackageUnit{{
				Size:   250,
				Amount: 1,
			}},
		},
		{
			quantity: 251,
			expectedRes: []model.PackageUnit{{
				Size:   500,
				Amount: 1,
			}},
		},
		{
			quantity: 501,
			expectedRes: []model.PackageUnit{
				{
					Size:   250,
					Amount: 1,
				},
				{
					Size:   500,
					Amount: 1,
				},
			},
		},
		{
			quantity: 12001,
			expectedRes: []model.PackageUnit{
				{
					Size:   250,
					Amount: 1,
				},
				{
					Size:   2000,
					Amount: 1,
				},
				{
					Size:   5000,
					Amount: 2,
				},
			},
		},
	}

	for _, testCase := range tests {
		res, err := service.CalculatePackages(context.TODO(), "ABC", testCase.quantity)
		if err != nil {
			t.Fail()
		}
		slices.SortFunc(res.PackageUnits, func(a model.PackageUnit, b model.PackageUnit) int {
			if a.Size < b.Size {
				return -1
			}
			return 1
		})

		if !slices.Equal(res.PackageUnits, testCase.expectedRes) {
			t.Fail()
		}
	}

}

func TestCalculatePackagesOnInvalidProduct(t *testing.T) {
	mockStorage := &mockPackageStorage{wantErr: storage.ErrProductNotFound}
	service := NewPackageService(mockStorage)

	_, err := service.CalculatePackages(context.TODO(), "ABC", 100)
	if err == nil || !errors.Is(err, ErrProductNotFound) {
		t.Fail()
	}
}

func TestCalculatePackagesOnProductWithoutPackageSizes(t *testing.T) {
	mockStorage := &mockPackageStorage{wantRes: &model.Product{ID: uuid.NewString(), Name: "ABC"}}
	service := NewPackageService(mockStorage)

	_, err := service.CalculatePackages(context.TODO(), "ABC", 100)
	if err == nil || !errors.Is(err, ErrProductWithoutPackages) {
		t.Fail()
	}
}
