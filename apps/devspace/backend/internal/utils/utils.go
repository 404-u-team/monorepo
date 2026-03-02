package utils

import (
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
)

func CreateToken(secret string, userId int, config *config.Config) (string, error) {
	secretConverted := []byte(secret)
	token, err := auth.CreateJWT(secretConverted, userId, config)
	if err != nil {
		return "", err
	}

	return token, nil
}
