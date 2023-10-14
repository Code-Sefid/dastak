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
func (s *BankService) CreateOrUpdate(ctx context.Context, req *dto.CreateBankRequest, userId int) (*dto.BankResponse, error) {
    tx := s.base.Database.WithContext(ctx).Begin()
    var bank models.BankAccounts

    // Defer the transaction handling
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }()

    // Attempt to find a bank account for the given user
    err := tx.Where("user_id = ?", userId).First(&bank).Error
    if err != nil {
        if !errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, err
        }

        // Create a new bank account if it doesn't exist
        bank = models.BankAccounts{
            UserID:      userId,
            CardName:   req.CardName,
            CardNumber: req.CardNumber,
        }

        if err := tx.Create(&bank).Error; err != nil {
            return nil, err
        }
    } else {
        // Update the existing bank account
        bank.CardName = req.CardName
        bank.CardNumber = req.CardNumber

        if err := tx.Save(&bank).Error; err != nil {
            return nil, err
        }
    }

    var wallet models.Wallet
    err = tx.Model(&models.Wallet{}).Where("user_id = ?", userId).First(&wallet).Error
    if err != nil {
        if !errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, err
        }

        wallet = models.Wallet{
            UserId:      userId,
            Amount: 0,
        }

        if err := tx.Create(&wallet).Error; err != nil {
            return nil, err
        }
    }
    return &dto.BankResponse{
        CardName:   bank.CardName,
        CardNumber: bank.CardNumber,
        Amount:    wallet.Amount,
    }, nil
}


// Update
func (s *BankService) Update(ctx context.Context, id int, userID int ,req *dto.UpdateBankRequest) (*dto.BankResponse, error) {
	return s.base.Update(ctx, id,userID, req)
}

// Get By Id
func (s *BankService) GetByUserId(ctx context.Context, userId int) (*dto.BankResponse, error) {
    tx := s.base.Database.WithContext(ctx).Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }()

    var wallet models.Wallet
    err := tx.Model(&models.Wallet{}).Where("user_id = ?", userId).First(&wallet).Error
    if err != nil {
        if !errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, err
        }

        wallet = models.Wallet{
            UserId:      userId,
            Amount: 0,
        }

        if err := tx.Create(&wallet).Error; err != nil {
            return nil, err
        }
    }
	response , err :=  s.base.GetByUserId(ctx, userId)
    response.Amount = wallet.Amount
    return response , err
}
