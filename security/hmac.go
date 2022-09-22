package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"os"
)

// Generate a salt string with 16 bytes of crypto/rand data.
func generateSalt() string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(randomBytes)
}

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
func IsValidPasswordHash(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func GeneratePasswordHash(message []byte, salt string) ([]byte, error) {
	if len(message) > 0 {
		secret := os.Getenv("SECRET")
		if len(secret) > 0 {
			mac := hmac.New(sha512.New, []byte(secret))
			mac.Write(append(message, salt...))
			return mac.Sum(nil), nil
		}
		return nil, errors.New("SECRET env variable not set")
	}
	return nil, errors.New("message is required")
}
