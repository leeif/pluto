package jwt

import (
	"crypto"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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
	headB64 := b64.StdEncoding.EncodeToString(jwt.Head)
	payloadB64 := b64.StdEncoding.EncodeToString(jwt.Payload)
	signB64 := b64.StdEncoding.EncodeToString(jwt.Sign)
	return fmt.Sprintf("%s.%s.%s", headB64, payloadB64, signB64)
}

func (jwt *JWT) B64String() string {
	headB64 := b64.StdEncoding.EncodeToString(jwt.Head)
	payloadB64 := b64.StdEncoding.EncodeToString(jwt.Payload)
	signB64 := b64.StdEncoding.EncodeToString(jwt.Sign)
	plain := fmt.Sprintf("%s.%s.%s", headB64, payloadB64, signB64)
	return b64.StdEncoding.EncodeToString([]byte(plain))
}

type Head struct {
	Type string `json:"type"`
	Alg  string `json:"alg"`
}

type Payload struct {
	Create int64 `json:"create_time"`
	Expire int64 `json:"expire_time"`
}

type UserPayload struct {
	Payload
	UserID   uint   `json:"userId"`
	DeviceID string `json:"deviceId"`
	AppID    string `json:"appId"`
}

type RegisterVerifyPayload struct {
	Payload
	UserID uint `json:"userId"`
}

type PasswordResetPayload struct {
	Payload
	Mail string `json:"mail"`
}

type PasswordResetResultPayload struct {
	Payload
	Message string `json:"message"`
}

func setTimeField(payload interface{}, expire int64) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()
	v := reflect.ValueOf(payload)
	t := reflect.TypeOf(payload)
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	} else {
		return errors.New("Not pointer type")
	}
	create := v.FieldByName("Create")
	if !create.IsValid() {
		return errors.New("Create field is not valid")
	}
	create.SetInt(time.Now().Unix())
	exp := v.FieldByName("Expire")
	if !exp.IsValid() {
		return errors.New("Expire field is not valid")
	}
	exp.SetInt(time.Now().Unix() + expire)
	return nil
}

func GenerateJWT(head Head, payload interface{}, expire int64) (*JWT, *perror.PlutoError) {
	jwt := &JWT{}
	head.Alg = ALGRAS
	h, err := json.Marshal(head)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	jwt.Head = h

	if err := setTimeField(payload, expire); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}
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
	b, err := b64.StdEncoding.DecodeString(b64JWTToken)
	if err != nil {
		return nil, perror.InvalidJWTToekn
	}
	parts := strings.Split(string(b), ".")
	head, err := b64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, perror.InvalidJWTToekn
	}

	payload, err := b64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, perror.InvalidJWTToekn
	}

	signed, err := b64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, perror.InvalidJWTToekn
	}

	if err := rsa.VerifySignWithPublicKey(append(head, payload...), signed, crypto.SHA256); err != nil {
		return nil, perror.InvalidJWTToekn
	}
	return jwt, nil
}
