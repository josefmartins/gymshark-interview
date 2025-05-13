package service

import (
	"context"
	"errors"
	"gymshark-interview/internal/model"
	"gymshark-interview/internal/storage"
	"math"
	"slices"
)

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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

// get the Least Common Multiple
func getLeastCommonMultiple(nums []int) int {
	result := nums[0]
	for _, num := range nums[1:] {
		result = lcm(result, num)
	}
	return result
}

func calculate(units int, packageSizes []int) []model.PackageUnit {
	slices.Sort(packageSizes) // sort ascending
	maxSize := packageSizes[len(packageSizes)-1]

	leastCommonMultiple := getLeastCommonMultiple(packageSizes)

	// if units is bigger than at least 2 times lcm, then let's calculate an offset by using the biggest packageSize
	biggestPackageOffset := 0
	if units > leastCommonMultiple*2 {
		biggestPackageInLCM := leastCommonMultiple / maxSize

		rest := units % leastCommonMultiple
		division := units / leastCommonMultiple

		units = rest + leastCommonMultiple*2
		biggestPackageOffset = (division - 2) * biggestPackageInLCM
	}

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

	// Use biggestPackageOffset if not zero
	if biggestPackageOffset != 0 {
		counts[maxSize] += biggestPackageOffset
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
