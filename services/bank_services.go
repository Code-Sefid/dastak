package services

import (
	"context"

	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
)

type BankService struct {
	base *BaseService[models.BankAccounts, dto.CreateBankRequest, dto.UpdateBankRequest, dto.BankResponse]
}

func NewBankService(cfg *config.Config) *BankService {
	return &BankService{
		base: &BaseService[models.BankAccounts, dto.CreateBankRequest, dto.UpdateBankRequest, dto.BankResponse]{
			Database: db.GetDb(),
			Logger:   logging.NewLogger(cfg),
		},
	}
}

// Create
func (s *BankService) CreateByUserId(ctx context.Context, req *dto.CreateBankRequest, userId int) (*dto.BankResponse, error) {
	return s.base.CreateByUserId(ctx, req, userId)
}

// Update
func (s *BankService) Update(ctx context.Context, id int, req *dto.UpdateBankRequest) (*dto.BankResponse, error) {
	return s.base.Update(ctx, id, req)
}

// Get By Id
func (s *BankService) GetByUserId(ctx context.Context, userId int) (*dto.BankResponse, error) {
	return s.base.GetByUserId(ctx, userId)
}
