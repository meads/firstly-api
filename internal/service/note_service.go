package service

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/meads/firstly-api/internal/domain"
)

// type NoteServicer interface {
// 	CreateNote(ctx context.Context, userID int64, title, content string) (*domain.Note, error)
// 	UpdateNote(ctx context.Context, noteID, userID int64, title, content string) (*domain.Note, error)
// 	ListNotes(ctx context.Context, userID int64) ([]domain.Note, error)
// 	DeleteNote(ctx context.Context, noteID, userID int64) error
// }

type NoteService struct {
	noteRepo domain.NoteRepository
	userRepo domain.UserRepository
}

func NewNoteService(
	noteRepo domain.NoteRepository,
	userRepo domain.UserRepository) *NoteService {
	return &NoteService{
		noteRepo: noteRepo,
		userRepo: userRepo,
	}
}

func (ns *NoteService) CreateNote(ctx context.Context, userID int64, title, content string) (*domain.Note, error) {
	userExists, err := ns.userRepo.UserExists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("note service error calling user exists: %w", err)
	}

	if !userExists {
		return nil, errors.New("note service error user not found")
	}

	dbNote, err := ns.noteRepo.CreateNote(ctx, domain.CreateNoteParams{
		Content: content,
		Title:   title,
		UserID:  userID,
	})
	if err != nil {
		return nil, fmt.Errorf("note service error calling create note: %w", err)
	}

	return &domain.Note{
		Content:   dbNote.Content,
		CreatedAt: dbNote.CreatedAt,
		ID:        dbNote.ID,
		Title:     dbNote.Title,
		UserID:    dbNote.UserID,
	}, nil
}

func (ns *NoteService) UpdateNote(ctx context.Context, noteID, userID int64, title, content string) (*domain.Note, error) {
	userExists, err := ns.userRepo.UserExists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("note service error calling user exists: %w", err)
	}

	if !userExists {
		return nil, errors.New("note service error user not found")
	}

	err = ns.noteRepo.UpdateNote(ctx, domain.UpdateNoteParams{
		Content: content,
		ID:      noteID,
		Title:   title,
		UserID:  userID,
	})
	if err != nil {
		return nil, fmt.Errorf("note service error calling update note: %w", err)
	}

	return &domain.Note{
		Content: content,
		ID:      noteID,
		Title:   title,
		UserID:  userID,
	}, nil
}

func (ns *NoteService) ListNotes(ctx context.Context, userID int64) ([]domain.Note, error) {
	userExists, err := ns.userRepo.UserExists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("note service error calling user exists: %w", err)
	}

	if !userExists {
		return nil, errors.New("note service error user not found")
	}

	notes, err := ns.noteRepo.ListNotesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("note service error calling list notes by userid: %w", err)
	}

	return notes, nil
}

func (ns *NoteService) DeleteNote(ctx context.Context, noteID, userID int64) error {
	userExists, err := ns.userRepo.UserExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("note service error calling user exists: %w", err)
	}

	if !userExists {
		return errors.New("note service error user not found")
	}

	err = ns.noteRepo.DeleteNote(ctx, domain.DeleteNoteParams{
		ID:     noteID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("note service error calling delete note: %w", err)
	}

	return nil
}
