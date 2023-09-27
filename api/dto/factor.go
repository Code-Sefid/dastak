package dto

type CreateFactor struct {
	Products   []int `json:"products" binding:"required"`
	OffPercent uint            `json:"offPercent"`
}

type FactorResponse struct {
	ID         int    `json:"id"`
	Code       string `json:"code"`
	OffPercent uint   `json:"offPercent"`
	Status     int     `json:"status"`
}

type UpdateFactor struct {
	OffPercent *uint    `json:"offPercent" binding:"required"`
}

type FactorProductResponse struct {
	ID         int                      `json:"id"`
	Code       string                   `json:"code"`
	OffPercent uint                     `json:"offPercent" binding:"required"`
	Status     int                      `json:"status"`
	Account    *AccountResponse         `json:"account"`
	Products   []*ProductFactorResponse `json:"products"`
}


type AccountResponse struct{
	Name string `json:"name"`
}
type FactorItem struct {
	ID uint `json:"id"`
}
