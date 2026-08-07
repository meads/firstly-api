package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword string, password string) error
}

type HashLib struct{}

func NewHasher() Hasher {
	return &HashLib{}
}

func (h *HashLib) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password %w", err)
	}

	return string(hashed), nil
}

func (h *HashLib) ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
