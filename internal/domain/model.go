package domain

import "time"

type RegisterResult struct {
	SessionID             string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	Username              string
	UserID                int64
}

type LoginResult struct {
	SessionID             string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	Username              string
	UserID                int64
}

type RenewAccessTokenResult struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
}
