package routes

import (
	"items-service/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group) {
	// GET /api/v1/items
	e.GET("/items", handlers.GetItems)

	// GET /api/v1/items/:id
	e.GET("/items/:id", handlers.GetItemById)

	// POST /api/v1/items
	e.POST("/items", handlers.CreateItem)
}
