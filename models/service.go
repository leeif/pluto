package models

// Service : Services which registered in the auth server
type Service struct {
	BaseModel
	Name    string `gorm:"type:varchar(100);size:100;unique;not null"`
	Webhook string `gorm:"type:varchar(255);size:255;unique;not null"`
}
