package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoswagger "github.com/swaggo/echo-swagger"

	"main-service/routes"

	_ "main-service/docs"
)

// @title Main Service API
// @version 1.0
// @description Main service provides authentication and user management endpoints.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /main/api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter your token with "Bearer " prefix, e.g. "Bearer eyJhbGciOiJIUzI1NiI..."
func main() {
	e := echo.New()

	// Use LoggerWithConfig for verbose, custom formatting
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339} | method=${method} | uri=${uri} | status=${status} | latency=${latency_human}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Main Service /")
	})

	e.GET("/swagger/*", echoswagger.WrapHandler)

	// group route /api/v1/ from routes.go
	mainGroup := e.Group("/api/v1")
	routes.RegisterRoutes(mainGroup)

	e.Logger.Fatal(e.Start(":8081"))
}
