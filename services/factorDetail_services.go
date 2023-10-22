package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soheilkhaledabdi/dastak/api/dto"
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
			Config:   cfg,
			Logger:   logging.NewLogger(cfg),
		},
	}
}

// Create
func (s *FactorDetailService) CreateOrUpdate(ctx context.Context, req *dto.CreateFactorDetailRequest) (*dto.FactorDetailResponse, error) {
	factor := models.Factors{}
	if err := s.base.Database.WithContext(ctx).Where("code = ?", req.Code).Preload("User").First(&factor).Error; err != nil {
		return nil, err
	}

	factorDetail := models.FactorDetail{}
	err := s.base.Database.WithContext(ctx).Where("factor_id = ?", factor.ID).First(&factorDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			factorDetail = models.FactorDetail{
				FactorID:     factor.ID,
				FullName:     req.FullName,
				Mobile:       req.Mobile,
				Province:     req.Province,
				City:         req.City,
				Address:      req.Address,
				PostalCode:   *req.PostalCode,
				TrackingCode: "",
			}

			err = s.base.Database.WithContext(ctx).Create(&factorDetail).Error
			if err != nil {
				return nil, err
			}

			factor.Status = models.PENDING
			err = s.base.Database.WithContext(ctx).Save(&factor).Error
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		factorDetail.FullName = req.FullName
		factorDetail.Mobile = req.Mobile
		factorDetail.Province = req.Province
		factorDetail.City = req.City
		factorDetail.Address = req.Address
		factorDetail.PostalCode = *req.PostalCode

		err = s.base.Database.WithContext(ctx).Save(&factorDetail).Error
		if err != nil {
			return nil, err
		}
	}

	return &dto.FactorDetailResponse{
		ID:         factorDetail.ID,
		FullName:   factorDetail.FullName,
		Mobile:     factor.User.Mobile,
		City:       factorDetail.City,
		Province:   factorDetail.Province,
		Address:    factorDetail.Address,
		PostalCode: factorDetail.PostalCode,
	}, nil
}

// Update
func (s *FactorDetailService) Update(ctx context.Context, id int, userID int, req *dto.UpdateFactorDetailRequest) (*dto.FactorDetailResponse, error) {
	return s.base.Update(ctx, id, userID, req)
}

// Get By Id
func (s *FactorDetailService) GetById(ctx context.Context, id int, userID int) (*dto.FactorDetailResponse, error) {
	return s.base.GetById(ctx, id, userID)
}

// Get By Filter
func (s *FactorDetailService) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter, userId int) (*dto.PagedList[dto.FactorDetailResponse], error) {
	return s.base.GetByFilter(ctx, req, userId)
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

	// randomName := common.GenerateRandomWord() + filepath.Ext(req.Image.Filename)

	// dstPath := filepath.Join("uploads", randomName)

	// url := s.base.Config.Server.AppUrl + "/api/v1/files/" + randomName

	// err := ctx.SaveUploadedFile(req.Image, dstPath)
	// if err != nil {
	// 	return err
	// }

	FactorPayment := models.FactorPayment{
		FactorID:   req.FactorID,
		FinalPrice: req.FinalPrice,
		// FactorImage:   url,
		// PaymentMethod: models.PaymentMethod(req.PaymentMethod),
	}

	err := tx.Create(&FactorPayment).Error
	if err != nil {
		return err
	}

	return nil
}


func (s *FactorDetailService) AddTrackingCode(ctx context.Context, req *dto.AddTrackingCode) error {
	tx := s.base.Database.WithContext(ctx).Begin()
	defer func() {
		if tx.Error == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	var factor models.Factors
	if err := tx.Where("code = ?" , req.Code).Preload("User").First(&factor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	factorDetail := models.FactorDetail{}
	err := tx.Where("factor_id = ?", factor.ID).First(&factorDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	factor.Status = models.POSTED
	err = tx.Save(&factor).Error
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"status": models.POSTED,
	}

	err = tx.Model(&factorDetail).Updates(update).Error
	if err != nil {
		return err
	}

	s.SendTrackingCode(factorDetail.FullName,factor.User.FullName,req.TrackingCode,factorDetail.Mobile)

	return nil
}

// Helper Functions
func (s *FactorDetailService) SendTrackingCode(fullName, shopName, trackingCode,mobile string) error {
	url := fmt.Sprintf("http://api.payamak-panel.com/post/Send.asmx/SendByBaseNumber3?username=09135882813&password=T13Y7&text=@168063@%s;%s;%s;&to=%s", fullName, shopName,trackingCode,mobile)

	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP: %d", response.StatusCode)
	}

	return nil
}

// End Helper Functions