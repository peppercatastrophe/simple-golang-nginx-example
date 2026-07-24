package handlers

import "github.com/labstack/echo/v4"

// Users
//
//	@Summary		List users
//	@Description	Returns all registered users
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{array}	UserResponse
//	@Failure		401	{object}	ErrorResponse
//	@Router			/users [get]
func Users(c echo.Context) error {
	userResponses := make([]UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = UserResponse{
			ID:    u.ID,
			Email: u.Email,
			Name:  u.Name,
		}
	}
	return c.JSON(200, userResponses)
}
