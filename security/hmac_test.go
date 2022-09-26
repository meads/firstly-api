package security

import (
	"os"
	"reflect"
	"testing"
)

func TestGeneratePasswordHash(t *testing.T) {
	tests := []struct {
		name         string
		salt         string
		secret       func()
		phrase       []byte
		expectedHash []uint8
		expectedErr  bool
	}{
		{
			name:         "valid key and phrase generate deterministic hash",
			phrase:       []byte("bar"),
			salt:         "salt the snail",
			secret:       func() { os.Setenv("SECRET", "z") },
			expectedHash: []uint8{17, 202, 120, 200, 17, 209, 112, 113, 156, 136, 172, 248, 165, 51, 216, 98, 102, 216, 28, 146, 70, 152, 0, 93, 50, 57, 23, 241, 62, 212, 86, 222, 37, 51, 224, 123, 128, 35, 49, 45, 161, 225, 233, 89, 41, 5, 19, 200, 165, 154, 98, 227, 169, 182, 200, 236, 117, 141, 43, 169, 10, 204, 201, 140},
			expectedErr:  false,
		},
		{
			name:        "empty SECRET environment variable generates error",
			phrase:      []byte("bar"),
			salt:        "salt the snail",
			secret:      func() { os.Setenv("SECRET", "") },
			expectedErr: true,
		},
		{
			name:        "empty phrase error",
			phrase:      []byte(""),
			salt:        "salt the snail",
			secret:      func() { os.Setenv("SECRET", "valid") },
			expectedErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.secret()
			sut := NewHasher()
			actualHash, err := sut.GeneratePasswordHash(test.phrase, test.salt)
			if !reflect.DeepEqual(test.expectedHash, actualHash) {
				t.Fatalf("expected hash: '%+v', doesn't match actual: '%+v'", test.expectedHash, actualHash)
			}
			if err != nil && !test.expectedErr {
				t.Fatalf("expected no error for test case but error was encountered: %s", err)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name          string
		phrase        []byte
		salt          string
		password      string
		secret        func()
		expectedMatch bool
		expectedErr   bool
	}{
		{
			name:          "valid password matches stored phrase when hashed with current salt",
			phrase:        []uint8{17, 202, 120, 200, 17, 209, 112, 113, 156, 136, 172, 248, 165, 51, 216, 98, 102, 216, 28, 146, 70, 152, 0, 93, 50, 57, 23, 241, 62, 212, 86, 222, 37, 51, 224, 123, 128, 35, 49, 45, 161, 225, 233, 89, 41, 5, 19, 200, 165, 154, 98, 227, 169, 182, 200, 236, 117, 141, 43, 169, 10, 204, 201, 140},
			salt:          "salt the snail",
			password:      "bar",
			secret:        func() { os.Setenv("SECRET", "z") },
			expectedMatch: true,
			expectedErr:   false,
		},
		{
			name:          "error is generated when SECRET isn't set",
			phrase:        []uint8{17, 202, 120, 200, 17, 209, 112, 113, 156, 136, 172, 248, 165, 51, 216, 98, 102, 216, 28, 146, 70, 152, 0, 93, 50, 57, 23, 241, 62, 212, 86, 222, 37, 51, 224, 123, 128, 35, 49, 45, 161, 225, 233, 89, 41, 5, 19, 200, 165, 154, 98, 227, 169, 182, 200, 236, 117, 141, 43, 169, 10, 204, 201, 140},
			salt:          "salt the snail",
			password:      "bar",
			secret:        func() { os.Setenv("SECRET", "") },
			expectedMatch: false,
			expectedErr:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.secret()
			sut := NewHasher()
			matches, err := sut.IsValidPassword(test.phrase, test.salt, test.password)
			if err != nil && !test.expectedErr {
				t.Fatalf("expected no error for test case but error was encountered: %s", err)
			}
			if !matches && test.expectedMatch {
				t.Fatalf("expected match but password doesn't match")
			}
		})
	}
}
