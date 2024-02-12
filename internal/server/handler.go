package server

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mauricioabreu/keep/internal/db"
	"go.uber.org/zap"
)

type Note struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type NoteResponse struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

type NoteHandler struct {
	sdb    db.NoteStorer
	logger *zap.SugaredLogger
}

func NewNoteHandler(s db.NoteStorer, logger *zap.SugaredLogger) *NoteHandler {
	return &NoteHandler{sdb: s, logger: logger}
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

	createdNote, err := h.sdb.CreateNote(c.Request().Context(), db.CreateNoteParams{
		Title:   note.Title,
		Content: note.Content,
	})

	if err != nil {
		h.logger.Errorw("Error creating note", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: Error{
				Code:    "INTERNAL_ERROR",
				Message: "Internal error",
			},
		})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Note created",
		Data: NoteResponse{
			ID:      createdNote.ID,
			Title:   createdNote.Title,
			Content: createdNote.Content,
		},
	})
}
