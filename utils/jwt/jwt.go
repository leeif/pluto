package jwt

import (
	"crypto"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	perror "github.com/leeif/pluto/datatype/pluto_error"

	"github.com/leeif/pluto/utils/rsa"
)

const (
	ALGRAS              = "rsa"
	ACCESS              = "access"
	REGISTERVERIFY      = "register_verify"
	PASSWORDRESET       = "password_reset"
	PASSWORDRESETRESULT = "password_reset_result"
)

type JWT struct {
	Head    []byte
	Payload []byte
	Sign    []byte
}

func (jwt *JWT) String() string {
	headB64 := b64.RawStdEncoding.EncodeToString(jwt.Head)
	payloadB64 := b64.RawStdEncoding.EncodeToString(jwt.Payload)
	signB64 := b64.RawStdEncoding.EncodeToString(jwt.Sign)
	return fmt.Sprintf("%s.%s.%s", headB64, payloadB64, signB64)
}

func (jwt *JWT) B64String() string {
	headB64 := b64.RawStdEncoding.EncodeToString(jwt.Head)
	payloadB64 := b64.RawStdEncoding.EncodeToString(jwt.Payload)
	signB64 := b64.RawStdEncoding.EncodeToString(jwt.Sign)
	plain := fmt.Sprintf("%s.%s.%s", headB64, payloadB64, signB64)
	return b64.RawStdEncoding.EncodeToString([]byte(plain))
}

type Head struct {
	Type string `json:"type"`
	Alg  string `json:"alg"`
}

type Payload struct {
	Type   string `json:"type"`
	Create int64  `json:"create_time"`
	Expire int64  `json:"expire_time"`
}

type UserPayload struct {
	Payload
	UserID    uint   `json:"userId"`
	DeviceID  string `json:"deviceId"`
	AppID     string `json:"appId"`
	LoginType string `json:"login_type"`
}

func NewUserPayload(userID uint, deviceID, appID, loginType string, expire int64) *UserPayload {
	up := &UserPayload{}
	up.UserID = userID
	up.DeviceID = deviceID
	up.AppID = appID
	up.LoginType = loginType

	up.Payload.Type = ACCESS
	up.Payload.Create = time.Now().Unix()
	up.Payload.Expire = time.Now().Unix() + expire
	return up
}

type RegisterVerifyPayload struct {
	Payload
	UserID uint `json:"userId"`
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

func NewPasswordResetResultPayload(successed bool, expire int64) *PasswordResetResultPayload {
	rrrp := &PasswordResetResultPayload{}
	rrrp.Successed = successed

	rrrp.Payload.Type = PASSWORDRESETRESULT
	rrrp.Payload.Create = time.Now().Unix()
	rrrp.Payload.Expire = time.Now().Unix() + expire
	return rrrp
}

func GenerateRSAJWT(payload interface{}) (*JWT, *perror.PlutoError) {
	jwt := &JWT{}
	head := Head{}
	head.Alg = ALGRAS
	head.Type = "jwt"
	h, err := json.Marshal(head)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	jwt.Head = h

	p, err := json.Marshal(payload)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	jwt.Payload = p

	sig, err := rsa.SignWithPrivateKey([]byte(string(h)+string(p)), crypto.SHA256)

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	jwt.Sign = sig

	return jwt, nil
}

func VerifyB64JWT(b64JWTToken string) (*JWT, *perror.PlutoError) {
	jwt := &JWT{}
	b, err := b64.RawStdEncoding.DecodeString(b64JWTToken)
	if err != nil {
		return nil, perror.InvalidJWTToekn
	}
	parts := strings.Split(string(b), ".")
	head, err := b64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, perror.InvalidJWTToekn
	}

	jwt.Head = head

	payload, err := b64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, perror.InvalidJWTToekn
	}

	jwt.Payload = payload

	sign, err := b64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, perror.InvalidJWTToekn
	}

	jwt.Sign = sign

	if err := rsa.VerifySignWithPublicKey(append(head, payload...), sign, crypto.SHA256); err != nil {
		return nil, perror.InvalidJWTToekn
	}
	return jwt, nil
}
