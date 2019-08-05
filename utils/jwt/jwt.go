package jwt

import (
	"crypto"
	b64 "encoding/base64"
	"encoding/json"
	"time"

	"github.com/leeif/pluto/utils/rsa"
)

const (
	ALGRAS = "rsa"
)

type Head struct {
	Type string `json:"type"`
	Alg  string `json:"alg"`
}

type UserPayload struct {
	UserID   uint   `json:"userId"`
	DeviceID string `json:"deviceId"`
	AppID    string `json:"appId"`
	Expire   int64  `json:"expire"`
}

func GenerateUserJWT(head Head, payload UserPayload) (string, error) {
	head.Type = "JWT"
	h, err := json.Marshal(head)
	if err != nil {
		return "", err
	}
	// expire to one hour later
	payload.Expire = time.Now().Unix() + 60*60
	p, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	sig, err := rsa.SignWithPrivateKey([]byte(string(h)+string(p)), crypto.SHA256)

	if err != nil {
		return "", err
	}

	hB64 := b64.StdEncoding.EncodeToString(h)
	pB64 := b64.StdEncoding.EncodeToString(p)
	sigB64 := b64.StdEncoding.EncodeToString(sig)

	return hB64 + "." + pB64 + "." + sigB64, nil
}
