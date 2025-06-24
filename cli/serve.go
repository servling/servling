package cli

import (
	"github.com/servling/servling/pkg/server"
	"github.com/spf13/cobra"
)

func ServeCommand() *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run the servling backend server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Run()
		},
	}

	return serveCmd
}
