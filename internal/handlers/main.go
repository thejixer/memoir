package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thejixer/memoir/internal/database"
	"github.com/thejixer/memoir/internal/redis"
)

type HandlerService struct {
	dbStore    *database.PostgresStore
	redisStore *redis.RedisStore
}

func NewHandlerService(
	dbStore *database.PostgresStore,
	redisStore *redis.RedisStore,
) *HandlerService {
	return &HandlerService{
		dbStore:    dbStore,
		redisStore: redisStore,
	}
}

func (h *HandlerService) HandleHelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World from memoir")
}
