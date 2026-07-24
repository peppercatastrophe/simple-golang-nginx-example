package routes

import (
	"main-service/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	e.POST("/auth/signin", handlers.Signin)

	protected := e.Group("")
	protected.Use(handlers.JWTMiddleware)

	protected.GET("/me", handlers.Me)

	protected.GET("/users", handlers.Users)
}
