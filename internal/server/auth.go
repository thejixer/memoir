package server

import "github.com/labstack/echo/v4"

func (s *APIServer) ApplyAuthRoutes(e *echo.Echo) {
	auth := e.Group("/auth")
	auth.POST("/signup", s.handlerService.HandleSignup)
	auth.POST("/request-verificationCode", s.handlerService.HandleRequestVerificationEmail)
	auth.GET("/verify-email", s.handlerService.HandleEmailVerification)
	auth.POST("/login", s.handlerService.HandleLogin)
}
