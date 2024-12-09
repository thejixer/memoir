package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyRoutes(e *echo.Echo) {
	e.GET("/", s.handlerService.HandleHelloWorld)

	s.ApplyAuthRoutes(e)
	s.ApplyPersonRoutes(e)
	s.ApplyTagRoutes(e)
	s.ApplyNoteRoutes(e)

}
