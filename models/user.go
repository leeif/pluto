package models

import (
	"encoding/json"
	"time"
)

type User struct {
	BaseModel
	Mail          *string    `gorm:"type:varchar(255);size:60" json:"mail"`
	Name          *string    `gorm:"type:varchar(60);size:60;not null" json:"name"`
	Role          string     `gorm:"type:varchar(60);size:60" json:"-"`
	Gender        *string    `gorm:"type:varchar(10);size:10" json:"gender"`
	Password      *string    `gorm:"type:varchar(255);not null" json:"-"`
	Birthday      *time.Time `json:"birthday"`
	Avatar        string     `gorm:"type:varchar(255)" json:"avatar"`
	Verified      bool       `json:"-"`
	LoginType     string     `gorm:"type:varchar(10);not null;UNIQUE_INDEX:login" json:"login_type"`
	IdentifyToken string     `gorm:"type:varchar(255);not null;UNIQUE_INDEX:login" json:"-"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	type user struct {
		ID        uint    `json:"id"`
		CreatedAt int64   `json:"create_at"`
		UpdatedAt int64   `json:"updated_at"`
		DeletedAt int64   `json:"delete_at"`
		Mail      *string `json:"mail"`
		Name      *string `json:"name"`
		Gender    *string `json:"gender"`
		Birthday  int64   `json:"birthday"`
		Avatar    string  `json:"avatar"`
		LoginType string  `json:"login_type"`
	}
	s := &user{
		ID:        u.ID,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
		Mail:      u.Mail,
		Name:      u.Name,
		Gender:    u.Gender,
		Avatar:    u.Avatar,
		LoginType: u.LoginType,
	}
	if u.Birthday != nil {
		s.Birthday = u.Birthday.Unix()
	}
	if u.DeletedAt != nil {
		s.DeletedAt = u.DeletedAt.Unix()
	}
	return json.Marshal(s)
}
