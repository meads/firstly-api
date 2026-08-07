package service

// import (
// 	"context"
// 	"time"

// 	domain "github.com/meads/firstly-api/internal/domain"
// 	r "github.com/meads/firstly-api/internal/repository"
// )

// type SessionService struct {
// 	repo r.SessionRepository
// }

// func NewSessionService(repo r.SessionRepository) *SessionService {
// 	return &SessionService{repo: repo}
// }

// func (s *SessionService) Create(ctx context.Context, id string, userID int64, refreshToken string, expiresAt time.Time) (*domain.Session, error) {

// 	arg := domain.CreateSessionParams{
// 		ID:           id,
// 		UserID:       userID,
// 		RefreshToken: refreshToken,
// 		IsRevoked:    false,
// 		ExpiresAt:    expiresAt,
// 	}

// 	return s.repo.CreateSession(ctx, arg)
// }
