package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/soheilkhaledabdi/dastak/api/dto"
	"github.com/soheilkhaledabdi/dastak/config"
	"github.com/soheilkhaledabdi/dastak/constants"
	"github.com/soheilkhaledabdi/dastak/data/db"
	"github.com/soheilkhaledabdi/dastak/data/models"
	"github.com/soheilkhaledabdi/dastak/pkg/logging"
	"github.com/soheilkhaledabdi/dastak/pkg/service_errors"
	"gorm.io/gorm"
)

type TokenService struct {
	logger logging.Logger
	cfg    *config.Config
}

type tokenDto struct {
	UserId int
	Mobile string
}

func NewTokenService(cfg *config.Config) *TokenService {
	logger := logging.NewLogger(cfg)
	return &TokenService{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *TokenService) GenerateToken(token *tokenDto) (*dto.TokenDetail, error) {
	var err error
	td := &dto.TokenDetail{}
	td.AccessTokenExpireTime = time.Now().Add(s.cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix()

	atc := jwt.MapClaims{}

	atc[constants.UserIdKey] = token.UserId
	atc[constants.ExpireTimeKey] = td.AccessTokenExpireTime

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)

	td.AccessToken, err = at.SignedString([]byte(s.cfg.JWT.Secret))

	if err != nil {
		return nil, err
	}

	rtc := jwt.MapClaims{}

	rtc[constants.UserIdKey] = token.UserId

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtc)

	td.RefreshToken, err = rt.SignedString([]byte(s.cfg.JWT.RefreshSecret))

	if err != nil {
		return nil, err
	}

	db := db.GetDb()

	existingToken := models.JWTToken{}
	if err := db.Where("user_id = ?", token.UserId).First(&existingToken).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		tokenModel := models.JWTToken{
			UserID:                token.UserId,
			AccessToken:           td.AccessToken,
			AccessTokenExpireTime: time.Unix(td.AccessTokenExpireTime, 0),
			RefreshToken:          td.RefreshToken,
		}
		if err := db.Create(&tokenModel).Error; err != nil {
			return nil, err
		}
	} else {
		existingToken.AccessToken = td.AccessToken
		existingToken.AccessTokenExpireTime = time.Unix(td.AccessTokenExpireTime, 0)
		existingToken.RefreshToken = td.RefreshToken
		if err := db.Save(&existingToken).Error; err != nil {
			return nil, err
		}
	}

	return td, nil
}

func (s *TokenService) VerifyToken(token string) (*jwt.Token, error) {

	db := db.GetDb()
	var tokenModel models.JWTToken
	if err := db.Where("access_token = ? AND user_id IS NOT NULL", token).First(&tokenModel).Error; err != nil {
		return nil, err
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
		}
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken, nil
}

func (s *TokenService) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}

	verifyToken, err := s.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, &service_errors.ServiceError{EndUserMessage: service_errors.ClaimsNotFound}
}

func (s *TokenService) DeleteToken(token string) error {
	db := db.GetDb()

	// Find the token in the database
	var tokenModel models.JWTToken
	if err := db.Where("access_token = ?", token).First(&tokenModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &service_errors.ServiceError{EndUserMessage: service_errors.TokenNotFound}
		}
		return err
	}

	// Delete the token from the database
	if err := db.Delete(&tokenModel).Error; err != nil {
		return err
	}

	return nil
}
