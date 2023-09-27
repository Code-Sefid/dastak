package dto

type CreateFactor struct {
	Products   []int `json:"products" binding:"required"`
	OffPercent uint            `json:"offPercent"`
}

type FactorResponse struct {
	ID         int    `json:"id"`
	Code       string `json:"code"`
	OffPercent uint   `json:"offPercent" binding:"required"`
	Status     int
}

type UpdateFactor struct {
	OffPercent *uint    `json:"offPercent" binding:"required"`
}

type FactorProductResponse struct {
	ID         int                      `json:"id"`
	Code       string                   `json:"code"`
	OffPercent uint                     `json:"offPercent" binding:"required"`
	Status     int                      `json:"status"`
	Products   []*ProductFactorResponse `json:"products"`
}

type FactorItem struct {
	ID uint `json:"id"`
}
