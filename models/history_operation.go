package models

type HistoryOperation struct {
	BaseModel
	UserID        uint   `gorm:"column:user_id;not null"`
	OperationType string `gorm:"column:type;type:varchar(20);not null"`
}
