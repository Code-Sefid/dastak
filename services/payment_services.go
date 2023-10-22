package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/api/helper"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
	"gorm.io/gorm"
)

type PaymentService struct {
	logger   logging.Logger
	cfg      *config.Config
	database *gorm.DB
}

func NewPaymentService(cfg *config.Config) *PaymentService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &PaymentService{
		cfg:      cfg,
		database: database,
		logger:   logger,
	}
}

func (p *PaymentService) PaymentURL(ctx context.Context, req *dto.Payment) (*dto.PaymentResponse, *dto.Alert, error) {
	tx := p.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var factor models.Factors
	err := tx.Model(&models.Factors{}).Where("code = ?", req.Code).First(&factor).Error
	if err != nil {
		return nil, &dto.Alert{Message: "فاکتور شما وجود ندارن یا به مشکل خورده"}, err
	}

	var FactorProducts []*models.FactorProducts

	err = tx.Model(models.FactorProducts{}).Where("factor_id = ?", factor.ID).Preload("Product").Preload("Factor").Find(&FactorProducts).Error

	if err != nil {
		return nil, nil, err
	}

	var sum int
	for _, item := range FactorProducts {
		sum += item.Count * item.Product.Price
	}

	if factor.OffPercent != 0 {
		sum = sum / 100 * (100 - int(factor.OffPercent))
	}
	onePercent := float32(sum) / 100
	onePercent = (onePercent * 5)
	sum += int(onePercent) + factor.PostalCost

	merchant := p.cfg.Zibal.Token
	data := fmt.Sprintf(`{
		"merchant": "%s",
		"amount": %d,
		"callbackUrl": "%s",
		"orderId": "%s"
	}`, merchant, sum*10, p.cfg.Zibal.CallbackUrl+"/factor/"+req.Code, factor.Code)

	result, err := p.postToZibal("v1/request", data)
	if err != nil {
		return nil, &dto.Alert{Message: "خطایی در ارتباط با درگاه پرداخت رخ داده است"}, err
	}

	var paymentResponse dto.PaymentResponse
	if err := json.Unmarshal([]byte(result), &paymentResponse); err != nil {
		tx.Rollback()
		return nil, &dto.Alert{Message: "خطایی در ارتباط با درگاه پرداخت رخ داده است"}, err
	}

	if paymentResponse.Message != "success" {
		tx.Rollback()
		return nil, &dto.Alert{Message: "خطایی در ارتباط با درگاه پرداخت رخ داده است"}, err
	}
	trackIDStr := strconv.Itoa(paymentResponse.TrackID)

	paymentResponse.Url = fmt.Sprintf("https://gateway.zibal.ir/start/%s", trackIDStr)

	return &paymentResponse, nil, nil
}

func (p *PaymentService) CheckPayment(ctx context.Context, req *dto.Verify) (bool, *dto.Alert, error) {
	tx := p.database.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var factor models.Factors
	err := tx.Model(&models.Factors{}).Where("code = ?", req.Code).Preload("User").First(&factor).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &dto.Alert{Message: "فاکتور شما وجود ندارد یا به مشکل خورده"}, err
		}
		return false, &dto.Alert{Message: "خطا در پایگاه داده"}, err
	}

	var factorDetail models.FactorDetail
	err = tx.Model(&models.FactorDetail{}).Where("factor_id = ?", factor.ID).First(&factorDetail).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &dto.Alert{Message: "فاکتور شما وجود ندارد یا به مشکل خورده"}, err
		}
		return false, &dto.Alert{Message: "خطا در پایگاه داده"}, err
	}

	merchant := p.cfg.Zibal.Token
	data := `{
		"merchant" : "` + merchant + `",
		"trackId" : ` + req.TrackID + `
	}`

	result, err := p.postToZibal("v1/verify", data)
	if err != nil {
		tx.Rollback()
		return false, &dto.Alert{Message: "خطایی در ارتباط با درگاه پرداخت رخ داده است"}, err
	}

	var verifyResponse dto.VerifyResponse
	err = json.Unmarshal([]byte(result), &verifyResponse)
	if err != nil {
		tx.Rollback()
		return false, &dto.Alert{Message: "خطایی در پردازش پاسخ از درگاه پرداخت رخ داده است"}, err
	}

	if !p.verifyResult(verifyResponse.Result) {
		tx.Rollback()
		return false, &dto.Alert{Message: "پاسخ از درگاه پرداخت نامعتبر است"}, err
	}

	// if (verifyResponse.Result == 201 || verifyResponse.Result == 100) && (factor.Status == models.PENDING) {
		var FactorProducts []*models.FactorProducts

		err = tx.Model(models.FactorProducts{}).Where("factor_id = ?", factor.ID).Preload("Product").Preload("Factor").Find(&FactorProducts).Error

		if err != nil {
			return false, nil, err
		}

		var sum int

		for _, item := range FactorProducts {
			sum += item.Count * item.Product.Price
		}

		if factor.OffPercent != 0 {
			sum = sum / 100 * (100 - int(factor.OffPercent))
		}

		var onePercent float32
		var referralAmount float32
		var dastak float32

		// var referral int
		onePercent = float32(sum) / 100

		if factor.User.ReferralMobile != "" {
			referralAmount = (onePercent * 2)
			log.Print(referralAmount)

			dastak = (onePercent * 3)
			log.Print(dastak)

			transactionDastak := models.Transactions{
				FactorID:        factor.ID,
				Description:     fmt.Sprintf("درصد رفرال دستک از طرف %s", factorDetail.FullName),
				UserID:          factor.UserID,
				Amount:          float64(dastak),
				TransactionType: models.DASTAK,
			}

			err = tx.Create(&transactionDastak).Error
			if err != nil {
				tx.Rollback()
				return false, &dto.Alert{Message: "خطایی در ارتباط با درگاه پرداخت رخ داده است"}, err
			}


			transactionReferral := models.Transactions{
				FactorID: factor.ID,
				Description: fmt.Sprintf("درصد رفرال از طرف %s" , factorDetail.FullName),
				UserID: factor.UserID,
				Amount: float64(referralAmount),
				TransactionType: models.Referral,
			}

			err = tx.Create(&transactionReferral).Error
			if err != nil {
				tx.Rollback()
				return false, &dto.Alert{Message: "خطایی در ارتباط با درگاه پرداخت رخ داده است"}, err
			}

			var referralUser models.Users
			err = tx.Model(&models.Users{}).Where("mobile = ?", factor.User.ReferralMobile).First(&referralUser).Error
			if err != nil {
				tx.Rollback()
				return false, &dto.Alert{Message: "مشکلی در گرفتن یوزر رفال وجود دارد"}, err
			}

			var referralWallet models.Wallet
			err = tx.Model(&models.Wallet{}).Where("user_id = ?", referralUser.ID).First(&referralWallet).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return false, &dto.Alert{Message: "مشکلی در افزایش موجودی کیف پول داریم"}, err
			}else if errors.Is(err, gorm.ErrRecordNotFound){
				wallet := models.Wallet{
					UserId: referralUser.ID,
					Amount: int(referralAmount),
				}

				err := tx.Create(&wallet).Error
				if err != nil {
					tx.Rollback()
					return false, &dto.Alert{Message: "مشکل در ساخت کیف پول است"}, err
				}
			} else {
				newAmount := (referralWallet.Amount + int(referralAmount))
				err = tx.Model(&models.Wallet{}).Where("user_id = ?", referralUser.ID).Updates(map[string]interface{}{"amount": newAmount}).Error
				if err != nil {
					tx.Rollback()
					return false, &dto.Alert{Message: "مشکلی در افزایش موجودی کیف پول داریم"}, err
				}
			}
		} else {
			onePercent = (onePercent * 5)

			transactionDastak := models.Transactions{
				FactorID:        factor.ID,
				Description:     fmt.Sprintf("سود دستک %s", factorDetail.FullName),
				UserID:          factor.UserID,
				Amount:          float64(onePercent),
				TransactionType: models.DASTAK,
			}

			err = tx.Create(&transactionDastak).Error
			if err != nil {
				tx.Rollback()
				return false, &dto.Alert{Message: "خطایی در ارتباط با درگاه پرداخت رخ داده است"}, err
			}
		}

		transaction := models.Transactions{
			FactorID:        factor.ID,
			Description:     fmt.Sprintf("پرداخت فاکتور توسط مشتری %s انجام شد", factorDetail.FullName),
			UserID:          factor.UserID,
			Amount:          float64(sum),
			TransactionType: models.SALES,
		}

		err = tx.Create(&transaction).Error
		if err != nil {
			tx.Rollback()
			return false, &dto.Alert{Message: "خطایی در ارتباط با درگاه پرداخت رخ داده است"}, err
		}

		var wallet models.Wallet
		err = tx.Model(&models.Wallet{}).Where("user_id = ?", factor.UserID).First(&wallet).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return false, &dto.Alert{Message: "مشکلی در افزایش موجودی کیف پول داریم"}, err
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet := models.Wallet{
				UserId: factor.UserID,
				Amount: sum,
			}
			err := tx.Create(&wallet).Error
			if err != nil {
				tx.Rollback()
				return false, &dto.Alert{Message: "مشکل در ساخت کیف پول است"}, err
			}
		} else {
			newAmount := (wallet.Amount + sum)
			err = tx.Model(&models.Wallet{}).Where("user_id = ?", factor.UserID).Updates(map[string]interface{}{"amount": newAmount}).Error
			if err != nil {
				tx.Rollback()
				return false, &dto.Alert{Message: "مشکلی در افزایش موجودی کیف پول داریم"}, err
			}
		}

		factor.Status = models.PAID
		err = tx.Save(&factor).Error
		if err != nil {
			tx.Rollback()
			return false, &dto.Alert{Message: "مشکلی در تغییر وضعیت فاکتور پیش امده است"}, err
		}

		p.SendPayment(factorDetail.Mobile, factor.User.FullName, factorDetail.FullName, factor.Code)
		amountStr := helper.Separate(sum)
		p.SendPaymentToUser(factorDetail.FullName, amountStr, factor.Code, factorDetail.Mobile, factor.User.Mobile)
		return true, nil, nil
	// }
	return false, nil, nil
}

// Helper

func (p *PaymentService) postToZibal(path string, parameters string) (string, error) {
	var jsonStr = []byte(parameters)
	var url = "https://gateway.zibal.ir/" + path

	fmt.Println(bytes.NewBuffer(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (p *PaymentService) verifyResult(result int) bool {
	switch result {
	case 100:
		return true
	case 201:
		return true
	default:
		return false

	}
}

const payamakAPIURL = "http://api.payamak-panel.com/post/Send.asmx/SendByBaseNumber3"
const username = "09135882813"
const password = "T13Y7"

func (s *PaymentService) SendPayment(mobile, fullName, shopName, codeFactor string) error {
	values := url.Values{}
	values.Set("username", username)
	values.Set("password", password)
	values.Set("text", fmt.Sprintf("@167441@%s;%s;%s;", shopName, fullName, codeFactor))
	values.Set("to", mobile)

	url := payamakAPIURL + "?" + values.Encode()
	print(url)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	log.Print(response)

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP: %d", response.StatusCode)
	}

	return nil
}

func (s *PaymentService) SendPaymentToUser(fullName, amount, code, mobile, userMobile string) error {
	values := url.Values{}
	values.Set("username", username)
	values.Set("password", password)
	values.Set("text", fmt.Sprintf("@167730@%s;%s;%s;%s;", fullName, amount, code, mobile))
	values.Set("to", userMobile)

	url := payamakAPIURL + "?" + values.Encode()

	print(url)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	log.Print(response)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP: %d", response.StatusCode)
	}

	return nil
}

// End Helper
