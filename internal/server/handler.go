package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mauricioabreu/keep/internal/db"
)

type Note struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteHandler struct {
	dbq *db.Queries
}

func NewNoteHandler(dbq *db.Queries) *NoteHandler {
	return &NoteHandler{dbq: dbq}
}

func (h *NoteHandler) CreateNote(c echo.Context) error {
	note := new(Note)
	if err := c.Bind(note); err != nil {
		return err
	}

	_, err := h.dbq.CreateNote(c.Request().Context(), db.CreateNoteParams{
		Title:   note.Title,
		Content: note.Content,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create note")
	}

	return c.String(http.StatusOK, "Note created")
}
