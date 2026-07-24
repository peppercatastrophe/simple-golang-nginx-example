package routes

import (
	"items-service/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group) {
	e.GET("/items", handlers.GetItems)

	e.GET("/items/:id", handlers.GetItemById)

	e.POST("/items", handlers.CreateItem)
}
