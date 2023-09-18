package services

import (
	"context"

	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
)

type ProductsService struct {
	base *BaseService[models.Products, dto.CreateProductsRequest, dto.UpdateProductsRequest, dto.ProductsResponse]
}

func NewProductsService(cfg *config.Config) *ProductsService {
	return &ProductsService{
		base: &BaseService[models.Products, dto.CreateProductsRequest, dto.UpdateProductsRequest, dto.ProductsResponse]{
			Database: db.GetDb(),
			Logger:   logging.NewLogger(cfg),
			Preloads: []preload{
				{string: "Categories"},
			},
		},
	}
}

// Create
func (s *ProductsService) CreateByUserId(ctx context.Context, req *dto.CreateProductsRequest, userId int) (*dto.ProductsResponse, error) {
	return s.base.CreateByUserId(ctx, req, userId)
}

// Update
func (s *ProductsService) Update(ctx context.Context, id int, req *dto.UpdateProductsRequest) (*dto.ProductsResponse, error) {
	return s.base.Update(ctx, id, req)
}

// Delete
func (s *ProductsService) Delete(ctx context.Context, id int) error {
	return s.base.Delete(ctx, id)
}

// Get By Id
func (s *ProductsService) GetById(ctx context.Context, id int) (*dto.ProductsResponse, error) {
	return s.base.GetById(ctx, id)
}

// Get By Filter
func (s *ProductsService) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter) (*dto.PagedList[dto.ProductsResponse], error) {
	return s.base.GetByFilter(ctx, req)
}

// Get By Filter
func (s *ProductsService) GetAll(ctx context.Context, req *dto.PaginationInputWithFilter) (*dto.PagedList[dto.ProductsResponse], error) {
	return s.base.GetByFilter(ctx, req)
}
