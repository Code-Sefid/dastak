package models

type FactorStatus string
type PaymentMethod string

const (
	PAID    FactorStatus = "paid"
	NOTPAID FactorStatus = "notPaid"
	CANCEL  FactorStatus = "cancel"
	CREATED FactorStatus = "created"
	PENDING FactorStatus = "pending"

	CARDBYCARD PaymentMethod = "cardByCard"
	DEPOSIT    PaymentMethod = "deposit"
)

type Factors struct {
	BaseModel
	User   Users `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID int
	Code          string
	OffPercent    uint
	Status        FactorStatus
	FinalPrice    float64
	PaymentMethod PaymentMethod
}

type FactorProducts struct {
	Product   Products `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductID int

	Factor   Factors `gorm:"foreignKey:FactorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FactorID int
}

type FactorPayment struct {
	User   Users `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID int

	Factor   Factors `gorm:"foreignKey:FactorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FactorID int
}
