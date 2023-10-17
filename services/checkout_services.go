package services

import (
	"context"
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


func (c *CheckOutService) CheckOutMony(ctx context.Context, userID int, req dto.CheckOut)(*dto.Alert,bool,error){
	tx := c.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}else{
			tx.Commit()
		}
	}()
	var user models.Users
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, false , err
	}

	var wallet models.Wallet
	err := tx.Model(&models.Wallet{}).Where("user_id = ?", user.ID).First(&wallet).Error
	if err != nil {
		tx.Rollback()
		return &dto.Alert{Message: "کیف پول شما خالی است"},false ,nil
	}
	
	if req.Amount <= 50000 && wallet.Amount >= req.Amount{
		amount := wallet.Amount - req.Amount
		err = tx.Model(&models.Wallet{}).Where("user_id = ?", user.ID).Update("amount" , amount).Error
		if err != nil {
			return nil, false,err
		}
	
		 checkOut := models.CheckOutRequest{
			UserID: user.ID,
			Amount: req.Amount,
			Status: models.PENDINGCHECKOUT,
		}
	
		err = tx.Model(&models.CheckOutRequest{}).Create(&checkOut).Error
		if err != nil {
			return nil, false,err
		}
	
		transactions := models.Transactions{
			UserID: user.ID,
			Description: "درخواست برداشت شما به مبلغ " + strconv.Itoa(req.Amount) + " تومان در حال بررسی است",
			Amount: float64(req.Amount),
			FactorID: 1,
			TransactionType: models.WITHDRAW,
		}
		
		err = tx.Model(&models.Transactions{}).Create(&transactions).Error
		if err != nil {
			return nil,false, err 
		}
		
	
		return &dto.Alert{Message: "موجودی شما با موفقیت برداشت شد"},true,nil
	}
	return &dto.Alert{Message: "مقدار موجودی شما کمتر از حد برداشت است"},false,nil
}