package models

type Device struct {
	BaseModel
	DeviceID string `gorm:"column:device_id;type:varchar(255);not null"`
	AppID    string `gorm:"column:app_id;type:varchar(255);not null"`
}
