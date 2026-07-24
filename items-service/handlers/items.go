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

var items = []Item{
	{ID: "I1", Name: "Pencil", Price: 2000, OwnerID: "U1"},
	{ID: "I2", Name: "Notebook", Price: 12000, OwnerID: "U2"},
}

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

func CreateItem(c echo.Context) error {
	token := extractToken(c)
	if token == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "missing authorization header"})
	}

	// get current user from main-service
	me, err := fetchMe(token)
	if err != nil {
		return authError(c, err)
	}

	name := c.FormValue("name")
	if name == "" {
		return c.JSON(400, echo.Map{"message": "name is required"})
	}
	price, _ := strconv.Atoi(c.FormValue("price"))
	if price == 0 {
		return c.JSON(400, echo.Map{"message": "price is required"})
	}

	item := Item{
		ID:      "I" + strconv.Itoa(len(items)+1),
		Name:    name,
		Price:   price,
		OwnerID: me.ID,
	}
	items = append(items, item)

	// fetch users to resolve owner details
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
