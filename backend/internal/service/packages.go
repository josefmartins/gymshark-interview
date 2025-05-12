package service

import (
	"context"
	"errors"
	"gymshark-interview/internal/model"
	"gymshark-interview/internal/storage"
	"math"
	"slices"
)

func NewPackageService(storage PackagesStorage) *Packages {
	return &Packages{
		storage: storage,
	}
}

type Packages struct {
	storage PackagesStorage
}

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

// CalculatePackages calculates the minimum amount of package units required to satisfy the requested amount of units.
func (s *Packages) CalculatePackages(ctx context.Context, productID string, units int) (*model.Package, error) {
	product, err := s.storage.GetProductWithPackageSizes(ctx, productID)
	if err != nil {
		if errors.Is(err, storage.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	if len(product.PackageSizes) == 0 {
		return nil, ErrProductWithoutPackages
	}
	return &model.Package{
		PackageUnits: calculate(units, product.PackageSizes),
	}, nil
}

func calculate(units int, packageSizes []int) []model.PackageUnit {
	slices.Sort(packageSizes) // sort ascending
	maxSize := packageSizes[len(packageSizes)-1]
	maxUnits := units + maxSize // we may need to overfill a bit

	type dpEntry struct {
		totalItems int
		packCount  int
		path       []int // indices of packageSizes
	}

	const maxInt = math.MaxInt64
	dp := make([]dpEntry, maxUnits+1)
	for i := range dp {
		dp[i] = dpEntry{totalItems: maxInt, packCount: maxInt}
	}
	dp[0] = dpEntry{totalItems: 0, packCount: 0, path: []int{}}

	for i := 1; i <= maxUnits; i++ {
		for idx, size := range packageSizes {
			if i < size {
				continue
			}
			prev := dp[i-size]
			if prev.totalItems == maxInt {
				continue
			}
			newTotal := prev.totalItems + size
			newCount := prev.packCount + 1
			newPath := append([]int{}, prev.path...)
			newPath = append(newPath, idx)

			if newTotal < dp[i].totalItems || (newTotal == dp[i].totalItems && newCount < dp[i].packCount) {
				dp[i] = dpEntry{
					totalItems: newTotal,
					packCount:  newCount,
					path:       newPath,
				}
			}
		}
	}

	// Find the best valid entry from units upwards
	best := dpEntry{totalItems: maxInt, packCount: maxInt}
	for i := units; i <= maxUnits; i++ {
		entry := dp[i]
		if entry.totalItems < best.totalItems || (entry.totalItems == best.totalItems && entry.packCount < best.packCount) {
			best = entry
		}
	}

	// Count usage of each pack size
	counts := map[int]int{}
	for _, idx := range best.path {
		counts[packageSizes[idx]]++
	}

	// Build result slice
	res := []model.PackageUnit{}
	for size, count := range counts {
		res = append(res, model.PackageUnit{
			Size:   size,
			Amount: count,
		})
	}

	return res
}
