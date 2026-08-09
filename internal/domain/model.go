package domain

import "time"

type RegisterResult struct {
	SessionID             string    // `json:"sessionId"`
	AccessToken           string    // `json:"accessToken"`
	RefreshToken          string    // `json:"refreshToken"`
	AccessTokenExpiresAt  time.Time // `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time // `json:"refreshTokenExpiresAt"`
	Username              string    // `json:"username"`
	UserID                int64     // `json:"userId"`
}

type LoginResult struct {
	SessionID             string    // `json:"sessionId"`
	AccessToken           string    // `json:"accessToken"`
	RefreshToken          string    // `json:"refreshToken"`
	AccessTokenExpiresAt  time.Time // `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time // `json:"refreshTokenExpiresAt"`
	Username              string    // `json:"username"`
	UserID                int64     // `json:"userId"`
}

type RenewAccessTokenResult struct {
	AccessToken          string    // `json:"accessToken"`
	AccessTokenExpiresAt time.Time // `json:"accessTokenExpiresAt"`
}
