package main

import (
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

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Items Service /")
	})
	e.GET("/api/v1/items/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Items Service /api/v1/items")
	})

	e.Logger.Fatal(e.Start(":8082"))
}
