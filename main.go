package main

import (
	"log"
	"os"

	"github.com/ReLium/crud/cmd"
)

func main() {
	logFilePath, ok := os.LookupEnv("CRUD_LOG_FILE")
	if !ok {
		log.Fatal("Please provide ENV CRUD_LOG_FILE")
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}

	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.LstdFlags)
	cmd.Root()
}
