package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
	"github.com/guneyeroglu/wander-wheels-be/utils"
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
			CASE 
				WHEN $1 = 'tr_TR' then name_tr
				ELSE name_en
			END as name
		FROM transmissions
		ORDER BY id ASC
	`

	rows, err := db.Query(query, lang)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var transmission Transmission

		err := rows.Scan(&transmission.Id, &transmission.Name)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		transmissions = append(transmissions, transmission)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    transmissions,
		"message": utils.GetTranslation(lang, "success"),
	})
}
