package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJWT(secret []byte, userID uint, expirationTime int) (string, error) {
	expiration := time.Second * time.Duration(expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,                            // стандарт в jwt, подразумевает userID
		"exp": time.Now().Add(expiration).Unix(), // время конца жизни
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// возвращает "sub" из токена, если валиден
func ValidateJWT(secret []byte, tokenString string) (uint, error) {
	// парсинг и валидация токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("ожидался другой способ создания подписи токена: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return 0, err
	}

	// получаем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("неккоректный токен")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if float64(time.Now().Unix()) > exp {
			return 0, fmt.Errorf("время жизни токена истекло")
		}
	} else {
		return 0, fmt.Errorf("отсутствует exp внутри токена")
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("отсутствует sub внутри токена")
	}

	return uint(sub), nil
}
