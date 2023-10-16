package services

import (
	"context"
	"errors"

	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
	"gorm.io/gorm"
)

type TransactionsService struct {
	logger       logging.Logger
	cfg          *config.Config
	tokenService *TokenService
	database     *gorm.DB
}

func NewTransactionsService(cfg *config.Config) *TransactionsService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &TransactionsService{
		cfg:          cfg,
		database:     database,
		logger:       logger,
		tokenService: NewTokenService(cfg),
	}
}



func (s *TransactionsService) GetByFilter(ctx context.Context, userId int) ([]*dto.TransactionsResponse, error) {
	tx := s.database.WithContext(ctx).Begin()
	defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }()

	var transaction []*models.Transactions

	err := tx.Where("user_id = ?",userId).Find(&transaction).Error
	if err != nil && errors.Is(err,gorm.ErrRecordNotFound) {
		return nil, err
	}
	
	responses := make([]*dto.TransactionsResponse,0)
	for _,item := range transaction{
		var title string

		if item.TransactionType == models.SALES  {
			title = "پرداخت مشتری"
		}else if item.TransactionType == models.WITHDRAW {
			title = "برداشت از کیف پول"
		}else {
			title = "سود همکاری"
		}
		responses = append(responses, &dto.TransactionsResponse{
			Title: title,
			Message: item.Description,
			Amount: item.Amount,
		})
	}
	return responses,nil	

}
