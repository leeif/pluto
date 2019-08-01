package request

type RefreshAccessToken struct {
	RefreshToken string `json:"refresh_token"`
	UseID        uint   `json:"user_id"`
	DeviceID     string `json:"device_id"`
	AppID        string `json:"app_id"`
}

func (rat RefreshAccessToken) Validation() bool {
	if rat.RefreshToken == "" || rat.UseID == 0 {
		return false
	}
	if rat.DeviceID == "" || rat.AppID == "" {
		return false
	}
	return true
}
