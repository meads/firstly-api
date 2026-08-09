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

func MapToRegisterResponse(registerResult *domain.RegisterResult) *RegisterResponse {
	return &RegisterResponse{
		SessionID:             registerResult.SessionID,
		AccessToken:           registerResult.AccessToken,
		RefreshToken:          registerResult.RefreshToken,
		AccessTokenExpiresAt:  registerResult.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: registerResult.RefreshTokenExpiresAt,
		Username:              registerResult.Username,
		UserID:                registerResult.UserID,
	}
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

func MapToRenewAccessTokenResponse(result *domain.RenewAccessTokenResult) *RenewAccessTokenResponse {
	return &RenewAccessTokenResponse{
		AccessToken:          result.AccessToken,
		AccessTokenExpiresAt: result.AccessTokenExpiresAt,
	}
}

// Users

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	// CreatedAt sql.NullTime `json:"createdAt"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

func MapUsersToListUsersResponse(users []domain.User) ListUsersResponse {
	userResponses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:       u.ID,
			Username: u.Username,
			Password: u.Password,
			// CreatedAt: u.CreatedAt,
		})
	}
	return ListUsersResponse{Users: userResponses}
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

type NoteResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int64  `json:"userId"`
}

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  int64  `json:"userId" binding:"required"`
}

type UpdateNoteRequest struct {
	ID      int64  `json:"id" binding:"required"`
	UserID  int64  `json:"userId" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ListNotesResponse struct {
	Notes []NoteResponse `json:"notes"`
}

type DeleteNoteResponse struct {
	Message string `json:"message"`
}

func MapNotesToListNotesResponse(notes []domain.Note) ListNotesResponse {
	noteResponses := make([]NoteResponse, 0, len(notes))
	for _, n := range notes {
		noteResponses = append(noteResponses, NoteResponse{
			ID:      n.ID,
			Title:   n.Title,
			Content: n.Content,
			UserID:  n.UserID,
		})
	}
	return ListNotesResponse{Notes: noteResponses}
}
