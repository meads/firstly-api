package handler

import (
	"time"

	"github.com/meads/firstly-api/internal/domain"
)

// Auth

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	SessionID             string    `json:"sessionId"`
	AccessToken           string    `json:"accessToken"`
	RefreshToken          string    `json:"refreshToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
	Username              string    `json:"username"`
	UserID                int64     `json:"userId"`
}

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type LoginResponse struct {
	SessionID             string    `json:"sessionId"`
	AccessToken           string    `json:"accessToken"`
	RefreshToken          string    `json:"refreshToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
	Username              string    `json:"username"`
	UserID                int64     `json:"userId"`
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`
}

// Users

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type ListUsersResponse struct {
	Users []domain.User `json:"users"`
}

type PatchUserRequest struct {
	ID              int64  `json:"id" binding:"required"`
	Username        string `json:"username" binding:"required"`
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

type PatchUserResponse struct {
	Message string `json:"message"`
}

// Notes

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  int64  `json:"userId" binding:"required"`
}

// type CreateNoteResponse struct {
// 	ID      int64  `json:"id"`
// 	Title   string `json:"title"`
// 	Content string `json:"content"`
// 	UserID  int64  `json:"userId"`
// }

type UpdateNoteRequest struct {
	ID      int64  `json:"id" binding:"required"`
	UserID  int64  `json:"userId" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ListNotesResponse struct {
	Notes []domain.Note `json:"notes"`
}

type DeleteNoteResponse struct {
	Message string `json:"message"`
}
