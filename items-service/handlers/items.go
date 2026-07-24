package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// placeholder items slice
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
	return c.JSON(200, items)

}

func GetItemById(c echo.Context) error {
	id := c.Param("id")
	for _, item := range items {
		if item.ID == id {
			return c.JSON(200, item)
		}
	}
	return c.JSON(404, echo.Map{"message": "item not found"})
}

func CreateItem(c echo.Context) error {
	// request: { name, price }
	name := c.FormValue("name")
	if name == "" {
		return c.JSON(400, echo.Map{"message": "name is required"})
	}
	price, _ := strconv.Atoi(c.FormValue("price"))
	if price == 0 {
		return c.JSON(400, echo.Map{"message": "price is required"})
	}

	item := new(Item)
	// if err := c.Bind(item); err != nil {
	// 	return err
	// }
	item.Name = name
	item.Price = price
	item.ID = "I" + strconv.Itoa(len(items)+1)
	item.OwnerID = "U1" // placeholder
	items = append(items, *item)
	return c.JSON(201, item)
}
