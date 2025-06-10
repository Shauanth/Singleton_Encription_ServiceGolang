package crypton

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"io"
)

type Config struct {
	EncryptionKey string `json:"encryption_key"`
	Salt          string `json:"salt"`
}

// Nota: No hay ningun beneficio practico al reusar el mismo salt en cada ciphertext, mejorar esquema a futuro.
func loadKey(configuracion Config) ([]byte, error) {
	if configuracion.EncryptionKey == "" {
		return nil, errors.New("[crypto] key is empty, have you set it apppropiately?")
	}
	if configuracion.Salt == "" {
		return nil, errors.New("[crypto] salt is empty, have you set it appropiately?")
	}
	key, _ := pbkdf2.Key(sha512.New, configuracion.EncryptionKey, []byte(configuracion.Salt), 4096, 32)
	return key, nil
}

func Encrypt(plaintext string, configuracion Config) (string, error) {
	plaintextBytes := []byte(plaintext)
	key, err := loadKey(configuracion)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := aesgcm.Seal(nonce, nonce, plaintextBytes, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(ciphertext string, configuracion Config) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	key, err := loadKey(configuracion)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return "", errors.New("[Crypto] Ciphertext is smaller than nonce, is the ciphertext correct?")
	}
	nonce, ciphertextBytes := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]
	plaintextBytes, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}
	return string(plaintextBytes), nil
}
