package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thejixer/memoir/internal/database"
)

type HandlerService struct {
	dbStore *database.PostgresStore
}

func NewHandlerService(dbStore *database.PostgresStore) *HandlerService {
	return &HandlerService{
		dbStore: dbStore,
	}
}

func (h *HandlerService) HandleHelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World from memoir")
}
