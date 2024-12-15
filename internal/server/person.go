package server

import (
	"github.com/labstack/echo/v4"
)

func (s *APIServer) ApplyPersonRoutes(e *echo.Echo) {
	g := e.Group("/person")
	g.POST("/create", s.handlerService.HandleCreatePerson, s.handlerService.Gaurd)
	g.GET("/query", s.handlerService.HandleQueryMyPersons, s.handlerService.Gaurd)
	g.GET("/s/:id", s.handlerService.HandleGetSinglePerson, s.handlerService.Gaurd)
	g.GET("/byMeeting/:meetingId", s.handlerService.HandleGetPersonsByMeetingId, s.handlerService.Gaurd)
}
