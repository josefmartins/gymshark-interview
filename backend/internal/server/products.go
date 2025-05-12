package server

import (
	"context"
	"errors"
	"gymshark-interview/internal/model"
	"gymshark-interview/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type ProductsService interface {
	List(ctx context.Context) ([]model.Product, error)
	Create(ctx context.Context, product model.Product) (*model.Product, error)
	DeleteByID(ctx context.Context, id string) error
}

func (s *Server) ListProducts(ctx context.Context, _ *ListProductsRequest) (*ListProductsResponse, error) {
	products, err := s.productService.List(ctx)
	if err != nil {
		return nil, err
	}

	data := make([]ProductResponseBody, len(products))
	for i, product := range products {
		data[i] = convertProductToResponseBody(product)
	}

	return &ListProductsResponse{
		Body: ListProductsResponseBody{
			Data: data,
		},
	}, nil
}

func (s *Server) CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductResponse, error) {
	for _, packageSize := range req.Body.PackageSizes {
		if packageSize < 1 {
			return nil, huma.Error400BadRequest("invalid package size")
		}
	}

	product, err := s.productService.Create(ctx, model.Product{
		Name:         req.Body.Name,
		PackageSizes: req.Body.PackageSizes,
	})
	if err != nil {
		if errors.Is(err, service.ErrConstraintViolation) {
			return nil, huma.Error400BadRequest("constraint violation")
		}
		return nil, err
	}

	return &CreateProductResponse{
		Body: ProductResponseBody{
			ID:           product.ID,
			Name:         product.Name,
			PackageSizes: product.PackageSizes,
		},
	}, nil
}

func (s *Server) DeleteProductByID(ctx context.Context, req *DeleteProductByIDRequest) (*DeleteProductByIDResponse, error) {
	err := s.productService.DeleteByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &DeleteProductByIDResponse{}, nil
}

func convertProductToResponseBody(product model.Product) ProductResponseBody {
	return ProductResponseBody{
		ID:           product.ID,
		Name:         product.Name,
		PackageSizes: product.PackageSizes,
	}

}
