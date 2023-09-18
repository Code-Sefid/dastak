package models

import "time"

type JWTToken struct {
	ID                    int
	UserID                int
	AccessToken           string
	AccessTokenExpireTime time.Time
	RefreshToken          string
}
