package models

type RefreshToken struct {
	BaseModel
	RefreshToken string `gorm:"type:varchar(255);size:255;not null"`
	UserID       uint   `gorm:"column:user_id;not null"`
	DeviceAPPID  uint   `gorm:"column:device_app_id;not null"`
}
