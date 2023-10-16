package models

type UsersType string

const (
	DIGITAL  UsersType = "digital"
	PHYSICAL UsersType = "physical"
)

type Users struct {
	BaseModel
	FullName  string `gorm:"type:VARCHAR(255)"`
	Type      UsersType
	Mobile    string `gorm:"type:VARCHAR(255)"`
	Password  string `gorm:"type:VARCHAR(255)"`
	ReferralMobile   string `gorm:"type:VARCHAR(255)"`
	SaleCount uint
}
