package request

import "time"

type MailRegister struct {
	Mail     string    `json:"mail"`
	Password string    `json:"password"`
	Gender   string    `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

func (mr MailRegister) Validation() bool {
	if mr.Mail == "" || mr.Password == "" {
		return false
	}
	return true
}

type MailLogin struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
	AppID    string `json:"app_id"`
}

func (ml MailLogin) Validation() bool {
	if ml.Mail == "" || ml.Password == "" {
		return false
	}

	if ml.DeviceID == "" || ml.AppID == "" {
		return false
	}

	return true
}
