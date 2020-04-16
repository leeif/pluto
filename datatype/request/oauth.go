package request

type OAuthAuthorize struct {
	ClientID     string `json:"client_id" schema:"client_id"`
	AppID        string `json:"app_id" schema:"app_id"`
	RedirectURI  string `json:"redirect_uri" schema:"redirect_uri"`
	ResponseType string `json:"response_type" schema:"response_type"`
	State        string `json:"state"`
	Scopes       string `json:"scopes"`
	LifeTime     int64  `json:"life_time"`
}

func (auth *OAuthAuthorize) Validation() bool {
	if auth.AppID == "" || auth.ResponseType == "" {
		return false
	}

	if auth.ClientID == "" {
		return false
	}
	return true
}

type OAuthTokens struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string
	ClientSecret string
	Scopes       string `json:"scopes"`
	AppID        string `json:"app_id"`
	DeviceID     string `json:"device_id"`
	Password     string `json:"password"`
	User         string `json:"user"`
	RefreshToken string `json:"refresh_token"`
}

func (token *OAuthTokens) Validation() bool {
	if token.GrantType == "" {
		return false
	}

	return true
}

func (token *OAuthTokens) ValidateAuthorizationCode() bool {
	if token.Code == "" || token.RedirectURI == "" {
		return false
	}
	return true
}

func (token *OAuthTokens) ValidatePasswordGrant() bool {
	if token.User == "" || token.Password == "" {
		return false
	}

	if token.AppID == "" {
		return false
	}

	return true
}

func (token *OAuthTokens) ValidateClientCredentials() bool {
	if token.AppID == "" {
		return false
	}
	return true
}

func (token *OAuthTokens) ValidateRefreshToken() bool {

	if token.AppID == "" || token.RefreshToken == "" {
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

type OAuthCreateClient struct {
	Key         string `json:"key"`
	Secret      string `json:"secret"`
	RedirectURI string `json:"redirect_uri"`
}

func (occ *OAuthCreateClient) Validation() bool {
	if occ.Key == "" || occ.RedirectURI == "" || occ.Secret == "" {
		return false
	}

	return true
}

type OAuthApproveClient struct {
	Key string `json:"key"`
}

func (occ *OAuthApproveClient) Validation() bool {
	if occ.Key == "" {
		return false
	}

	return true
}
