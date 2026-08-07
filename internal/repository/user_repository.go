package repository

import (
	"context"
	"database/sql"
	"fmt"

	sqlc "github.com/meads/firstly-api/db/sqlc"
	"github.com/meads/firstly-api/internal/domain"
)

// Repository defines what the database layer must do.
type UserRepository interface {
	CreateUser(ctx context.Context, username, password string) (*domain.User, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	ListUsers(ctx context.Context, params domain.ListUsersParams) ([]domain.User, error)
	UpdateUserPassword(ctx context.Context, id int64, password string) error
	UserExists(ctx context.Context, id int64) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
}

// SQLRepository wraps the standard sql.DB pool and the sqlc Querier.
type UserSQLRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserSQLRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *UserSQLRepository) CreateUser(ctx context.Context, username, password string) (*domain.User, error) {
	params := sqlc.CreateUserParams{Username: username, Password: password}
	sqlcUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        sqlcUser.ID,
		Username:  sqlcUser.Username,
		Password:  sqlcUser.Password,
		CreatedAt: sqlcUser.CreatedAt.Time,
	}, nil
}

func (r *UserSQLRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return r.queries.UsernameExists(ctx, username)
}

func (r *UserSQLRepository) DeleteUser(ctx context.Context, id int64) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *UserSQLRepository) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	sqlcUser, err := r.queries.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}

	return &domain.User{
		ID:        sqlcUser.ID,
		Username:  sqlcUser.Username,
		Password:  sqlcUser.Password,
		CreatedAt: sqlcUser.CreatedAt.Time,
	}, nil
}

func (r *UserSQLRepository) ListUsers(ctx context.Context, params domain.ListUsersParams) ([]domain.User, error) {
	sqlcUsers, err := r.queries.ListUsers(
		ctx, sqlc.ListUsersParams{Limit: params.Limit, Offset: params.Offset})
	if err != nil {
		return nil, err
	}

	// avoiding expensive array resizing operations during the loop
	domainUsers := make([]domain.User, 0, len(sqlcUsers))

	for _, u := range sqlcUsers {
		domainUsers = append(domainUsers, domain.User{
			ID:        u.ID,
			Username:  u.Username,
			Password:  u.Password,
			CreatedAt: u.CreatedAt.Time,
		})
	}

	return domainUsers, nil
}

func (r *UserSQLRepository) UpdateUserPassword(ctx context.Context, id int64, password string) error {
	sqlcParams := sqlc.UpdateUserPasswordParams{ID: id, Password: password}
	return r.queries.UpdateUserPassword(ctx, sqlcParams)
}

func (r *UserSQLRepository) UserExists(ctx context.Context, id int64) (bool, error) {
	return r.queries.UserExists(ctx, id)
}

func (r *UserSQLRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	sqlcUser, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}

	return &domain.User{
		ID:        sqlcUser.ID,
		Username:  sqlcUser.Username,
		Password:  sqlcUser.Password,
		CreatedAt: sqlcUser.CreatedAt.Time,
	}, nil
}
