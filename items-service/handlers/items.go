package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Item struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	OwnerID string `json:"owner_id"`
}

type CreateItemRequest struct {
	Name  string `json:"name" form:"name" example:"Eraser"`
	Price int    `json:"price" form:"price" example:"3000"`
}

var items = []Item{
	{ID: "I1", Name: "Pencil", Price: 2000, OwnerID: "U1"},
	{ID: "I2", Name: "Notebook", Price: 12000, OwnerID: "U2"},
}

// GetItems
//
//	@Summary		List items
//	@Description	Returns all items with owner details. Calls main-service /me and /users to resolve user data.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200	{object}	ItemsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		502	{object}	ErrorResponse
//	@Router			/items [get]
func GetItems(c echo.Context) error {
	token := extractToken(c)
	if token == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "missing authorization header"})
	}

	me, err := fetchMe(token)
	if err != nil {
		return authError(c, err)
	}

	users, err := fetchUsers(token)
	if err != nil {
		return authError(c, err)
	}
	ownerMap := buildOwnerMap(users)

	itemResponses := make([]ItemResponse, 0, len(items))
	for _, item := range items {
		itemResponses = append(itemResponses, ItemResponse{
			ID:      item.ID,
			Name:    item.Name,
			Price:   item.Price,
			OwnerID: item.OwnerID,
			Owner:   lookupOwner(ownerMap, item.OwnerID),
		})
	}

	return c.JSON(200, ItemsResponse{
		RequestedBy: *me,
		Items:       itemResponses,
	})
}

// GetItemById
//
//	@Summary		Get item by ID
//	@Description	Returns a single item with owner details
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Item ID"
//	@Security		Bearer
//	@Success		200	{object}	ItemResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		502	{object}	ErrorResponse
//	@Router			/items/{id} [get]
func GetItemById(c echo.Context) error {
	token := extractToken(c)
	if token == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "missing authorization header"})
	}

	id := c.Param("id")
	var found *Item
	for _, item := range items {
		if item.ID == id {
			found = &item
			break
		}
	}
	if found == nil {
		return c.JSON(404, echo.Map{"message": "item not found"})
	}

	users, err := fetchUsers(token)
	if err != nil {
		return authError(c, err)
	}
	ownerMap := buildOwnerMap(users)

	return c.JSON(200, ItemResponse{
		ID:      found.ID,
		Name:    found.Name,
		Price:   found.Price,
		OwnerID: found.OwnerID,
		Owner:   lookupOwner(ownerMap, found.OwnerID),
	})
}

// CreateItem
//
//	@Summary		Create item
//	@Description	Creates a new item. Owner is derived from the authenticated user via main-service /me.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			body	body	CreateItemRequest	true	"Item data"
//	@Security		Bearer
//	@Success		201	{object}	ItemResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		502	{object}	ErrorResponse
//	@Router			/items [post]
func CreateItem(c echo.Context) error {
	token := extractToken(c)
	if token == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "missing authorization header"})
	}

	var req CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, echo.Map{"message": "invalid request body"})
	}

	if req.Name == "" {
		return c.JSON(400, echo.Map{"message": "name is required"})
	}

	me, err := fetchMe(token)
	if err != nil {
		return authError(c, err)
	}

	item := Item{
		ID:      "I" + strconv.Itoa(len(items)+1),
		Name:    req.Name,
		Price:   req.Price,
		OwnerID: me.ID,
	}
	items = append(items, item)

	users, err := fetchUsers(token)
	if err != nil {
		return authError(c, err)
	}
	ownerMap := buildOwnerMap(users)

	return c.JSON(201, ItemResponse{
		ID:      item.ID,
		Name:    item.Name,
		Price:   item.Price,
		OwnerID: item.OwnerID,
		Owner:   lookupOwner(ownerMap, item.OwnerID),
	})
}
