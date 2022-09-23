package security

import (
	"os"
	"reflect"
	"testing"
)

func TestGeneratePasswordHash(t *testing.T) {
	tests := []struct {
		name         string
		key          []byte
		salt         string
		secret       func()
		message      []byte
		expectedHash []uint8
		expectedErr  bool
	}{
		{
			name:         "valid key and message generate deterministic hash",
			key:          []byte("foo"),
			message:      []byte("bar"),
			salt:         "salt the snail",
			secret:       func() { os.Setenv("SECRET", "z") },
			expectedHash: []uint8{17, 202, 120, 200, 17, 209, 112, 113, 156, 136, 172, 248, 165, 51, 216, 98, 102, 216, 28, 146, 70, 152, 0, 93, 50, 57, 23, 241, 62, 212, 86, 222, 37, 51, 224, 123, 128, 35, 49, 45, 161, 225, 233, 89, 41, 5, 19, 200, 165, 154, 98, 227, 169, 182, 200, 236, 117, 141, 43, 169, 10, 204, 201, 140},
			expectedErr:  false,
		},
		{
			name:        "empty SECRET environment variable generates error",
			key:         []byte("foo"),
			message:     []byte("bar"),
			salt:        "salt the snail",
			secret:      func() { os.Setenv("SECRET", "") },
			expectedErr: true,
		},
		{
			name:        "empty message error",
			key:         []byte("foo"),
			message:     []byte(""),
			salt:        "salt the snail",
			secret:      func() { os.Setenv("SECRET", "valid") },
			expectedErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.secret()
			sut := NewHasher()
			actualHash, err := sut.GeneratePasswordHash(test.message, test.salt)
			if !reflect.DeepEqual(test.expectedHash, actualHash) {
				t.Fatalf("expected hash: '%+v', doesn't match actual: '%+v'", test.expectedHash, actualHash)
			}
			if err != nil && !test.expectedErr {
				t.Fatalf("expected no error for test case but error was encountered: %s", err)
			}
		})
	}
}
