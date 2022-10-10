package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"os"
)

type Hasher interface {
	GenerateSalt() string
	IsValidPassword(phrase []byte, salt, password string) (bool, error)
	GeneratePasswordHash(phrase []byte, salt string) ([]byte, error)
}

type HashLib struct{}

func NewHasher() Hasher {
	return &HashLib{}
}

// Generate a salt string with 4096 bytes of crypto/rand data.
func (*HashLib) GenerateSalt() string {
	randomBytes := make([]byte, 4096)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(randomBytes)
}

// IsValidPassword reports whether password generates a valid HMAC matching the stored phrase.
func (*HashLib) IsValidPassword(phrase []byte, salt, password string) (bool, error) {
	secret := os.Getenv("SECRET")
	if len(secret) == 0 {
		return false, errors.New("SECRET env variable not set")
	}
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(append([]byte(password), salt...))
	return hmac.Equal(mac.Sum(nil), phrase), nil
}

func (*HashLib) GeneratePasswordHash(phrase []byte, salt string) ([]byte, error) {
	if len(phrase) > 0 {
		secret := os.Getenv("SECRET")
		if len(secret) > 0 {
			mac := hmac.New(sha512.New, []byte(secret))
			mac.Write(append(phrase, salt...))
			return mac.Sum(nil), nil
		}
		return nil, errors.New("SECRET env variable not set")
	}
	return nil, errors.New("phrase is required")
}
