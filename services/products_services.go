package services

import (
	"context"
	"fmt"

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
		product.CategoryID = req.CategoryID
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create product: %v", err)
	}

	var productResponse models.Products

	if err := tx.Model(&models.Products{}).Preload("Category").Order("created_at desc").First(&productResponse, product.ID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to retrieve product after creation: %v", err)
	}

	var categoryResponse *dto.CategoriesResponse
	if productResponse.Category != nil {
		categoryResponse = &dto.CategoriesResponse{
			ID:   productResponse.Category.ID,
			Name: productResponse.Category.Name,
		}
	}

	tx.Commit()

	if req.CategoryID == nil {
		categoryResponse = nil
	}

	return &dto.ProductsResponse{
		ID:        productResponse.ID,
		Title:     productResponse.Title,
		Price:     float64(productResponse.Price),
		Category:  categoryResponse,
		Inventory: productResponse.Inventory,
	}, nil
}

// Update
func (s *ProductsService) Update(ctx context.Context, req *dto.UpdateProductsRequest, productID int) (*dto.ProductsResponse, error) {
    tx := s.base.Database.WithContext(ctx).Begin()

    var existingProduct models.Products
    if err := tx.Where("id = ?", productID).First(&existingProduct).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to find product: %v", err)
    }




    if req.Title != nil {
        existingProduct.Title = *req.Title
    }
    if req.Price != nil {
        existingProduct.Price = int(*req.Price)
    }
    if req.Inventory != nil {
        existingProduct.Inventory = *req.Inventory
    }
    if req.CategoryID != nil {
		var category models.Categories
		err := tx.Model(&models.Categories{}).Where("id = ?",req.CategoryID).First(&category).Error
		if err != nil {
			return nil,err
		}
        existingProduct.CategoryID = req.CategoryID
    }

    if err := tx.Save(&existingProduct).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to update product: %v", err)
    }

    // Load the updated product with category information
    var productResponse models.Products
    if err := tx.Model(&models.Products{}).Preload("Category").Where("id = ?", productID).First(&productResponse).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to retrieve updated product: %v", err)
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

// Delete
func (s *ProductsService) Delete(ctx context.Context, id int, userID int) error {
	return s.base.Delete(ctx, id, userID)
}

// Get By Id
func (s *ProductsService) GetById(ctx context.Context, id int, userID int) (*dto.ProductsResponse, error) {
	return s.base.GetById(ctx, id, userID)
}

// Get By Filter
func (s *ProductsService) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter, userId int) (*dto.PagedList[dto.ProductsResponse], error) {
	return s.base.GetByFilter(ctx, req, userId)
}

// Get By Filter
func (s *ProductsService) GetAll(ctx context.Context, req *dto.PaginationInputWithFilter, userId int) (*dto.PagedList[dto.ProductsResponse], error) {
	return s.base.GetByFilter(ctx, req, userId)
}
