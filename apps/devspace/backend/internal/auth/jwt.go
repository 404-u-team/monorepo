package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJWT(secret []byte, userID uint, expirationTime int) (string, error) {
	expiration := time.Second * time.Duration(expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.FormatUint(uint64(userID), 10), // стандарт в jwt, подразумевает userID
		"exp": time.Now().Add(expiration).Unix(),      // время конца жизни
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
