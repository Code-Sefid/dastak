package dto

type CheckOut struct {
	Amount int `json:"amount" binding:"required"`
}
