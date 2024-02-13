package db

import (
	"context"

	"github.com/google/uuid"
)

type NoteStorer interface {
	CreateNote(context.Context, CreateNoteParams) (Note, error)
	GetNote(context.Context, uuid.UUID) (Note, error)
	ListNotes(context.Context) ([]Note, error)
}
