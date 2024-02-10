package cmd

import (
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
		fx.New(server.Module).Run()
	},
}
