package utils

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/google/uuid"
)

func CreateToken(secret string, userId uuid.UUID, expirationTime int) (string, error) {
	secretConverted := []byte(secret)
	token, err := auth.CreateJWT(secretConverted, userId, expirationTime)
	if err != nil {
		return "", err
	}

	return token, nil
}
