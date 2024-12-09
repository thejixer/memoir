package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyTagRoutes(e *echo.Echo) {
	g := e.Group("/tag")
	g.POST("/create", s.handlerService.HandleCreateTag, s.handlerService.Gaurd)
	g.GET("/query-note-tags", s.handlerService.HandleQueryNoteTags, s.handlerService.Gaurd)
	g.GET("/query-meeting-tags", s.handlerService.HandleQueryMeetingTags, s.handlerService.Gaurd)
}
