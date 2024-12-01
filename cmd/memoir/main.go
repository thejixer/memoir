package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/thejixer/memoir/internal/handlers"
	server "github.com/thejixer/memoir/internal/server"
)

func init() {
	godotenv.Load()
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")

	handlerService := handlers.NewHandlerService()

	server := server.NewAPIServer(listenAddr, handlerService)
	server.Run()
}
