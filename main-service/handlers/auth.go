package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// placeholder map slice
// [
//
//	{ "id": "U1", "email": "admin@example.com", "password": "admin123", "name": "Admin" },
//	{ "id": "U2", "email": "user@example.com",  "password": "user123",  "name": "User"  }
//
// ]
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// placeholder
var users = []User{
	{ID: "U1", Email: "admin@example.com", Password: "admin123", Name: "Admin"},
	{ID: "U2", Email: "user@example.com", Password: "user123", Name: "User"},
}

type SigninResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func Signin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.JSON(400, map[string]string{
			"message": "email and password are required",
		})
	}

	for _, user := range users {
		if user.Email == email && user.Password == password {
			token, err := GenerateToken(user.ID, user.Email)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"message": "failed to generate token",
				})
			}
			return c.JSON(200, SigninResponse{
				AccessToken: token,
				TokenType:   "Bearer",
				ExpiresIn:   int(getJWTExpiresMinutes().Seconds()),
			})
		}
	}
	return c.JSON(401, map[string]string{
		"message": "invalid credentials",
	})
}

func Me(c echo.Context) error {
	claims := c.Get("claims").(*JWTClaims)
	return c.JSON(200, map[string]string{
		"id":    claims.Subject,
		"email": claims.Email,
	})
}
