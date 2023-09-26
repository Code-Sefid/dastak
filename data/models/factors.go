package models

type FactorStatus string
type PaymentMethod string

const (
	PAID    FactorStatus = "paid"
	NOTPAID FactorStatus = "notPaid"
	CANCEL  FactorStatus = "cancel"
	CREATED FactorStatus = "created"
	PENDING FactorStatus = "pending"

	Offline PaymentMethod = "offline"
	Online  PaymentMethod = "online"
)

type Factors struct {
	BaseModel
	User       Users `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID     int
	Code       string
	OffPercent uint
	Status     FactorStatus
	FinalPrice float64
}

type FactorProducts struct {
	Product   Products `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductID int

	Factor   Factors `gorm:"foreignKey:FactorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FactorID int
}

type FactorPayment struct {
	Factor   Factors `gorm:"foreignKey:FactorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FactorID int

	FactorImage   string
	PaymentMethod PaymentMethod `gorm:"type:VARCHAR(255)"`
}

type FactorDetail struct {
	BaseModel
	Factor       Factors `gorm:"foreignKey:FactorID;"`
	FactorID     int
	FullName     string `gorm:"type:VARCHAR(255)"`
	Mobile       string `gorm:"type:VARCHAR(255)"`
	Province     string `gorm:"type:VARCHAR(255)"`
	City         string `gorm:"type:VARCHAR(255)"`
	Address      string
	PostalCode   string `gorm:"type:VARCHAR(255)"`
	TrackingCode string `gorm:"type:VARCHAR(255)"`
}
