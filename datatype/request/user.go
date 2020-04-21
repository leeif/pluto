package request

const (
	defaultDeviceID = "UnKnown Device"
)

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

type PasswordLogin struct {
	Account  string `json:"account" schema:"account"`
	Password string `json:"password" schema:"password"`
	DeviceID string `json:"device_id" schema:"deviceID"`
	AppID    string `json:"app_id" schema:"appID"`
}

func (ml *PasswordLogin) Validation() bool {
	if ml.Account == "" || ml.Password == "" {
		return false
	}

	if ml.AppID == "" {
		return false
	}

	// set default deviceID
	if ml.DeviceID == "" {
		ml.DeviceID = defaultDeviceID
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

	if gml.AppID == "" {
		return false
	}

	// set default deviceID
	if gml.DeviceID == "" {
		gml.DeviceID = defaultDeviceID
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

	if wml.AppID == "" {
		return false
	}

	// set default deviceID
	if wml.DeviceID == "" {
		wml.DeviceID = defaultDeviceID
	}

	return true
}

type AppleMobileLogin struct {
	Code     string `json:"code"`
	DeviceID string `json:"device_id"`
	AppID    string `json:"app_id"`
}

func (aml *AppleMobileLogin) Validation() bool {
	if aml.Code == "" {
		return false
	}

	if aml.AppID == "" {
		return false
	}

	// set default deviceID
	if aml.DeviceID == "" {
		aml.DeviceID = defaultDeviceID
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

type ResetPasswordWeb struct {
	Password string `json:"password" schema:"password"`
}

func (rp *ResetPasswordWeb) Validation() bool {
	if rp.Password == "" {
		return false
	}

	return true
}

type UpdateUserInfo struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (uui *UpdateUserInfo) Validation() bool {
	if uui.Name == "" && uui.Avatar == "" {
		return false
	}

	return true
}

type VerifyAccessToken struct {
	Token string `json:"token"`
}

func (vat *VerifyAccessToken) Validation() bool {
	if vat.Token == "" {
		return false
	}

	return true
}

type VerifyIDToken struct {
	Token string `json:"token"`
}

func (vit *VerifyIDToken) Validation() bool {
	if vit.Token == "" {
		return false
	}

	return true
}

type MailBinding struct {
	Mail string `json:"mail"`
}

func (mb *MailBinding) Validation() bool {
	if mb.Mail == "" {
		return false
	}

	return true
}

type WechatBinding struct {
	Code string `json:"code"`
}

func (wb *WechatBinding) Validation() bool {
	if wb.Code == "" {
		return false
	}

	return true
}

type GoogleBinding struct {
	IDToken string `json:"id_token"`
}

func (gb *GoogleBinding) Validation() bool {
	if gb.IDToken == "" {
		return false
	}

	return true
}

type AppleBinding struct {
	Code string `json:"code"`
}

func (ab *AppleBinding) Validation() bool {
	if ab.Code == "" {
		return false
	}

	return true
}

type UnBinding struct {
	Type string `json:"type"`
}

func (ub *UnBinding) Validation() bool {
	if ub.Type == "" {
		return false
	}

	return true
}
