package server

import (
	"context"
	"errors"
	"gymshark-interview/internal/model"
	"gymshark-interview/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type PackagesService interface {
	AddPackageSize(ctx context.Context, productID string, size int) (*model.Product, error)
	RemovePackageSize(ctx context.Context, productID string, size int) (*model.Product, error)
	CalculatePackages(ctx context.Context, productID string, units int) (*model.Package, error)
}

func (s *Server) AddPackageSize(ctx context.Context, req *AddPackageSizeRequest) (*AddPackageSizeResponse, error) {
	if req.PackageSize < 1 {
		return nil, huma.Error400BadRequest("invalid package size")
	}

	product, err := s.packagesService.AddPackageSize(ctx, req.ProductID, req.PackageSize)
	if err != nil {
		if errors.Is(err, service.ErrConstraintViolation) {
			return nil, huma.Error400BadRequest("constraint violation")
		} else if errors.Is(err, service.ErrProductNotFound) {
			return nil, huma.Error404NotFound("product not found")
		}
		return nil, err
	}

	return &AddPackageSizeResponse{
		Body: ProductResponseBody{
			ID:           product.ID,
			Name:         product.Name,
			PackageSizes: product.PackageSizes,
		},
	}, nil
}

func (s *Server) RemovePackageSize(ctx context.Context, req *RemovePackageSizeRequest) (*RemovePackageSizeResponse, error) {
	product, err := s.packagesService.RemovePackageSize(ctx, req.ProductID, req.PackageSize)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return nil, huma.Error404NotFound("product not found")
		}
		return nil, err
	}

	return &RemovePackageSizeResponse{
		Body: ProductResponseBody{
			ID:           product.ID,
			Name:         product.Name,
			PackageSizes: product.PackageSizes,
		},
	}, nil
}

func (s *Server) CalculatePackages(ctx context.Context, req *CalculatePackageSizeRequest) (*CalculatePackageSizeResponse, error) {
	if req.ProductUnits < 1 {
		return nil, huma.Error400BadRequest("invalid units request")
	}

	pack, err := s.packagesService.CalculatePackages(ctx, req.ProductID, req.ProductUnits)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return nil, huma.Error404NotFound("product not found")
		} else if errors.Is(err, service.ErrProductWithoutPackages) {
			return nil, huma.Error400BadRequest("product has no available package sizes")
		}
		return nil, err
	}

	return &CalculatePackageSizeResponse{
		Body: CalculatePackageSizeResponseBody{
			Packages: convertPackages(*pack),
		},
	}, nil
}

func convertPackages(pack model.Package) []PackageResponseBody {
	res := make([]PackageResponseBody, len(pack.PackageUnits))
	for i, packageUnit := range pack.PackageUnits {
		res[i] = PackageResponseBody{
			Amount: packageUnit.Amount,
			Size:   packageUnit.Size,
		}
	}
	return res
}
