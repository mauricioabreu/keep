package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mauricioabreu/keep/internal/config"
	"github.com/mauricioabreu/keep/internal/db"
	"go.uber.org/fx"
)

func New(q *db.Queries) *echo.Echo {
	e := echo.New()
	e.POST("/notes", func(c echo.Context) error {
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
