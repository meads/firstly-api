package domain

import "time"

type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateNoteParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int64  `json:"userId"`
}

type CreateNoteResult struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int64  `json:"userId"`
}

type DeleteNoteParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
}

type UpdateNoteParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	ID      int64  `json:"id"`
	UserID  int64  `json:"userId"`
}
