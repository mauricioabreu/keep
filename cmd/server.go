package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mauricioabreu/keep/internal/config"
	"github.com/mauricioabreu/keep/internal/db"
	"github.com/mauricioabreu/keep/internal/logging"
	"github.com/mauricioabreu/keep/internal/server"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Keep server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Could not load .env file")
		}

		app := fx.New(
			fx.Provide(config.New),
			fx.Provide(logging.New),
			fx.Provide(db.NewPool),
			fx.Provide(db.NewQuerier),
			server.Module)
		app.Run()
	},
}
