package salt

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/scrypt"
)

const DEFAULTRANDLEN = 10

func EncodePassword(password string, salt string) (string, error) {
	dk, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}
	dst := base64.URLEncoding.EncodeToString(dk)
	return dst, nil
}

func RandomSalt(prefix ...string) string {
	salts := encryptRandSequence(DEFAULTRANDLEN, prefix)
	return salts
}

func DecodeSalt(salt string) (string, error) {
	dst, err := base64.URLEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}
	return string(dst), nil
}

func encryptRandSequence(n int, prefix []string) string {
	var src []byte
	b := make([]byte, n)
	rand.Read(b)
	if prefix != nil {
		src = append([]byte(strings.Join(prefix, ".")), b...)
	} else {
		src = b
	}
	dst := base64.URLEncoding.EncodeToString(src)
	return dst
}