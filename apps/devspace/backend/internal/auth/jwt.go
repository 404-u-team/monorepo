package auth

import (
	"strconv"
	"time"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/golang-jwt/jwt"
)

func CreateJWT(secret []byte, userID int, config *config.Config) (string, error) {
	expiration := time.Second * time.Duration(config.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
