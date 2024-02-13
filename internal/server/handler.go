package server

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (h *NoteHandler) GetNote(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: Error{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request",
				Details: []ErrorDetail{
					{
						Reason:  "id",
						Message: "Invalid UUID",
					},
				},
			},
		})
	}

	note, err := h.sdb.GetNote(c.Request().Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: Error{
					Code:    "NOT_FOUND",
					Message: "Note not found",
				},
			})
		}

		h.logger.Errorw("Error getting note", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: Error{
				Code:    "INTERNAL_ERROR",
				Message: "Internal error",
			},
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Note found",
		Data: NoteResponse{
			ID:      note.ID,
			Title:   note.Title,
			Content: note.Content,
		},
	})
}

func (h *NoteHandler) ListNotes(c echo.Context) error {
	notes, err := h.sdb.ListNotes(c.Request().Context())
	if err != nil {
		h.logger.Errorw("Error listing notes", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: Error{
				Code:    "INTERNAL_ERROR",
				Message: "Internal error",
			},
		})
	}

	notesResponse := []NoteResponse{}
	for _, note := range notes {
		notesResponse = append(notesResponse, NoteResponse{
			ID:      note.ID,
			Title:   note.Title,
			Content: note.Content,
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Notes found",
		Data:    notesResponse,
	})
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
