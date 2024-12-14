package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyMeetingRoutes(e *echo.Echo) {
	g := e.Group("/meeting")
	g.POST("/create", s.handlerService.HandleCreateMeeting, s.handlerService.Gaurd)
}
