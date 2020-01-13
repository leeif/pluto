package request

type OAuth2Authorize struct {
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	ResponseType string `json:"response_type"`
	State        string `json:"state"`
	Scope        string `json:"scope"`
}

func (auth *OAuth2Authorize) Validation() bool {
	if auth.ClientID == "" || auth.ResponseType == "" {
		return false
	}
	return true
}

type OAuth2Tokens struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
	ClientID    string `json:"client_id"`
}

func (token *OAuth2Tokens) Validation() bool {
	if token.Code == "" || token.GrantType == "" || token.RedirectURI == "" || token.ClientID == "" {
		return false
	}
	return true
}

type RefreshAccessToken struct {
	RefreshToken string `json:"refresh_token"`
	UseID        uint   `json:"user_id"`
	DeviceID     string `json:"device_id"`
	AppID        string `json:"app_id"`
}

func (rat *RefreshAccessToken) Validation() bool {
	if rat.RefreshToken == "" || rat.UseID == 0 {
		return false
	}
	if rat.DeviceID == "" || rat.AppID == "" {
		return false
	}
	return true
}
