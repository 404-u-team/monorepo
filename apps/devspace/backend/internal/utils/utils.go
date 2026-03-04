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

// создает срез с заданным запасом и вставляет элементы
func MakeSlice[T any](cap int, elems ...T) []T {
	slice := make([]T, 0, cap)

	for _, elem := range elems {
		slice = append(slice, elem)
	}

	return slice
}
