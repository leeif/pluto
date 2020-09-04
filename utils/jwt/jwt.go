package jwt

import (
	"crypto"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	perror "github.com/MuShare/pluto/datatype/pluto_error"

	"github.com/MuShare/pluto/utils/rsa"
)

const (
	RS256ALGRAS    = "RS256"
	ACCESS         = "access"
	REGISTERVERIFY = "register_verify"
	PASSWORDRESET  = "password_reset"
)

type JWT struct {
	Head    string
	Payload string
	Sign    string
}

func (jwt *JWT) String() string {
	return fmt.Sprintf("%s.%s.%s", jwt.Head, jwt.Payload, jwt.Sign)
}

func (jwt *JWT) B64String() string {
	plain := fmt.Sprintf("%s.%s.%s", jwt.Head, jwt.Payload, jwt.Sign)
	return b64.RawURLEncoding.EncodeToString([]byte(plain))
}

func (jwt *JWT) UnmarshalPayload(v interface{}) *perror.PlutoError {
	b, err := b64.RawURLEncoding.DecodeString(jwt.Payload)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	return nil
}

type Head struct {
	Type string `json:"type"`
	Alg  string `json:"alg"`
}

type Payload struct {
	Type   string `json:"type"`
	Create int64  `json:"iat"`
	Expire int64  `json:"exp"`
}

type AccessPayload struct {
	Payload
	UserID uint     `json:"sub"`
	AppID  string   `json:"iss"`
	Scopes []string `json:"scopes"`
}

func NewAccessPayload(userID uint, scopes string, appID string, expire int64) *AccessPayload {
	up := &AccessPayload{}
	up.UserID = userID
	up.AppID = appID
	up.Scopes = strings.Split(scopes, ",")

	up.Payload.Type = ACCESS
	up.Payload.Create = time.Now().Unix()
	up.Payload.Expire = time.Now().Unix() + expire
	return up
}

type IDTokenPayload struct {
	Payload
	UserID uint   `json:"sub"`
	AppID  string `json:"iss"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	Mail   string `json:"mail"`
}

func NewIDTokenPayload(userID uint, role, name, mail, appID string, expire int64) *IDTokenPayload {
	idToken := &IDTokenPayload{}
	idToken.UserID = userID
	idToken.AppID = appID
	idToken.Role = role
	idToken.Name = name
	idToken.Mail = mail

	idToken.Payload.Type = ACCESS
	idToken.Payload.Create = time.Now().Unix()
	idToken.Payload.Expire = time.Now().Unix() + expire
	return idToken
}

type RegisterVerifyPayload struct {
	Payload
	UserID uint `json:"sub"`
}

func NewRegisterVerifyPayload(userID uint, expire int64) *RegisterVerifyPayload {
	rvp := &RegisterVerifyPayload{}
	rvp.UserID = userID

	rvp.Payload.Type = REGISTERVERIFY
	rvp.Payload.Create = time.Now().Unix()
	rvp.Payload.Expire = time.Now().Unix() + expire
	return rvp
}

type PasswordResetPayload struct {
	Payload
	Mail string `json:"mail"`
}

func NewPasswordResetPayload(mail string, expire int64) *PasswordResetPayload {
	rrp := &PasswordResetPayload{}
	rrp.Mail = mail

	rrp.Payload.Type = PASSWORDRESET
	rrp.Payload.Create = time.Now().Unix()
	rrp.Payload.Expire = time.Now().Unix() + expire
	return rrp
}

type PasswordResetResultPayload struct {
	Payload
	Successed bool `json:"successed"`
}

func GenerateRSA256JWT(payload interface{}) (*JWT, *perror.PlutoError) {
	jwt := &JWT{}
	head := Head{}
	head.Alg = RS256ALGRAS
	head.Type = "jwt"
	h, err := json.Marshal(head)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	jwt.Head = b64.RawURLEncoding.EncodeToString(h)

	p, err := json.Marshal(payload)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	jwt.Payload = b64.RawURLEncoding.EncodeToString(p)

	sig, err := rsa.SignWithPrivateKey([]byte(fmt.Sprintf("%s.%s", jwt.Head, jwt.Payload)), crypto.SHA256)

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	jwt.Sign = b64.RawURLEncoding.EncodeToString(sig)

	return jwt, nil
}

func VerifyRS256JWT(token string) (*JWT, *perror.PlutoError) {
	jwt := &JWT{}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, perror.InvalidJWTToken
	}

	jwt.Head = parts[0]
	jwt.Payload = parts[1]

	decodedSign, err := b64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, perror.InvalidJWTToken
	}

	jwt.Sign = parts[2]

	concat := fmt.Sprintf("%s.%s", parts[0], parts[1])

	if err := rsa.VerifySignWithPublicKey([]byte(concat), decodedSign, crypto.SHA256); err != nil {
		return nil, perror.InvalidJWTToken
	}
	return jwt, nil
}

func VerifyB64RS256JWT(b64JWTToken string) (*JWT, *perror.PlutoError) {
	b, err := b64.RawURLEncoding.DecodeString(b64JWTToken)
	if err != nil {
		return nil, perror.InvalidJWTToken
	}
	jwt, perr := VerifyRS256JWT(string(b))
	if perr != nil {
		return nil, perr
	}

	return jwt, nil
}
