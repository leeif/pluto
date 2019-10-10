package models

type Salt struct {
	BaseModel
	UserID uint   `gorm:"column:user_id;not null"`
	Salt   string `gorm:"type:varchar(255);size:255;not null"`
}
