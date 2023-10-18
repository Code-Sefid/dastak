package models

type TransactionType string
type CheckOutType string

const (
	SALES     TransactionType = "sales"
	WITHDRAW  TransactionType = "withdraw"
	Referral  TransactionType = "referral"
	CORRECTED TransactionType = "corrected"
	DASTAK    TransactionType = "dastak"

	PENDINGCHECKOUT CheckOutType = "pending"
	REJECT          CheckOutType = "reject"
	DONE            CheckOutType = "done"
)

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
	CardName    string `gorm:"type:VARCHAR(255)"`
	CardNumber  string `gorm:"type:VARCHAR(255)"`
	ShabaNumber string `gorm:"type:VARCHAR(255)"`
}

type Transactions struct {
	BaseModel

	Description string

	FactorID int
	Factor   Factors `gorm:"foreignKey:FactorID"`

	User   Users `gorm:"foreignKey:UserID"`
	UserID int

	TransactionType TransactionType `gorm:"type:VARCHAR(255)"`

	Amount float64
}

type CheckOutRequest struct {
	BaseModel
	UserID int
	Amount int
	Status CheckOutType
}
