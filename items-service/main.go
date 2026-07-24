package main

import (
	"items-service/routes"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoswagger "github.com/swaggo/echo-swagger"

	_ "items-service/docs"
)

// @title Items Service API
// @version 1.0
// @description Items service provides CRUD operations for items. It communicates with main-service via REST to authenticate and resolve user data.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /items/api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter your token with "Bearer " prefix, e.g. "Bearer eyJhbGciOiJIUzI1NiI..."
func main() {
	// Load .env file (ignore error if not present in production)
	_ = godotenv.Load()

	e := echo.New()

	// Use LoggerWithConfig for verbose, custom formatting
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339} | method=${method} | uri=${uri} | status=${status} | latency=${latency_human}\n",
	}))

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Items Service /")
	})

	e.GET("/swagger/*", echoswagger.WrapHandler)

	itemsGroup := e.Group("/api/v1")
	routes.RegisterRoutes(itemsGroup)

	e.Logger.Fatal(e.Start(":8082"))
}
