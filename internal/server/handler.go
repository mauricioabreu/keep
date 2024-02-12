package server

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mauricioabreu/keep/internal/db"
)

type Note struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type NoteResponse struct {
	ID      int32  `json:"id"`
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

	validate := validator.New()
	err := validate.Struct(note)

	if err != nil {
		validationErrors := []ErrorDetail{}
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ErrorDetail{
				Reason:  err.StructField(),
				Message: "Field validation for '" + err.StructField() + "' failed on the '" + err.Tag() + "' tag",
			})
		}

		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: Error{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request",
				Details: validationErrors,
			},
		})
	}

	createdNote, err := h.dbq.CreateNote(c.Request().Context(), db.CreateNoteParams{
		Title:   note.Title,
		Content: note.Content,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: Error{
				Code:    "INTERNAL_ERROR",
				Message: "Internal error",
			},
		})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Note created",
		Data:    NoteResponse{ID: createdNote.ID, Title: createdNote.Title, Content: createdNote.Content},
	})
}
