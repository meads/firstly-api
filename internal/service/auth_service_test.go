package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	domain "github.com/meads/firstly-api/internal/domain"
	"github.com/meads/firstly-api/internal/security"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Register(t *testing.T) {
	var testRegisterResult *domain.RegisterResult
	tests := []struct {
		name              string
		username          string
		password          string
		expectedError     error
		setupExpectations func(
			tokener *security.MockTokener, hasher *security.MockHasher,
			srepo *MockSessionRepository, urepo *MockUserRepository,
		)
	}{
		{
			name:          "auth service register fails given username exists returns an error",
			username:      "username",
			password:      "password",
			expectedError: errors.New("server error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, errors.New("server error"))

				testRegisterResult = nil
			},
		},
		{
			name:          "auth service register fails given username not available",
			username:      "username",
			password:      "password",
			expectedError: errors.New("please choose another username"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(true, nil)

				testRegisterResult = nil
			},
		},
		{
			name:          "auth service register fails given error hashing password",
			username:      "username",
			password:      "password",
			expectedError: errors.New("error hashing password"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("", errors.New("error hashing password"))

				testRegisterResult = nil
			},
		},
		{
			name:          "auth service register fails given error creating user",
			username:      "username",
			password:      "password",
			expectedError: errors.New("user repo create error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("hashpassword", nil)
				urepo.EXPECT().CreateUser(gomock.Any(), "username", "hashpassword").
					Return(nil, errors.New("user repo create error"))

				testRegisterResult = nil
			},
		},
		{
			name:          "auth service register fails given error generating access token",
			username:      "username",
			password:      "password",
			expectedError: errors.New("tokener error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("hashpassword", nil)
				urepo.EXPECT().CreateUser(gomock.Any(), "username", "hashpassword").
					Return(&domain.User{ID: int64(1), Username: "username", Password: "hashpassword"}, nil)

				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("", nil, errors.New("tokener error"))

				testRegisterResult = nil
			},
		},
		{
			name:          "auth service register fails given error generating refresh token",
			username:      "username",
			password:      "password",
			expectedError: errors.New("tokener error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("hashpassword", nil)
				urepo.EXPECT().CreateUser(gomock.Any(), "username", "hashpassword").
					Return(&domain.User{ID: int64(1), Username: "username", Password: "hashpassword"}, nil)

				accessClaims, _ := security.NewUserClaims(int64(1), "username", "access", 15*time.Minute)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("accesstoken", accessClaims, nil)

				tokener.EXPECT().
					GenerateToken(int64(1), "username", "refresh", 24*time.Hour).
					Return("", nil, errors.New("tokener error"))

				testRegisterResult = nil
			},
		},
		{
			name:          "auth service register fails given error creating session",
			username:      "username",
			password:      "password",
			expectedError: errors.New("server error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("hashpassword", nil)
				urepo.EXPECT().CreateUser(gomock.Any(), "username", "hashpassword").
					Return(&domain.User{ID: int64(1), Username: "username", Password: "hashpassword"}, nil)

				accessClaims, _ := security.NewUserClaims(int64(1), "username", "access", 15*time.Minute)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("accesstoken", accessClaims, nil)

				refreshClaims, _ := security.NewUserClaims(int64(1), "username", "access", 24*time.Hour)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "refresh", 24*time.Hour).
					Return("refreshtoken", refreshClaims, nil)

				createSessionParams := domain.CreateSessionParams{
					ID:           refreshClaims.RegisteredClaims.ID,
					UserID:       int64(1),
					RefreshToken: "refreshtoken",
					IsRevoked:    false,
					ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
				}
				srepo.EXPECT().
					CreateSession(gomock.Any(), createSessionParams).
					Return(nil, errors.New("server error"))
				testRegisterResult = nil
			},
		},
		{
			name:          "auth service register succeeds given valid data supplied",
			username:      "username",
			password:      "password",
			expectedError: nil,
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("hashpassword", nil)
				urepo.EXPECT().CreateUser(gomock.Any(), "username", "hashpassword").
					Return(&domain.User{ID: int64(1), Username: "username", Password: "hashpassword"}, nil)

				accessClaims, _ := security.NewUserClaims(int64(1), "username", "access", 15*time.Minute)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("accesstoken", accessClaims, nil)

				refreshClaims, _ := security.NewUserClaims(int64(1), "username", "access", 24*time.Hour)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "refresh", 24*time.Hour).
					Return("refreshtoken", refreshClaims, nil)

				createSessionParams := domain.CreateSessionParams{
					ID:           refreshClaims.RegisteredClaims.ID,
					UserID:       int64(1),
					RefreshToken: "refreshtoken",
					IsRevoked:    false,
					ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
				}
				srepo.EXPECT().
					CreateSession(gomock.Any(), createSessionParams).
					Return(&domain.Session{
						ID:           refreshClaims.RegisteredClaims.ID,
						UserID:       int64(1),
						RefreshToken: "refreshtoken",
						IsRevoked:    false,
						ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
					}, nil)
				testRegisterResult = &domain.RegisterResult{
					SessionID:             refreshClaims.RegisteredClaims.ID,
					AccessToken:           "accesstoken",
					RefreshToken:          "refreshtoken",
					AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
					RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
					Username:              "username",
					UserID:                int64(1),
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)

			tokener := security.NewMockTokener(ctrl)
			hasher := security.NewMockHasher(ctrl)
			sessionRepo := NewMockSessionRepository(ctrl)
			userRepo := NewMockUserRepository(ctrl)
			test.setupExpectations(tokener, hasher, sessionRepo, userRepo)

			authService := NewAuthService(userRepo, sessionRepo, tokener, hasher)

			// Act
			registerResult, err := authService.Register(context.Background(), test.username, test.password)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}

			if !reflect.DeepEqual(testRegisterResult, registerResult) {
				t.Fatalf("register result does not match, expected: \n%v, got: \n%v", testRegisterResult, registerResult)
			}
		})
	}

}
