package models

type Products struct {
	BaseModel
	User       Users `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID     int
	Title      string `gorm:"type:VARCHAR(255)"`
	Price      int
	Category   *Categories `gorm:"foreignKey:CategoryID;"`
	CategoryID int
	Inventory  int
}
