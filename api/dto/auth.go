package dto

type Register struct {
	FullName    string `json:"full_name" binding:"required,min=3"`
	AccountType int    `json:"account_type" binding:"required"`
	SaleCount   uint   `json:"sale_count" binding:"required"`
	Mobile      string `json:"mobile" binding:"required,min=11,max=11"`
}

type Login struct {
	Mobile   string `json:"mobile" binding:"required,min=11,max=11"`
	Password string `json:"password" binding:"required,min=6"`
}

type Mobile struct {
	Mobile string `json:"mobile" binding:"required,min=11,max=11"`
}
