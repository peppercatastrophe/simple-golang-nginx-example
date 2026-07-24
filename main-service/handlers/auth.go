package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SigninRequest struct {
	Email    string `json:"email" form:"email" example:"admin@example.com"`
	Password string `json:"password" form:"password" example:"admin123"`
}

type SigninResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"3600"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// placeholder
var users = []User{
	{ID: "U1", Email: "admin@example.com", Password: "admin123", Name: "Admin"},
	{ID: "U2", Email: "user@example.com", Password: "user123", Name: "User"},
}

// Signin
//
//	@Summary		Sign in
//	@Description	Authenticate with email and password to receive a JWT Bearer token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	SigninRequest	true	"Credentials"
//	@Success		200	{object}	SigninResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Router			/auth/signin [post]
func Signin(c echo.Context) error {
	var req SigninRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, ErrorResponse{Message: "invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(400, ErrorResponse{Message: "email and password are required"})
	}

	for _, user := range users {
		if user.Email == req.Email && user.Password == req.Password {
			token, err := GenerateToken(user.ID, user.Email)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to generate token"})
			}
			return c.JSON(200, SigninResponse{
				AccessToken: token,
				TokenType:   "Bearer",
				ExpiresIn:   int(getJWTExpiresMinutes().Seconds()),
			})
		}
	}
	return c.JSON(401, ErrorResponse{Message: "invalid credentials"})
}

// Me
//
//	@Summary		Get current user
//	@Description	Returns the authenticated user's profile based on the JWT token
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{object}	UserResponse
//	@Failure		401	{object}	ErrorResponse
//	@Router			/me [get]
func Me(c echo.Context) error {
	claims := c.Get("claims").(*JWTClaims)
	for _, user := range users {
		if user.ID == claims.Subject {
			return c.JSON(200, UserResponse{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
			})
		}
	}
	return c.JSON(200, UserResponse{
		ID:    claims.Subject,
		Email: claims.Email,
	})
}
