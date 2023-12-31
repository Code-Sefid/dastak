package dto

type CreateFactorDetailRequest struct {
	Code       string  `json:"code" binding:"required"`
	FullName   string  `json:"fullName" binding:"required"`
	Mobile     string  `json:"mobile" binding:"required,min=11,max=11,mobile"`
	Province   string  `json:"province" binding:"required,max=100"`
	City       string  `json:"city" binding:"required,max=100"`
	Address    string  `json:"address" binding:"required"`
	PostalCode *string `json:"postalCode" binding:"required"`
}

type UpdateFactorDetailRequest struct {
	FullName   string `json:"fullName"`
	Mobile     string `json:"mobile"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Address    string `json:"address" binding:"required"`
	PostalCode string `json:"postalCode"`
}

type FactorDetailResponse struct {
	ID         int    `json:"id"`
	FullName   string `json:"fullName"`
	Mobile     string `json:"mobile"`
	Address    string `json:"address"`
	Province   string `json:"province" binding:"required"`
	City       string `json:"city" binding:"required"`
	PostalCode string `json:"postalCode"`
	TrackingCode *string `json:"trackingCode"`
}

type FactorPayment struct {
	FactorID   int     `form:"factorId" binding:"required"`
	FinalPrice float64 `json:"finalPrice"`
	// Image         *multipart.FileHeader `form:"image" binding:"required"`
	// PaymentMethod string                `form:"paymentMethod" binding:"required"`
}


type AddTrackingCode struct {
	Code       string  `json:"code" binding:"required"`
	TrackingCode string `json:"trackingCode" binding:"required"`
}