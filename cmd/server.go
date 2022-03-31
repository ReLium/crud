package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ReLium/crud/cmd/server"
	"github.com/ReLium/crud/internal/mongodb"
	"github.com/ReLium/crud/internal/repository"
)

const DefaultMongoUrl = "mongodb://admin:pass@127.0.0.1:27017/"
const DefaultMongoTimeoutMilliseconds = 1000

func Server() error {
	mongodbHost, ok := os.LookupEnv("CRUD_MONGODB_HOST")
	if !ok {
		log.Fatal("Please provide ENV CRUD_MONGODB_HOST")
	}
	mongodbTimeoutMsec, ok := os.LookupEnv("CRUD_MONGODB_TIMEOUT_MSEC")
	if !ok {
		log.Fatal("Please provide ENV CRUD_MONGODB_TIMEOUT_MSEC")
	}
	timeout, err := strconv.Atoi(mongodbTimeoutMsec)
	if err != nil {
		return err
	}

	mongoDBClient, err := mongodb.NewClient(mongodbHost, timeout)
	if err != nil {
		return err
	}
	s := server.NewServer(repository.NewMongoDBRepo(mongoDBClient))

	port, ok := os.LookupEnv("CRUD_SERVER_PORT")
	if !ok {
		log.Fatal("Please provide ENV CRUD_SERVER_PORT")
	}
	return s.Serve(fmt.Sprintf(":%s", port))
}
