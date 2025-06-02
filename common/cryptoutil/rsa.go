package cryptoutil

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// RSAKeyPair 保存RSA公私钥
type RSAKeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	PublicKeyString string
}

// LoadOrGenerateRSAKeys 从文件加载或生成新的RSA密钥对
func LoadOrGenerateRSAKeys(privateKeyPath, publicKeyPath string) (*RSAKeyPair, error) {
	// 尝试从文件加载
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err == nil {
		publicKey, pubKeyStr, err := loadPublicKey(publicKeyPath)
		if err == nil {
			return &RSAKeyPair{
				PrivateKey: privateKey,
				PublicKey:  publicKey,
				PublicKeyString: pubKeyStr,
			}, nil
		}
	}

	// 生成新的密钥对
	keyPair, err := generateRSAKeyPair(2048)
	if err != nil {
		return nil, err
	}

	// 保存到文件
	if err := savePrivateKey(privateKeyPath, keyPair.PrivateKey); err != nil {
		return nil, err
	}
	if err := savePublicKey(publicKeyPath, &keyPair.PrivateKey.PublicKey); err != nil {
		return nil, err
	}

	return keyPair, nil
}

func generateRSAKeyPair(bits int) (*RSAKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	publicKey := &privateKey.PublicKey
	pubKeyStr, err := publicKeyToPEM(publicKey)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		PublicKeyString: pubKeyStr,
	}, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func loadPublicKey(path string) (*rsa.PublicKey, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, "", errors.New("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, "", err
	}

	return pub.(*rsa.PublicKey), string(data), nil
}

func savePrivateKey(path string, key *rsa.PrivateKey) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	return pem.Encode(file, block)
}

func savePublicKey(path string, key *rsa.PublicKey) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	der := x509.MarshalPKCS1PublicKey(key)

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	}

	return pem.Encode(file, block)
}

func publicKeyToPEM(pubKey *rsa.PublicKey) (string, error) {
	der := x509.MarshalPKCS1PublicKey(pubKey)
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, block); err != nil {
		return "", err
	}

	return buf.String(), nil
}
