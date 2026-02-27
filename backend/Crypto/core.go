package crypto

import (
	system "DevSpace/System"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

/*
type Argon2Params struct {
	Memory      uint
	Iterations  uint
	Parallelism uint
	SaltLength  uint
	KeyLength   uint
}

func CreateArgonConf(c *system.Config) Argon2Params {
	return Argon2Params{
		Memory:      c.Memory,
		Iterations:  c.Iterations,
		Parallelism: c.Parallelism,
		SaltLength:  c.SaltLength,
		KeyLength:   c.KeyLength,
	}
}

*/

func EncodePassword(pass string, c *system.Config) (string, error) {
	salt := make([]byte, c.SaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(pass),
		salt,
		uint32(c.Iterations),
		uint32(c.Memory),
		uint8(c.Parallelism),
		uint32(c.KeyLength),
	)

	// Форматируем строку в стандартном виде
	parts := []string{
		"$argon2id",
		"v=19",
		fmt.Sprintf("m=%d,t=%d,p=%d", c.Memory, c.Iterations, c.Parallelism),
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	}

	return strings.Join(parts, "$"), nil
}

func VerifyPassword(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false, errors.New("invalid argon2 format")
	}

	var m, t, p uint32
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	expected, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	computed := argon2.IDKey([]byte(password), salt, t, m, uint8(p), uint32(len(expected)))

	return hmac.Equal(computed, expected), nil
}
