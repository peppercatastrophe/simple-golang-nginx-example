package routes

import (
	"main-service/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// POST /api/v1/auth/signin
	// Signin godoc
	//
	//	@Summary		Signin
	//	@Description	Signin
	//	@Tags			accounts
	//	@Accept			json
	//	@Produce		json
	//	@Success		200	{object}	model.User
	//	@Failure		400	{object}	httputil.HTTPError
	//	@Failure		401	{object}	httputil.HTTPError
	//	@Failure		404	{object}	httputil.HTTPError
	//	@Failure		500	{object}	httputil.HTTPError
	//	@Security		ApiKeyAuth
	//	@Router			/auth/signin [post]
	e.POST("/auth/signin", handlers.Signin)

	// TODO: response 401 { "message": "invalid credentials" } for endpoints below

	// GET /api/v1/me
	e.GET("/me", handlers.Me)

	// GET /api/v1/users
	e.GET("/users", handlers.Users)
}
