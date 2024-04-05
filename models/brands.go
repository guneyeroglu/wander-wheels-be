package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
)

type Brand struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllBrands(c *fiber.Ctx) error {
	db := database.ConnectDb()
	defer db.Close()

	var brands []Brand

	rows, err := db.Query(`
		SELECT 
			id, 
			name
		FROM brands
	`)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for rows.Next() {
		var brand Brand

		err := rows.Scan(&brand.Id, &brand.Name)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		brands = append(brands, brand)
	}

	return c.Status(fiber.StatusOK).JSON(brands)
}
