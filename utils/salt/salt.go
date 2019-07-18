package salt

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/scrypt"
)

const DEFAULTRANDLEN = 10

func SaltEncode(userID string) string {
	id := []string{userID}
	return encryptRandSequence(DEFAULTRANDLEN, id...)
}

func SaltDecode(salt string) (string, error) {
	dst, err := base64.URLEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}
	return string(dst), nil
}

func PsdHandler(psd string, salts []byte) (string, error) {
	dk, err := scrypt.Key([]byte(psd), salts, 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}
	dst := base64.URLEncoding.EncodeToString(dk)
	return dst, nil
}

func encryptRandSequence(n int, userID ...string) string {
	var src []byte
	b := make([]byte, n)
	rand.Read(b)
	if userID != nil {
		src = append([]byte(userID[0]+":"), b...)
	} else {
		src = b
	}
	dst := base64.URLEncoding.EncodeToString(src)
	return dst
}
