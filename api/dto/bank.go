package dto

type CreateBankRequest struct {
	CardName    string `json:"cardName" binding:"required,min=3"`
	CardNumber  string `json:"cardNumber" binding:"required,min=12,max=12"`
	ShabaNumber string `json:"shabaNumber"`
}

type UpdateBankRequest struct {
	CardName    string `json:"cardName" binding:"min=3"`
	CardNumber  string `json:"cardNumber" binding:"min=12,max=12"`
	ShabaNumber string `json:"shabaNumber"`
}

type BankResponse struct {
	CardName   string `json:"cardName"`
	CardNumber string `json:"cardNumber"`
}
