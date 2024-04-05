package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
)

type Transmission struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllTransmissions(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var transmissions []Transmission
	query := `
			SELECT 
			id, 
			name_en as name
			FROM transmissions
		`

	if lang == "tr_TR" {
		query = `
				SELECT 
				id, 
				name_tr as name
				FROM transmissions
			`
	}

	rows, err := db.Query(query)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for rows.Next() {
		var transmission Transmission

		err := rows.Scan(&transmission.Id, &transmission.Name)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		transmissions = append(transmissions, transmission)
	}

	return c.Status(fiber.StatusOK).JSON(transmissions)
}
