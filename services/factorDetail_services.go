package services

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/common"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
	"gorm.io/gorm"
)

type FactorDetailService struct {
	base *BaseService[models.FactorDetail, dto.CreateFactorDetailRequest, dto.UpdateFactorDetailRequest, dto.FactorDetailResponse]
}

func NewFactorDetailService(cfg *config.Config) *FactorDetailService {
	return &FactorDetailService{
		base: &BaseService[models.FactorDetail, dto.CreateFactorDetailRequest, dto.UpdateFactorDetailRequest, dto.FactorDetailResponse]{
			Database: db.GetDb(),
			Config:  cfg,
			Logger:   logging.NewLogger(cfg),
		},
	}
}

// Create
func (s *FactorDetailService) Create(ctx context.Context, req *dto.CreateFactorDetailRequest) (*dto.FactorDetailResponse, error) {
	return s.base.Create(ctx, req)
}

// Update
func (s *FactorDetailService) Update(ctx context.Context, id int, req *dto.UpdateFactorDetailRequest) (*dto.FactorDetailResponse, error) {
	return s.base.Update(ctx, id, req)
}

// Get By Id
func (s *FactorDetailService) GetById(ctx context.Context, id int) (*dto.FactorDetailResponse, error) {
	return s.base.GetById(ctx, id)
}

// Get By Filter
func (s *FactorDetailService) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter) (*dto.PagedList[dto.FactorDetailResponse], error) {
	return s.base.GetByFilter(ctx, req)
}

func (s *FactorDetailService) FactorPayment(ctx *gin.Context, req *dto.FactorPayment) error {
	tx := s.base.Database.WithContext(ctx).Begin()
	defer func() {
		if tx.Error == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()


	var factor models.Factors
	if err := tx.First(&factor, req.FactorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}


	
	randomName :=  common.GenerateRandomWord() + filepath.Ext(req.Image.Filename)

	dstPath := filepath.Join("uploads", randomName)

	url := s.base.Config.Server.AppUrl + "/api/v1/files/" + randomName

	err := ctx.SaveUploadedFile(req.Image, dstPath)
	if err != nil {
		return err
	}


	FactorPayment := models.FactorPayment{
		FactorID: req.FactorID,
		FactorImage: url,
		PaymentMethod: models.PaymentMethod(req.PaymentMethod),
	}

	err = tx.Create(&FactorPayment).Error
	if err != nil {
		return err
	}
	
	return nil
}