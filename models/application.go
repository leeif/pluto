package models

// Applictaion is which registered in the auth server
type Applictaion struct {
	BaseModel
	Name    string `gorm:"type:varchar(100);size:100;unique;not null"`
	Webhook string `gorm:"type:varchar(255);size:255;unique;not null"`
}
