package models

type Wallet struct {
	BaseModel
	User   Users `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserId int
	Amount int
}

type BankAccounts struct {
	BaseModel
	User        Users `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID      int
	CardName       string `gorm:"type:VARCHAR(255)"`
	CardNumber  string `gorm:"type:VARCHAR(255)"`
	ShabaNumber string `gorm:"type:VARCHAR(255)"`
}

type Transactions struct {
	BaseModel
	User        Users `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserId      int
	Credit      int
	Debit       int
	Description string
}
