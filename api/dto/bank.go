package dto


type CreateBankRequest struct {
	CardName string `json:"cardName" binding:"required"`
	CardNumber string `json:"cardNumber" binding:"required"`
	ShabaNumber string `json:"shabaNumber"`
}

type UpdateBankRequest struct {
	CardName string `json:"cardName" `
	CardNumber string `json:"cardNumber"`
	ShabaNumber string `json:"shabaNumber"`
}

type BankResponse struct {
	CardName string `json:"cardName"`
	CardNumber string `json:"cardNumber"`
}