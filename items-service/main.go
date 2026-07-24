package main

import (
	"items-service/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Use LoggerWithConfig for verbose, custom formatting
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339} | method=${method} | uri=${uri} | status=${status} | latency=${latency_human}\n",
	}))

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Items Service /")
	})
	itemsGroup := e.Group("/api/v1")
	routes.RegisterRoutes(itemsGroup)

	e.Logger.Fatal(e.Start(":8082"))
}
