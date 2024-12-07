package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thejixer/memoir/internal/database"
	"github.com/thejixer/memoir/internal/handlers"
	"github.com/thejixer/memoir/internal/redis"
	"github.com/thejixer/memoir/internal/server"
)

func init() {
	godotenv.Load()
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")
	dbStore, err := database.NewPostgresStore()

	if err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	if err := dbStore.Init(); err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	redisStore, err := redis.NewRedisStore()
	if err != nil {
		log.Fatal("could not connect to the redis: ", err)
	}

	handlerService := handlers.NewHandlerService(dbStore, redisStore)

	s := server.NewAPIServer(listenAddr, handlerService)
	s.Run()
}
