package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
)

type Fuel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllFuels(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var fuels []Fuel
	query := `
			SELECT 
			id, 
			name_en as name
			FROM fuels
		`

	if lang == "tr_TR" {
		query = `
				SELECT 
				id, 
				name_tr as name
				FROM fuels
			`
	}

	rows, err := db.Query(query)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for rows.Next() {
		var fuel Fuel

		err := rows.Scan(&fuel.Id, &fuel.Name)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		fuels = append(fuels, fuel)
	}

	return c.Status(fiber.StatusOK).JSON(fuels)
}
