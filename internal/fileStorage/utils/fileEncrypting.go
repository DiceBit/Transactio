package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize  = 16
	keySize   = 32
	iteration = 10000
)

func EncryptData(data []byte, password string) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(password), salt, iteration, keySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	padding := aes.BlockSize - len(data)%aes.BlockSize
	for i := 0; i < padding; i++ {
		data = append(data, byte(padding))
	}

	ciphertext := make([]byte, saltSize+aes.BlockSize+len(data))
	copy(ciphertext[:saltSize], salt)
	copy(ciphertext[saltSize:saltSize+aes.BlockSize], iv)

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[saltSize+aes.BlockSize:], data)

	return ciphertext, nil
}
func DecryptData(encryptedData []byte, password string) ([]byte, error) {
	if len(encryptedData) < saltSize+aes.BlockSize {
		return nil, errors.New("invalid encrypted data")
	}

	salt := encryptedData[:saltSize]
	iv := encryptedData[saltSize : saltSize+aes.BlockSize]
	ciphertext := encryptedData[saltSize+aes.BlockSize:]

	key := pbkdf2.Key([]byte(password), salt, iteration, keySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("invalid ciphertext length")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	padding := int(plaintext[len(plaintext)-1])

	if padding > aes.BlockSize || padding == 0 {
		return nil, errors.New("invalid padding size")
	}
	for i := len(plaintext) - padding; i < len(plaintext); i++ {
		if plaintext[i] != byte(padding) {
			return nil, errors.New("invalid padding")
		}
	}
	plaintext = plaintext[:len(plaintext)-padding]

	return plaintext, err
}
