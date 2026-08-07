package service

import (
	"context"
	"fmt"

	domain "github.com/meads/firstly-api/internal/domain"
	"github.com/meads/firstly-api/internal/repository"
	security "github.com/meads/firstly-api/internal/security"
)

type UserServicer interface {
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, params domain.ListUsersParams) ([]domain.User, error)
	ChangePassword(ctx context.Context, params domain.UpdateUserPasswordParams) error
	// UpdateUserPassword(ctx context.Context, params domain.UpdateUserPasswordParams) error
}

type UserService struct {
	userRepo repository.UserRepository
	hasher   security.Hasher
}

func NewUserService(repo repository.UserRepository) UserServicer {
	return &UserService{
		userRepo: repo,
	}
}

func (u *UserService) DeleteUser(ctx context.Context, id int64) error {
	return u.userRepo.DeleteUser(ctx, id)
}

func (u *UserService) ListUsers(ctx context.Context, params domain.ListUsersParams) ([]domain.User, error) {
	return u.userRepo.ListUsers(ctx, params)
}
func (u *UserService) ChangePassword(ctx context.Context, params domain.UpdateUserPasswordParams) error {
	user, err := u.userRepo.GetUser(ctx, params.ID)
	if err != nil {
		return err
	}

	err = u.hasher.ComparePassword(user.Password, params.CurrentPassword)
	if err != nil {
		return fmt.Errorf("current password invalid: %w", err)
	}

	newPasswordHash, err := u.hasher.HashPassword(params.NewPassword)
	if err != nil {
		return fmt.Errorf("error creating password hash: %w", err)
	}

	err = u.userRepo.UpdateUserPassword(ctx, user.ID, newPasswordHash)
	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	return nil
}
