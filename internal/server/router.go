package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyRoutes(e *echo.Echo) {
	e.GET("/", s.handlerService.HandleHelloWorld)
}
