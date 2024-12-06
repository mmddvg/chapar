package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"mmddvg/chapar/pkg/errs"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Encrypt(username, password string) (string, error) {

	key := sha256.Sum256([]byte(password))

	// Generate an AES cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", errs.NewUnexpected(err)
	}

	// Create a GCM (Galois/Counter Mode) cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errs.NewUnexpected(err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	copy(nonce, key[:aesGCM.NonceSize()])

	ciphertext := aesGCM.Seal(nil, nonce, []byte(username), nil)

	encrypted := base64.StdEncoding.EncodeToString(ciphertext)

	return encrypted, nil
}

func CheckPassword(encryptedUsername, password, originalUsername string) (bool, error) {

	key := sha256.Sum256([]byte(password))

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedUsername)
	if err != nil {
		return false, errs.NewUnexpected(errors.New("failed to decode base64 : " + err.Error()))
	}

	// Generate an AES cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return false, errs.NewUnexpected(err)
	}

	// Create a GCM (Galois/Counter Mode) cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return false, errs.NewUnexpected(err)
	}

	// Derive the same nonce as used during encryption
	nonce := make([]byte, aesGCM.NonceSize())
	copy(nonce, key[:aesGCM.NonceSize()])

	// Decrypt the ciphertext
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false, nil // wrong password
	}

	if string(plaintext) == originalUsername {
		return true, nil
	}

	return false, nil
}

func GenerateJWT(id uint64) (string, error) {
	claims := jwt.MapClaims{
		"sub": fmt.Sprint(id),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		"iat": jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", errs.NewUnexpected(err)
	}

	return signedToken, nil
}
