package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "crud",
		Short: "Simple demonstration CRUD application",
		Long:  `Simple demonstration CRUD application that includes web server and command to fill database with dummy data`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start http server",
		Run: func(cmd *cobra.Command, args []string) {
			Server()
		},
	}
)

// Execute executes the root command.
func Root() error {
	rootCmd.AddCommand(serverCmd)
	return rootCmd.Execute()
}
