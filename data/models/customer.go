package models

type Customer struct {
	BaseModel
	FullName     string `gorm:"type:VARCHAR(255)"`
	Mobile       string `gorm:"type:VARCHAR(255)"`
	Province     string `gorm:"type:VARCHAR(255)"`
	City         string `gorm:"type:VARCHAR(255)"`
	PostalCode   string `gorm:"type:VARCHAR(255)"`
	TrackingCode string `gorm:"type:VARCHAR(255)"`
}
