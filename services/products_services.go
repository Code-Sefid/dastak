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
				{string: "Category"},
			},
		},
	}
}

// Create
func (s *ProductsService) CreateByUserId(ctx context.Context, req *dto.CreateProductsRequest, userId int) (*dto.ProductsResponse, error) {
    tx := s.base.Database.WithContext(ctx).Begin()

    product := models.Products{
        UserID:    userId,
        Title:     req.Title,
        Price:     int(req.Price),
        Inventory: req.Inventory,
    }

    if req.CategoryID != nil {
        product.CategoryID = *req.CategoryID
    }

    if err := tx.Create(&product).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    var productResponse models.Products

    if err := tx.Preload("Category").First(&productResponse, product.ID).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    var categoryResponse *dto.CategoriesResponse
    if productResponse.Category != nil {
        categoryResponse = &dto.CategoriesResponse{
            ID:   productResponse.Category.ID,
            Name: productResponse.Category.Name,
        }
    }

    tx.Commit()

    return &dto.ProductsResponse{
        ID:        productResponse.ID,
        Title:     productResponse.Title,
        Price:     float64(productResponse.Price),
        Category:  categoryResponse,
        Inventory: productResponse.Inventory,
    }, nil
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
