package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"main-service/routes"
)

func main() {
	e := echo.New()

	// Use LoggerWithConfig for verbose, custom formatting
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339} | method=${method} | uri=${uri} | status=${status} | latency=${latency_human}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Main Service /")
	})
	// group route /api/v1/ from routes.go
	mainGroup := e.Group("/api/v1")
	routes.RegisterRoutes(mainGroup)

	e.Logger.Fatal(e.Start(":8081"))
}
