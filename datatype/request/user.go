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

type GoogleLogin struct {
	IDToken  string `json:"id_token"`
	DeviceID string `json:"device_id"`
	AppID    string `json:"app_id"`
}

func (gl *GoogleLogin) Validation() bool {
	if gl.IDToken == "" {
		return false
	}

	if gl.DeviceID == "" || gl.AppID == "" {
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
