package domain

import (
	"context"
	"time"
)

type Note struct {
	ID        int64     // `json:"id"`
	Title     string    // `json:"title"`
	Content   string    // `json:"content"`
	UserID    int64     // `json:"userId"`
	CreatedAt time.Time // `json:"createdAt"`
}

type CreateNoteParams struct {
	Title   string // `json:"title"`
	Content string // `json:"content"`
	UserID  int64  // `json:"userId"`
}

type CreateNoteResult struct {
	ID      int64  // `json:"id"`
	Title   string // `json:"title"`
	Content string // `json:"content"`
	UserID  int64  // `json:"userId"`
}

type DeleteNoteParams struct {
	ID     int64 // `json:"id"`
	UserID int64 // `json:"userId"`
}

type UpdateNoteParams struct {
	ID      int64  // `json:"id"`
	Title   string // `json:"title"`
	Content string // `json:"content"`
	UserID  int64  // `json:"userId"`
}

type NoteRepository interface {
	CreateNote(ctx context.Context, param CreateNoteParams) (*Note, error)
	GetNote(ctx context.Context, id int64) (*Note, error)
	DeleteNote(ctx context.Context, param DeleteNoteParams) error
	ListNotesByUserID(ctx context.Context, userID int64) ([]Note, error)
	UpdateNote(ctx context.Context, param UpdateNoteParams) error
}
