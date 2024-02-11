package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mauricioabreu/keep/internal/config"
	"github.com/mauricioabreu/keep/internal/db"
	"go.uber.org/fx"
)

type Note struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func New(q *db.Queries) *echo.Echo {
	e := echo.New()
	e.POST("/notes", func(c echo.Context) error {
		note := new(Note)
		if err := c.Bind(note); err != nil {
			return err
		}

		_, err := q.CreateNote(context.Background(), db.CreateNoteParams{
			Title:   note.Title,
			Content: note.Content,
		})

		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create note")
		}

		return c.String(http.StatusOK, "Note created")
	})
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	return e
}

func registerHooks(lc fx.Lifecycle, cfg *config.Config, e *echo.Echo) {
	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go func() {
					if err := e.Start(":" + cfg.ServerPort); err != nil {
						e.Logger.Info("shutting down the server")
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return e.Shutdown(ctx)
			},
		},
	)
}
