package services

import (
	"context"

	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
)

type CategoriesService struct {
	base *BaseService[models.Categories, dto.CreateCategoriesRequest, dto.UpdateCategoriesRequest, dto.CategoriesResponse]
}

func NewCategoriesService(cfg *config.Config) *CategoriesService {
	return &CategoriesService{
		base: &BaseService[models.Categories, dto.CreateCategoriesRequest, dto.UpdateCategoriesRequest, dto.CategoriesResponse]{
			Database: db.GetDb(),
			Logger:   logging.NewLogger(cfg),
		},
	}
}

// Create
func (s *CategoriesService) Create(ctx context.Context, req *dto.CreateCategoriesRequest, userId int) (*dto.CategoriesResponse, error) {
	return s.base.CreateByUserId(ctx, req, userId)
}

// Update
func (s *CategoriesService) Update(ctx context.Context, id int, req *dto.UpdateCategoriesRequest) (*dto.CategoriesResponse, error) {
	return s.base.Update(ctx, id, req)
}

// Delete
func (s *CategoriesService) Delete(ctx context.Context, id int) error {
	return s.base.Delete(ctx, id)
}

// Get By Id
func (s *CategoriesService) GetById(ctx context.Context, id int) (*dto.CategoriesResponse, error) {
	return s.base.GetById(ctx, id)
}

// Get By Filter
func (s *CategoriesService) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter, userId int) (*dto.PagedList[dto.CategoriesResponse], error) {
	return s.base.GetByFilter(ctx, req, userId)
}

// Get By Filter
func (s *CategoriesService) GetAll(ctx context.Context, req *dto.PaginationInputWithFilter, userId int) (*dto.PagedList[dto.CategoriesResponse], error) {
	return s.base.GetByFilter(ctx, req, userId)
}
