package services

import (
	"context"

	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
	"gorm.io/gorm"
)


type FactorService struct {
	logger       logging.Logger
	cfg          *config.Config
	database     *gorm.DB
}

func NewFactorService(cfg *config.Config) *FactorService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &FactorService{
		cfg:          cfg,
		database:     database,
		logger:       logger,
	}
}


func (f *FactorService) Create(ctx context.Context, userId int, request dto.CreateFactor) error {
    tx := f.database.WithContext(ctx).Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }()

    code := helper.GenerateFactorCode()

    newFactor := models.Factors{
        Code:         code,
        UserID:      userId,
        OffPercent:   request.OffPercent,
        Status:       models.CREATED,
        FinalPrice:   request.FinalPrice,
        PaymentMethod: models.PaymentMethod(models.NOTPAID),
    }

    err := tx.Create(&newFactor).Error
    if err != nil {
        return err
    }

    for _, product := range request.Products {
        err := tx.Create(&models.FactorProducts{
            ProductID: int(product.ID),
            FactorID:  int(newFactor.ID),
        }).Error
        if err != nil {
            return err
        }
    }

    return nil
}

func (f *FactorService) GetAll(ctx context.Context, userId int) ([]*dto.FactorResponse , error) {
    tx := f.database.WithContext(ctx).Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }()


    var factors []models.Factors
    err := tx.Where("user_id = ?", userId).Find(&factors).Error
    if err != nil {
        return nil, err
    }

    var factorResponses []*dto.FactorResponse
    for _, factor := range factors {
        factorResponses = append(factorResponses, &dto.FactorResponse{
            ID:          factor.ID,
            Code:        factor.Code,
            OffPercent:  factor.OffPercent,
            Status:      f.ConvertStringToStatus(factor.Status),
        })
    }

    return factorResponses, nil
}




// Helper functions
func (f *FactorService) ConvertIntToStatus(status int) models.FactorStatus {

	switch status {
	case 1 : 
		return models.CREATED
	case 2 :
		return models.PENDING
	case 3 :
		return models.PAID
	case 4 :
		return models.NOTPAID
	case 5 :
		return models.CANCEL
	default:
		return models.CREATED
	}
}

func (f *FactorService) ConvertStringToStatus(status models.FactorStatus) int {
        switch status {
        case models.CREATED : 
            return 1
        case models.PENDING :
            return 2
        case models.PAID :
            return 3
        case models.NOTPAID :
            return 4
        case models.CANCEL :
            return 5
        default:
            return 1
        }
}
// End of helper functions