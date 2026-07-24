package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type RequestedBy struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Owner struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ItemResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	OwnerID string `json:"owner_id"`
	Owner   Owner  `json:"owner"`
}

type ItemsResponse struct {
	RequestedBy RequestedBy    `json:"requested_by"`
	Items       []ItemResponse `json:"items"`
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

func getMainServiceURL() string {
	if v := os.Getenv("MAIN_SERVICE_URL"); v != "" {
		return v
	}
	return "http://main-service:8081"
}

func extractToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

func fetchMe(token string) (*RequestedBy, error) {
	req, err := http.NewRequest("GET", getMainServiceURL()+"/api/v1/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("main-service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid or expired token")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("main-service /me returned status %d", resp.StatusCode)
	}

	var me RequestedBy
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
		return nil, fmt.Errorf("failed to decode /me response: %w", err)
	}
	return &me, nil
}

func fetchUsers(token string) ([]Owner, error) {
	req, err := http.NewRequest("GET", getMainServiceURL()+"/api/v1/users", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("main-service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: invalid or expired token")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("main-service /users returned status %d", resp.StatusCode)
	}

	var users []Owner
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode /users response: %w", err)
	}
	return users, nil
}

func buildOwnerMap(users []Owner) map[string]Owner {
	m := make(map[string]Owner, len(users))
	for _, u := range users {
		m[u.ID] = u
	}
	return m
}

func lookupOwner(ownerMap map[string]Owner, ownerID string) Owner {
	if owner, ok := ownerMap[ownerID]; ok {
		return owner
	}
	return Owner{}
}

func authError(c echo.Context, err error) error {
	msg := err.Error()
	if msg == "unauthorized: invalid or expired token" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": msg})
	}
	return c.JSON(http.StatusBadGateway, echo.Map{"message": "main-service unavailable"})
}
