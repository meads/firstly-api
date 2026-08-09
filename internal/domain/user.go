package domain

import "time"

type User struct {
	ID        int64     // `json:"id"`
	Username  string    // `json:"username"`
	Password  string    // `json:"password"`
	CreatedAt time.Time // `json:"createdAt"`
}

type ListUsersParams struct {
	Limit  int32
	Offset int32
}

type UpdateUserPasswordParams struct {
	ID              int64  // `json:"id"`
	Username        string // `json:"username"`
	CurrentPassword string // `json:"currentPassword"`
	NewPassword     string // `json:"newPassword"`
}
