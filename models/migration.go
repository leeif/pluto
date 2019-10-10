package models

type Migration struct {
	BaseModel
	Name string `gorm:"type:varchar(100);size:100;not null"`
}
