package services

import (
	"context"
	"errors"

	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
	"gorm.io/gorm"
)

type FactorService struct {
	logger   logging.Logger
	cfg      *config.Config
	database *gorm.DB
}

func NewFactorService(cfg *config.Config) *FactorService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &FactorService{
		cfg:      cfg,
		database: database,
		logger:   logger,
	}
}

func (f *FactorService) Create(ctx context.Context, userId int, request dto.CreateFactor) (*dto.FactorProductResponse, error) {
	tx := f.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	code := helper.GenerateFactorCode()
		var postalCost int
		if request.PostalCost == 0 {
			postalCost = 35000
		} else {
			postalCost = request.PostalCost
		}

	newFactor := models.Factors{
		Code:       code,
		UserID:     userId,
		Status:     models.CREATED,
		OffPercent: request.OffPercent,
		Description: request.Description,
		PostalCost: postalCost,	
	}

	err := tx.Create(&newFactor).Error
	if err != nil {
		return nil, err
	}

	for _, product := range request.Products {
		err := tx.Create(&models.FactorProducts{
			ProductID: int(product.ID),
			FactorID:  int(newFactor.ID),
			Count:     product.Count,
		}).Error
		if err != nil {
			return nil, err
		}
	}
	tx.Commit()

	factorResponse, err := f.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return factorResponse, nil
}

func (f *FactorService) GetAll(ctx context.Context, req *dto.PaginationInput, userId int) (*dto.PagedList[dto.FactorResponse], error) {
	tx := f.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var factors []models.Factors
	err := tx.Where("user_id = ?", userId).Order("created_at desc").Find(&factors).Error
	if err != nil {
		return nil, err
	}

	totalRecords := len(factors)
	pageNumber := req.PageNumber
	pageSize := req.PageSize

	totalPages := (totalRecords + pageSize - 1) / pageSize
	hasPreviousPage := pageNumber > 1
	hasNextPage := pageNumber < totalPages

	startIndex := (pageNumber - 1) * pageSize
	endIndex := startIndex + pageSize

	if endIndex > totalRecords {
		endIndex = totalRecords
	}

	paginatedModel := factors[startIndex:endIndex]
	paginatedResponse := make([]dto.FactorResponse, len(paginatedModel))

	for i, factor := range paginatedModel {
		paginatedResponse[i] = dto.FactorResponse{
			ID:         factor.ID,
			Code:       factor.Code,
			OffPercent: factor.OffPercent,
			PostalCost: factor.PostalCost,
			Status:     f.ConvertStringToStatus(factor.Status),
			Description: factor.Description,
		}
	}

	pagedList := dto.PagedList[dto.FactorResponse]{
		PageNumber:      req.PageNumber,
		TotalRows:       int64(totalRecords),
		TotalPages:      int(totalPages),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		Items:           &paginatedResponse,
	}

	return &pagedList, nil
}

func (f *FactorService) Update(ctx context.Context, factorID int, request dto.UpdateFactor) error {
	tx := f.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var existingFactor models.Factors
	err := tx.First(&existingFactor, factorID).Error
	if err != nil {
		return err
	}

	if request.OffPercent != nil {
		existingFactor.OffPercent = *request.OffPercent
	}
	err = tx.Save(&existingFactor).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *FactorService) Delete(ctx context.Context, factorID int) error {
	tx := f.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var existingFactor models.Factors
	err := tx.First(&existingFactor, factorID).Error
	if err != nil {
		return err
	}

	err = tx.Delete(&existingFactor).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *FactorService) GetByCode(ctx context.Context, code string) (*dto.FactorProductResponse, error) {
	tx := f.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var factor models.Factors
	err := tx.Where("code = ?", code).Preload("User").First(&factor).Error
	if err != nil {
		return nil, err
	}

	var factorProducts []models.FactorProducts
	err = tx.Where("factor_id = ?", factor.ID).Preload("Product").Find(&factorProducts).Error
	if err != nil {
		return nil, err
	}

	var FactorDetailResponse *dto.FactorDetailResponse

	var factorDetail *models.FactorDetail
	err = tx.Where("factor_id = ?", factor.ID).First(&factorDetail).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if factorDetail != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		FactorDetailResponse = &dto.FactorDetailResponse{
			ID:         factorDetail.ID,
			FullName:   factorDetail.FullName,
			Mobile:     factorDetail.Mobile,
			Address:    factorDetail.Address,
			Province:   factorDetail.Province,
			City:       factorDetail.City,
			PostalCode: factorDetail.PostalCode,
			TrackingCode: &factorDetail.TrackingCode,
		}
	}

	var products []*dto.ProductFactorResponse

	for _, item := range factorProducts {
		products = append(products, &dto.ProductFactorResponse{
			ID:    item.Product.ID,
			Title: item.Product.Title,
			Price: float64(item.Product.Price),
			Count: item.Count,
		})
	}

	factorResponse := &dto.FactorProductResponse{
		ID:         factor.ID,
		Code:       factor.Code,
		OffPercent: factor.OffPercent,
		Description: factor.Description,
		PostalCost: factor.PostalCost,
		Status:     f.ConvertStringToStatus(factor.Status),
		Products:   products,
		Account: &dto.AccountResponse{
			Name: factor.User.FullName,
		},
		Factor: FactorDetailResponse,
	}

	return factorResponse, nil
}

func (f *FactorService) DeleteItem(ctx context.Context, factorID int, request dto.FactorItem) error {
	tx := f.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err := tx.Where("factor_id = ? AND product_id = ?", factorID, request.ID).Delete(&models.FactorProducts{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *FactorService) AddItem(ctx context.Context, factorID int, request dto.FactorItem) error {
	tx := f.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var product models.Products
	err := tx.Where("id = ?", request.ID).First(&product).Error
	if err != nil {
		return err
	}

	factorItem := &models.FactorProducts{
		ProductID: int(request.ID),
		FactorID:  factorID,
		Count:     request.Count,
	}

	if err := tx.Create(factorItem).Error; err != nil {
		return err
	}

	return nil
}

// Helper functions
func (f *FactorService) ConvertIntToStatus(status int) models.FactorStatus {
	switch status {
	case 1:
		return models.CREATED
	case 2:
		return models.PENDING
	case 3:
		return models.CANCEL
	case 4:
		return models.EXPIRED
	case 5:
		return models.NOTPAID
	case 6:
		return models.ACCEPTED
	case 7:
		return models.POSTED
	case 8:
		return models.FINISH
	default:
		return models.CREATED
	}
}

func (f *FactorService) ConvertStringToStatus(status models.FactorStatus) int {
	switch status {
	case models.CREATED:
		return 1
	case models.PENDING:
		return 2
	case models.CANCEL:
		return 3
	case models.EXPIRED:
		return 4
	case models.PAID:
		return 5
	case models.ACCEPTED:
		return 6
	case models.POSTED:
		return 7
	case models.FINISH:
		return 8
	default:
		return 1
	}
}

// End of helper functions
