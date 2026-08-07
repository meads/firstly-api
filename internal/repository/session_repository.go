package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sqlc "github.com/meads/firstly-api/db/sqlc"
	"github.com/meads/firstly-api/internal/domain"
)

// Repository defines what the database layer must do.
type SessionRepository interface {
	CreateSession(ctx context.Context, arg domain.CreateSessionParams) (*domain.Session, error)
	GetSession(ctx context.Context, id string) (*domain.Session, error)
	DeleteSession(ctx context.Context, id string) error
	RevokeSession(ctx context.Context, id string) error
	// RevokeUserSessions(ctx context.Context, userID int64) error

	// CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	// DeleteSession(ctx context.Context, id string) error
	// GetSession(ctx context.Context, id string) (Session, error)
	// RevokeSession(ctx context.Context, id string) error
	// RevokeUserSessions(ctx context.Context, userID int64) error
}

// SQLRepository wraps the standard sql.DB pool and the sqlc Querier.
type SessionSQLRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &SessionSQLRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *SessionSQLRepository) CreateSession(ctx context.Context, param domain.CreateSessionParams) (*domain.Session, error) {
	sqlcParam := sqlc.CreateSessionParams{
		ID:           param.ID,
		UserID:       param.UserID,
		RefreshToken: param.RefreshToken,
		IsRevoked:    param.IsRevoked,
		ExpiresAt:    param.ExpiresAt,
	}
	sqlcSession, err := r.queries.CreateSession(ctx, sqlcParam)
	if err != nil {
		return nil, err
	}

	return &domain.Session{
		ID:           sqlcSession.ID,
		UserID:       sqlcSession.UserID,
		RefreshToken: sqlcSession.RefreshToken,
		IsRevoked:    sqlcSession.IsRevoked,
		ExpiresAt:    sqlcSession.ExpiresAt,
		// CreatedAt:    sqlcParam.CreatedAt,
	}, nil
}

func (r *SessionSQLRepository) GetSession(ctx context.Context, id string) (*domain.Session, error) {
	sqlcSession, err := r.queries.GetSession(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid session: %w", err)
		}
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	return &domain.Session{
		ID:           sqlcSession.ID,
		UserID:       sqlcSession.UserID,
		RefreshToken: sqlcSession.RefreshToken,
		IsRevoked:    sqlcSession.IsRevoked,
		ExpiresAt:    sqlcSession.ExpiresAt,
		// CreatedAt:    sqlcParam.CreatedAt,
	}, nil
}

func (r *SessionSQLRepository) DeleteSession(ctx context.Context, id string) error {
	return r.queries.DeleteSession(ctx, id)
}

func (r *SessionSQLRepository) RevokeSession(ctx context.Context, id string) error {
	return r.queries.RevokeSession(ctx, id)
}
