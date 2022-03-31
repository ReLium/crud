package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
			err := Server()
			fmt.Print(err)
		},
	}
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Add dummy items to DB",
		Run: func(cmd *cobra.Command, args []string) {
			err := Generate()
			fmt.Print(err)
		},
	}
)

// Execute executes the root command.
func Root() error {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(generateCmd)
	return rootCmd.Execute()
}

// @TODO: Refine config injection logic.
func getMongoDBSettings() (host string, timeoutMsec int) {
	host, ok := os.LookupEnv("CRUD_MONGODB_HOST")
	if !ok {
		log.Fatal("Please provide ENV CRUD_MONGODB_HOST")
	}
	mongodbTimeoutMsec, ok := os.LookupEnv("CRUD_MONGODB_TIMEOUT_MSEC")
	if !ok {
		log.Fatal("Please provide ENV CRUD_MONGODB_TIMEOUT_MSEC")
	}
	timeoutMsec, err := strconv.Atoi(mongodbTimeoutMsec)
	if err != nil {
		log.Fatal("Invalid CRUD_MONGODB_TIMEOUT_MSEC")
	}
	return host, timeoutMsec
}
