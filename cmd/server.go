package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ReLium/crud/cmd/server"
	"github.com/ReLium/crud/internal/mongodb"
	"github.com/ReLium/crud/internal/repository"
)

func Server() error {
	host, timeout := getMongoDBSettings()
	mongoDBClient, err := mongodb.NewClient(host, timeout)
	if err != nil {
		return err
	}
	s := server.NewServer(repository.NewLogWrapper(repository.NewMongoDBRepo(mongoDBClient)))

	port, ok := os.LookupEnv("CRUD_SERVER_PORT")
	if !ok {
		log.Fatal("Please provide ENV CRUD_SERVER_PORT")
	}
	return s.Serve(fmt.Sprintf(":%s", port))
}
