package request

import (
	"github.com/MuShare/pluto/utils/general"
	"strings"
	"unicode"
)

const (
	defaultDeviceID = "UnKnown Device"
)

type MailRegister struct {
	Mail     string `json:"mail"`
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	AppName  string `json:"app_id"`
}

func (mr *MailRegister) Validation() bool {
	if !general.ValidMail(mr.Mail) || strings.TrimSpace(mr.Password) == "" || strings.TrimSpace(mr.Name) == "" {
		return false
	}
	return validateUserID(mr.UserID)
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
	Mail    string `json:"mail"`
	AppName string `json:"app_id"`
	UserID  string `json:"user_id"`
}

func (rvm *RegisterVerifyMail) Validation() bool {
	if rvm.Mail == "" && rvm.UserID == "" {
		return false
	}

	return true
}

type ResetPasswordMail struct {
	Mail    string `json:"mail"`
	AppName string `json:"app_id"`
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
	UserID string `json:"user_id"`
}

func (uui *UpdateUserInfo) Validation() bool {
	if strings.TrimSpace(uui.Name) == "" && uui.Avatar == "" && uui.UserID == "" {
		return false
	}
	return validateUserID(uui.UserID)
}

func validateUserID(userID string) bool {
	if len(userID) != 0 && strings.TrimSpace(userID) == "" || len(userID) > 100 {
		return false
	}
	//can only be consisted by digit, letter, '-' or '_'
	for _, r := range userID {
		if r == '_' || r == '-' {
			return true
		}
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			return false
		}
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

type Binding struct {
	Mail    string `json:"mail"`
	Code    string `json:"code"`
	IDToken string `json:"id_token"`
	Type    string `json:"type"`
}

func (binding *Binding) Validation() bool {
	if binding.Type == "" {
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

type PublicUserInfos struct {
	IDs []string `schema:"ids"`
}
