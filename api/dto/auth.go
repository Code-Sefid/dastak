package dto

// Register struct
type Register struct {
	FullName    string  `json:"full_name" binding:"required,min=3" `
	AccountType int     `json:"account_type" binding:"required" `
	SaleCount   uint    `json:"sale_count" binding:"required" `
	Mobile      string  `json:"mobile" binding:"required,min=11,max=11,mobile"`
	Referral    *string `json:"referral" binding:"omitempty,min=11,max=11,mobile"`
}

// Login struct
type Login struct {
	Mobile   string `json:"mobile" binding:"required,min=11,max=11,mobile" `
	Password string `json:"password" binding:"required,min=4" `
}

// Mobile struct
type Mobile struct {
	Mobile string `json:"mobile" binding:"required,min=11,max=11,mobile" `
}
