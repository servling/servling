package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "servling",
		Short: "The container deployment platform that just works",
		Long:  "Servling transforms container deployment from a chore into a breeze. Whether you're a solo developer or managing enterprise applications, Servling adapts to your workflow.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	rootCmd.AddCommand(ServeCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
