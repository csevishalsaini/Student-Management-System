package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
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
