package dto


type Payment struct {
	Code         string     `json:"code" binding:"required"`
	FinalPrice      string    `json:"finalPrice" binding:"required"`
}

type PaymentResponse struct {
	Message string `json:"message"`
	Result  int    `json:"result"`
	TrackID int    `json:"trackId"`
	Url    string `json:"url"`
}


type Verify struct {
	Success int    `json:"success" binding:"required"`
	Status  int    `json:"status" binding:"required"`
	TrackID string    `json:"trackId" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type VerifyResponse struct {
    Message         string  `json:"message"`
    Result          int     `json:"result"`
    RefNumber       *string `json:"refNumber"`
    PaidAt          string  `json:"paidAt"`
    Status          int     `json:"status"`
    Amount          int     `json:"amount"`
    OrderID         string  `json:"orderId"`
    Description     string  `json:"description"`
    CardNumber      *string `json:"cardNumber"`
}
