package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Config struct {
	EncryptionKey string `json:"encryption_key"`
	Salt          string `json:"salt"`
}

// Nota: No hay ningun beneficio practico al reusar el mismo salt en cada ciphertext, mejorar esquema a futuro.
func loadKey() ([]byte, error) {
	file, err := os.ReadFile("crypto/crip.json")
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	if config.EncryptionKey == "" {
		return nil, errors.New("[crypto] key is empty, have you set it apppropiately?")
	}
	if config.Salt == "" {
		return nil, errors.New("[crypto] salt is empty, have you set it appropiately?")
	}
	key, _ := pbkdf2.Key(sha512.New, config.EncryptionKey, []byte(config.Salt), 4096, 32)
	return key, nil
}

func Encrypt(plaintext string) (string, error) {
	plaintextBytes := []byte(plaintext)
	key, err := loadKey()
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

func Decrypt(ciphertext string) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	key, err := loadKey()
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
