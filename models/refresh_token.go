package models

type RefreshToken struct {
	BaseModel
	RefreshToken string `gorm:"type:varchar(255);size:255"`
	UserID       uint   `gorm:"column:user_id"`
	DeviceID     string `gorm:"column:device_id;type:varchar(255);not null"`
	AppID        string `gorm:"column:app_id;type:varchar(255);not null"`
}
