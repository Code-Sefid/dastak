package dto

type CreateFactor struct {
	Products    []FactorItem `json:"products" binding:"required"`
	OffPercent  uint         `json:"offPercent"`
	Description string       `json:"description"`
	PostalCost  int          `json:"postalCost"`
}

type FactorResponse struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	OffPercent  uint   `json:"offPercent"`
	Status      int    `json:"status"`
	PostalCost  int    `json:"postalCost"`
	Description string `json:"description"`
}

type UpdateFactor struct {
	OffPercent *uint `json:"offPercent" binding:"required"`
}

type FactorProductResponse struct {
	ID          int                      `json:"id"`
	Code        string                   `json:"code"`
	OffPercent  uint                     `json:"offPercent" binding:"required"`
	Status      int                      `json:"status"`
	Description string                   `json:"description"`
	PostalCost  int                      `json:"postalCost"`
	Account     *AccountResponse         `json:"account"`
	Products    []*ProductFactorResponse `json:"products"`
	Factor      *FactorDetailResponse    `json:"detail"`
}

type AccountResponse struct {
	Name string `json:"name"`
}
type FactorItem struct {
	ID    uint `json:"id"`
	Count int  `json:"count"`
}
