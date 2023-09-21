package dto

type CreateFactor struct {
	Products []FactorProduct `json:"products" binding:"required"`
	OffPercent uint `json:"off_percent" binding:"required"`
	FinalPrice float64 `json:"final_price" binding:"required"`
}

type FactorProduct struct {
	ID uint `json:"id"`
}

type FactorResponse struct{
	ID  int `json:"id"`
	Code string `json:"code"`
	OffPercent uint `json:"off_percent" binding:"required"`
	Status int
}

type UpdateFactor struct {
	OffPercent *uint `json:"off_percent" binding:"required"`
	FinalPrice *float64 `json:"final_price" binding:"required"`
}

type FactorProductResponse struct{
	ID  int `json:"id"`
	Code string `json:"code"`
	OffPercent uint `json:"off_percent" binding:"required"`
	Status int `json:"status"`
	Products []*ProductFactorResponse  `json:"products"`
}