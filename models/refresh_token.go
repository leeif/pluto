package models

type RefreshToken struct {
	BaseModel
	UserID       uint   `gorm:"column:user_id;not null;index:user_id_refresh_token"`
	RefreshToken string `gorm:"type:varchar(255);size:255;not null;index:user_id_refresh_token"`
	DeviceAPPID  uint   `gorm:"column:device_app_id;not null"`
}
