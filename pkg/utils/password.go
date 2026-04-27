package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func VerifyPassword(password, encodedHash string) error {
	parts := strings.Split(encodedHash, ".")

	if len(parts) != 2 {
		return ErrorHandler(errors.New("Invalid Hash Format "), "Internal server error")
	}
	saltBase64 := parts[0]
	hashPasswordBase64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)

	if err != nil {
		return ErrorHandler(errors.New("failed to decode salt"), "failed decode the salt")
	}

	hashPassword, err := base64.StdEncoding.DecodeString(hashPasswordBase64)

	if err != nil {
		return ErrorHandler(errors.New("failed to decode hash"), "failed decode the hash")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	if len(hash) != len(hashPassword) {
		return ErrorHandler(errors.New("hash length mismatch "), "incorrect Password")
	}

	if subtle.ConstantTimeCompare(hash, hashPassword) != 1 {
		return ErrorHandler(errors.New("incorrect password "), "incorrect Password")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrorHandler(errors.New("password is blank "), "Please enter password")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)

	if err != nil {
		return "", ErrorHandler(errors.New("failed to generate salt "), "internal error")
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hasBase64 := base64.StdEncoding.EncodeToString(hash)

	encodeHash := fmt.Sprintf("%s.%s", saltBase64, hasBase64)
	password = encodeHash
	return password, nil
}


