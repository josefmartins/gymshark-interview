package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

const (
	v1                            = "/v1"
	listProductsEndpointPath      = v1 + "/products"
	createProductEndpointPath     = v1 + "/products"
	deleteProductByIDEndpointPath = v1 + "/products/{productID}"

	modifyPackageSizeEndpointPath = v1 + "/products/{productID}/packageSizes/{packageSize}"
	calculatePackagesEndpointPath = v1 + "/products/{productID}/calculate/{productUnits}"
)

func (s *Server) declareRoutes() {
	var listProductsResponse *ListProductsResponse
	huma.Register(s.api, huma.Operation{
		OperationID:   huma.GenerateOperationID(http.MethodGet, listProductsEndpointPath, listProductsResponse),
		Summary:       "v1 - List Products",
		Method:        http.MethodGet,
		Path:          listProductsEndpointPath,
		DefaultStatus: http.StatusOK,
	}, s.ListProducts)
	var createProductResponse *CreateProductResponse
	huma.Register(s.api, huma.Operation{
		OperationID:   huma.GenerateOperationID(http.MethodPost, createProductEndpointPath, createProductResponse),
		Summary:       "v1 - Create Product",
		Method:        http.MethodPost,
		Path:          createProductEndpointPath,
		DefaultStatus: http.StatusCreated,
	}, s.CreateProduct)
	huma.Register(s.api, huma.Operation{
		Method:        http.MethodOptions,
		Path:          createProductEndpointPath,
		DefaultStatus: http.StatusNoContent,
		Hidden:        true,
	}, s.CreateProduct)
	var deleteProductResponse *DeleteProductByIDResponse
	huma.Register(s.api, huma.Operation{
		OperationID:   huma.GenerateOperationID(http.MethodDelete, deleteProductByIDEndpointPath, deleteProductResponse),
		Summary:       "v1 - Delete Product",
		Method:        http.MethodDelete,
		Path:          deleteProductByIDEndpointPath,
		DefaultStatus: http.StatusNoContent,
	}, s.DeleteProductByID)
	huma.Register(s.api, huma.Operation{
		Method:        http.MethodOptions,
		Path:          deleteProductByIDEndpointPath,
		DefaultStatus: http.StatusNoContent,
		Hidden:        true,
	}, s.DeleteProductByID)

	var addPackageResponse *AddPackageSizeResponse
	huma.Register(s.api, huma.Operation{
		OperationID:   huma.GenerateOperationID(http.MethodPost, modifyPackageSizeEndpointPath, addPackageResponse),
		Summary:       "v1 - Add Package Size",
		Method:        http.MethodPost,
		Path:          modifyPackageSizeEndpointPath,
		DefaultStatus: http.StatusCreated,
	}, s.AddPackageSize)
	var removePackageResponse *RemovePackageSizeResponse
	huma.Register(s.api, huma.Operation{
		OperationID:   huma.GenerateOperationID(http.MethodDelete, modifyPackageSizeEndpointPath, removePackageResponse),
		Summary:       "v1 - Remove Package Size",
		Method:        http.MethodDelete,
		Path:          modifyPackageSizeEndpointPath,
		DefaultStatus: http.StatusOK,
	}, s.RemovePackageSize)
	huma.Register(s.api, huma.Operation{
		Method:        http.MethodOptions,
		Path:          modifyPackageSizeEndpointPath,
		DefaultStatus: http.StatusNoContent,
		Hidden:        true,
	}, s.AddPackageSize)
	var calculatePackageResponse *CalculatePackageSizeResponse
	huma.Register(s.api, huma.Operation{
		OperationID:   huma.GenerateOperationID(http.MethodPost, calculatePackagesEndpointPath, calculatePackageResponse),
		Summary:       "v1 - Calculate Package",
		Method:        http.MethodPost,
		Path:          calculatePackagesEndpointPath,
		DefaultStatus: http.StatusOK,
	}, s.CalculatePackages)
	huma.Register(s.api, huma.Operation{
		Method:        http.MethodOptions,
		Path:          calculatePackagesEndpointPath,
		DefaultStatus: http.StatusNoContent,
		Hidden:        true,
	}, s.AddPackageSize)
}

type ListProductsRequest struct{}

type ListProductsResponse struct {
	Body ListProductsResponseBody
}

type ListProductsResponseBody struct {
	Data []ProductResponseBody
}

type ProductResponseBody struct {
	ID           string `json:"id" example:"018ef16a-31a7-7e11-a77d-78b2eea91e2f" doc:"Product ID"`
	Name         string `json:"name" example:"My First Product" doc:"Name of the Product"`
	PackageSizes []int  `json:"package_sizes,omitempty" doc:"Available Package Sizes"`
}

type CreateProductRequest struct {
	Body CreateProductRequestBody `required:"true"`
}

type CreateProductRequestBody struct {
	Name         string `json:"name" minLength:"5" required:"true" example:"My First Product" doc:"Name of the Product"`
	PackageSizes []int  `json:"package_sizes" required:"false" example:"[100]" doc:"Available Package Sizes"`
}

type CreateProductResponse struct {
	Body ProductResponseBody
}

type DeleteProductByIDRequest struct {
	ID string `path:"productID" example:"018ef16a-31a7-7e11-a77d-78b2eea91e2f" doc:"Product ID"`
}

type DeleteProductByIDResponse struct{}

type AddPackageSizeRequest struct {
	ProductID   string `path:"productID" example:"018ef16a-31a7-7e11-a77d-78b2eea91e2f" doc:"Product ID"`
	PackageSize int    `path:"packageSize" example:"250" doc:"Package Size"`
}

type AddPackageSizeResponse struct {
	Body ProductResponseBody
}

type RemovePackageSizeRequest struct {
	ProductID   string `path:"productID" example:"018ef16a-31a7-7e11-a77d-78b2eea91e2f" doc:"Product ID"`
	PackageSize int    `path:"packageSize" example:"250" doc:"Package Size"`
}

type RemovePackageSizeResponse struct {
	Body ProductResponseBody
}

type CalculatePackageSizeRequest struct {
	ProductID    string `path:"productID" example:"018ef16a-31a7-7e11-a77d-78b2eea91e2f" doc:"Product ID"`
	ProductUnits int    `path:"productUnits" example:"250" doc:"Product Units"`
}

type CalculatePackageSizeResponse struct {
	Body CalculatePackageSizeResponseBody
}

type CalculatePackageSizeResponseBody struct {
	Packages []PackageResponseBody `json:"packages" doc:"List of Packages"`
}

type PackageResponseBody struct {
	Amount int `json:"units"  example:"3" doc:"Units of Package"`
	Size   int `json:"size"  example:"250" doc:"Package Size"`
}
