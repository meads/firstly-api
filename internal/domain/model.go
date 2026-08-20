package domain

type RegisterResult struct {
	RefreshToken string
	AccessToken  string
	SessionID    string
	UserID       int64
}

type LoginResult struct {
	RefreshToken string
	AccessToken  string
	SessionID    string
	UserID       int64
}

type RenewAccessTokenResult struct {
	AccessToken string
}
