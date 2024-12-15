package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyNoteRoutes(e *echo.Echo) {
	g := e.Group("/note")
	g.POST("/create-person-note", s.handlerService.HandleCreatePersonNote, s.handlerService.Gaurd)
	g.GET("/byPerson/:personId", s.handlerService.HandleGetNotesByPersonId, s.handlerService.Gaurd)
	g.POST("/create-meeting-note", s.handlerService.HandleCreateMeetingNote, s.handlerService.Gaurd)
	g.GET("/byMeeting/:meetingId", s.handlerService.HandleGetNotesByMeetingId, s.handlerService.Gaurd)
}
