package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
	"gorm.io/gorm"
)

type CheckOutService struct {
	logger   logging.Logger
	cfg      *config.Config
	database *gorm.DB
}

func NewCheckOutService(cfg *config.Config) *CheckOutService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &CheckOutService{
		cfg:      cfg,
		database: database,
		logger:   logger,
	}
}

func (c *CheckOutService) CheckOutMony(ctx context.Context, userID int, req dto.CheckOut) (*dto.Alert, bool, error) {
	tx := c.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	var user models.Users
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, false, err
	}

	var wallet models.Wallet
	err := tx.Model(&models.Wallet{}).Where("user_id = ?", user.ID).First(&wallet).Error
	if err != nil {
		tx.Rollback()
		return &dto.Alert{Message: "کیف پول شما خالی است"}, false, nil
	}

	var checkout models.CheckOutRequest
	err = tx.Where("user_id = ?", userID).Order("created_at desc").First(&checkout).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		if checkout.Status == models.PENDINGCHECKOUT {
			return &dto.Alert{Message: "شما نمیتوانید بیش از یک درخواست برداشت فعال داشته باشید"}, false, nil
		}
	}

	var GetCountTransactions int64
	err = tx.Model(&models.Transactions{}).Where("user_id = ?", userID).Where("transaction_type = ?", models.WITHDRAW).Count(&GetCountTransactions).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	if req.Amount < 500000 && GetCountTransactions == 0 {
		return &dto.Alert{
			Message: "حداقل میزان برداشت باید 500,000 تومان باشد",
		}, false, nil
	}

	if req.Amount <= 100000 && GetCountTransactions > 0 {
		return &dto.Alert{
			Message: "حداقل میزان برداشت باید 100,000 تومان باشد",
		}, false, nil
	}


	

	if req.Amount >= 50000 && wallet.Amount >= req.Amount {
		amount := wallet.Amount - req.Amount
		err = tx.Model(&models.Wallet{}).Where("user_id = ?", user.ID).Update("amount", amount).Error
		if err != nil {
			return nil, false, err
		}

		checkOut := models.CheckOutRequest{
			UserID: user.ID,
			Amount: req.Amount,
			Status: models.PENDINGCHECKOUT,
		}

		err = tx.Model(&models.CheckOutRequest{}).Create(&checkOut).Error
		if err != nil {
			return nil, false, err
		}

		transactions := models.Transactions{
			UserID:          userID,
			Description:     "درخواست برداشت شما به مبلغ " + strconv.Itoa(req.Amount) + " تومان در حال بررسی است",
			Amount:          float64(req.Amount),
			TransactionType: models.WITHDRAW,
		}

		err = tx.Model(&models.Transactions{}).Create(&transactions).Error
		if err != nil {
			return nil, false, err
		}

		if user.ReferralMobile != "" {
			var refUser models.Users
			err = tx.Model(&models.Users{}).Where("mobile = ?", user.ReferralMobile).First(&refUser).Error
			if err != nil {
				return nil, false, err
			}

			var refWallet models.Wallet
			err = tx.Model(&models.Wallet{}).Where("user_id = ?", refUser.ID).First(&refWallet).Error
			if err != nil {
				return nil, false, err
			}

			refWallet.LockAmount -= 15000
			refWallet.Amount += 15000
			
			err := tx.Save(&refWallet).Error
			if err != nil {
				return nil,false,err
			}
		}

		return &dto.Alert{Message: "موجودی شما با موفقیت برداشت شد"}, true, nil
	}
	return &dto.Alert{Message: "مقدار موجودی شما کمتر از حد برداشت است"}, false, nil
}
