package models

type UsersType string
type UsersAssist string

const (
	DIGITAL  UsersType = "digital"
	PHYSICAL UsersType = "physical"

	NORMAL  UsersAssist = "normal"
	PARTNER UsersAssist = "partner"
)

type Users struct {
	BaseModel
	FullName       string `gorm:"type:VARCHAR(255)"`
	Type           UsersType
	Assist         UsersAssist `gorm:"type:VARCHAR(255)"`
	Mobile         string      `gorm:"type:VARCHAR(255)"`
	Password       string      `gorm:"type:VARCHAR(255)"`
	ReferralMobile string      `gorm:"type:VARCHAR(255)"`
	SaleCount      uint
}
