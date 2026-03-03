package utils

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
)

func CreateToken(secret string, userId uint, expirationTime int) (string, error) {
	secretConverted := []byte(secret)
	token, err := auth.CreateJWT(secretConverted, userId, expirationTime)
	if err != nil {
		return "", err
	}

	return token, nil
}
