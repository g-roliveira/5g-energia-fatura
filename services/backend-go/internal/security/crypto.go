package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type Cipher struct {
	key []byte
}

func NewCipher(secret string) (*Cipher, error) {
	if secret == "" {
		return nil, fmt.Errorf("BACKEND_ENCRYPTION_KEY não configurada")
	}
	sum := sha256.Sum256([]byte(secret))
	return &Cipher{key: sum[:]}, nil
}

func (c *Cipher) Encrypt(plaintext string) (ciphertext string, nonce string, err error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}
	nonceBytes := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonceBytes); err != nil {
		return "", "", err
	}
	encrypted := gcm.Seal(nil, nonceBytes, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(encrypted), base64.StdEncoding.EncodeToString(nonceBytes), nil
}

func (c *Cipher) Decrypt(ciphertext string, nonce string) (string, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return "", err
	}
	plain, err := gcm.Open(nil, nonceBytes, cipherBytes, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
