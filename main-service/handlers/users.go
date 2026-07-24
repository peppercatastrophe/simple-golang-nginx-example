package handlers

import "github.com/labstack/echo/v4"

func Users(c echo.Context) error {
	// placeholder
	return c.JSON(200, users)
}
