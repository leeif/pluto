package request

type MailRegister struct {
	Mail     string `json:"mail"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (mr *MailRegister) Validation() bool {
	if mr.Mail == "" || mr.Password == "" || mr.Name == "" {
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

func (ml *MailLogin) Validation() bool {
	if ml.Mail == "" || ml.Password == "" {
		return false
	}

	if ml.DeviceID == "" || ml.AppID == "" {
		return false
	}

	return true
}

type GoogleMobileLogin struct {
	IDToken  string `json:"id_token"`
	DeviceID string `json:"device_id"`
	AppID    string `json:"app_id"`
}

func (gml *GoogleMobileLogin) Validation() bool {
	if gml.IDToken == "" {
		return false
	}

	if gml.DeviceID == "" || gml.AppID == "" {
		return false
	}

	return true
}

type WechatMobileLogin struct {
	Code     string `json:"code"`
	DeviceID string `json:"device_id"`
	AppID    string `json:"app_id"`
}

func (wml *WechatMobileLogin) Validation() bool {
	if wml.Code == "" {
		return false
	}

	if wml.DeviceID == "" || wml.AppID == "" {
		return false
	}

	return true
}

type RegisterVerifyMail struct {
	Mail string `json:"mail"`
}

func (rvm *RegisterVerifyMail) Validation() bool {
	if rvm.Mail == "" {
		return false
	}

	return true
}

type ResetPasswordMail struct {
	Mail string `json:"mail"`
}

func (rpm *ResetPasswordMail) Validation() bool {
	if rpm.Mail == "" {
		return false
	}

	return true
}

type ResetPassword struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (rp *ResetPassword) Validation() bool {
	if rp.Token == "" || rp.Password == "" {
		return false
	}

	return true
}

type UpdateUserInfo struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Avatar string `json:"avatar"`
}

func (uui *UpdateUserInfo) Validation() bool {
	if uui.Name == "" && uui.Gender == "" && uui.Avatar == "" {
		return false
	}

	return true
}
