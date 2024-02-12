package db

import "context"

type NoteStorer interface {
	CreateNote(context.Context, CreateNoteParams) (Note, error)
}
