// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"
)

type Querier interface {
	AccountExists(ctx context.Context, id int64) (bool, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateImage(ctx context.Context, data string) (Image, error)
	DeleteAccount(ctx context.Context, id int64) error
	DeleteImage(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountByUsername(ctx context.Context, username string) (Account, error)
	GetImage(ctx context.Context, id int64) (Image, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListImages(ctx context.Context, arg ListImagesParams) ([]Image, error)
	SoftDeleteAccount(ctx context.Context, id int64) error
	SoftDeleteImage(ctx context.Context, id int64) error
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) error
	UpdateImage(ctx context.Context, arg UpdateImageParams) error
}

var _ Querier = (*Queries)(nil)
