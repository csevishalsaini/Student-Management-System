package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(userId int, userName, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtEpiresIn := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"uid":  userId,
		"user": userName,
		"role": role,
	}

	if jwtEpiresIn != "" {
		duration, err := time.ParseDuration(jwtEpiresIn)
		if err != nil {
			return "", ErrorHandler(err, "Internal server error")
		}
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(duration))
	} else {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(50 * time.Minute))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", ErrorHandler(err, "internal error")
	}

	return signedToken, nil

}
