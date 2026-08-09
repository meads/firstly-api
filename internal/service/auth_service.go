package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/meads/firstly-api/internal/domain"
	security "github.com/meads/firstly-api/internal/security"
)

// type AuthServicer interface {
// 	Register(ctx context.Context, username, password string) (*domain.RegisterResult, error)
// 	Login(ctx context.Context, username, password string) (*domain.LoginResult, error)
// 	Logout(ctx context.Context, sessionID string) error
// 	RenewAccessToken(ctx context.Context, refreshToken string) (*domain.RenewAccessTokenResult, error)
// 	RevokeSession(ctx context.Context, sessionID string) error
// }

type AuthService struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
	tokener     security.Tokener
	hasher      security.Hasher
}

func NewAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	tokener security.Tokener,
	hasher security.Hasher,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		tokener:     tokener,
		hasher:      hasher,
	}
}

func (s *AuthService) Register(ctx context.Context, username, password string) (*domain.RegisterResult, error) {
	usernameExists, err := s.userRepo.UsernameExists(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error calling user repository from auth service: %w", err)
	}
	if usernameExists {
		return nil, errors.New("please choose another username")
	}

	password, err = s.hasher.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	dbUser, err := s.userRepo.CreateUser(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("error calling user repository create user in auth service: %w", err)
	}

	accessToken, accessClaims, err := s.tokener.GenerateToken(dbUser.ID, dbUser.Username, "access", 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("error generating access token in auth service: %w", err)
	}

	refreshToken, refreshClaims, err := s.tokener.GenerateToken(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token in auth service: %w", err)
	}

	session, err := s.sessionRepo.CreateSession(ctx, domain.CreateSessionParams{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserID:       dbUser.ID,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating session in auth service: %w", err)
	}

	return &domain.RegisterResult{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		Username:              dbUser.Username,
		UserID:                dbUser.ID,
	}, nil

}

func (s *AuthService) Login(ctx context.Context, username, password string) (*domain.LoginResult, error) {
	dbUser, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = s.hasher.ComparePassword(dbUser.Password, password)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password: %w", err)
	}

	accessToken, accessClaims, err := s.tokener.GenerateToken(dbUser.ID, dbUser.Username, "access", 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("error creating access token: %w", err)
	}

	refreshToken, refreshClaims, err := s.tokener.GenerateToken(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("error creating refresh token: %w", err)
	}

	session, err := s.sessionRepo.CreateSession(ctx, domain.CreateSessionParams{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserID:       dbUser.ID,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}
	// ctx.Writer.Header().Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	return &domain.LoginResult{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		Username:              dbUser.Username,
		UserID:                dbUser.ID,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	err := s.sessionRepo.DeleteSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("logout error, error deleting session %w", err)
	}
	return nil
}

func (s *AuthService) RenewAccessToken(ctx context.Context, refreshToken string) (*domain.RenewAccessTokenResult, error) {
	refreshClaims, err := s.tokener.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying refresh token: %w", err)
	}

	session, err := s.sessionRepo.GetSession(ctx, refreshClaims.RegisteredClaims.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	if session.IsRevoked {
		return nil, errors.New("session revoked")
	}

	accessToken, accessClaims, err := s.tokener.GenerateToken(
		refreshClaims.ID, refreshClaims.Username, "access", 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("error creating token: %w", err)
	}

	return &domain.RenewAccessTokenResult{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, sessionID string) error {
	err := s.sessionRepo.RevokeSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("error revoking session: %w", err)
	}

	return nil
}
