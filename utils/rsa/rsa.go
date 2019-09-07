package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/leeif/pluto/config"
)

var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
)

// Init : Init the rsa setting, generate new public private files or load from existing files
func Init(config *config.Config) error {
	privateKeyPath := path.Join(*config.RSA.Path, *config.RSA.Name)
	publicKeyPath := path.Join(*config.RSA.Path, *config.RSA.Name+".pub")

	_, privateKeyPathErr := os.Stat(privateKeyPath)

	_, publicKeyPathErr := os.Stat(publicKeyPath)

	if os.IsNotExist(privateKeyPathErr) && os.IsNotExist(publicKeyPathErr) {
		// generate private and public key
		privateKey, publicKey, err := generateKey()
		if err != nil {
			return err
		}
		if err := writePrivatekeyToFile(privateKey, privateKeyPath); err != nil {
			return err
		}

		if err := writePublicKeyToFile(publicKey, publicKeyPath); err != nil {
			return err
		}

		return nil
	}

	if !os.IsNotExist(privateKeyPathErr) && !os.IsNotExist(publicKeyPathErr) {
		// load private and public key from file
		var err error
		publicKey, err = loadPublicKeyFromFile(publicKeyPath)
		if err != nil {
			return err
		}

		privateKey, err = loadPrivateKeyFromFile(privateKeyPath)
		if err != nil {
			return err
		}

		return nil
	}

	// not pair only public key or private key exists
	return errors.New("Public and private key not in pair")
}

func loadPublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	keybuffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keybuffer)
	if block == nil {
		return nil, errors.New("Public key error")
	}

	pubkeyinterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publickey := pubkeyinterface.(*rsa.PublicKey)
	return publickey, nil
}

func loadPrivateKeyFromFile(filePath string) (*rsa.PrivateKey, error) {
	keybuffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(keybuffer))
	if block == nil {
		return nil, errors.New("Private key error")
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Parse private key error")
	}

	return privatekey, nil
}

func writePublicKeyToFile(publicKey *rsa.PublicKey, filePath string) error {
	keybytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keybytes,
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func writePrivatekeyToFile(privateKey *rsa.PrivateKey, filePath string) error {
	var keybytes = x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keybytes,
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func generateKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	publickey := &privatekey.PublicKey
	return privatekey, publickey, nil
}

// SignWithPrivtaeKey : Sign the src data with private key
func SignWithPrivateKey(src []byte, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	signed, err := rsa.SignPKCS1v15(rand.Reader, privateKey, hash, hashed)
	if err != nil {
		return []byte{}, err
	}

	return signed, nil
}

// VerifySignWithPublicKey : verify the signed data with public key
func VerifySignWithPublicKey(src, signed []byte, hash crypto.Hash) error {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	err := rsa.VerifyPKCS1v15(publicKey, hash, hashed, signed)
	if err != nil {
		return err
	}
	return nil
}

// GetPublicKey : Get the public key
func GetPublicKey() (string, error) {
	keybytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keybytes,
	}
	keybuffer := pem.EncodeToMemory(block)
	return string(keybuffer), nil
}
