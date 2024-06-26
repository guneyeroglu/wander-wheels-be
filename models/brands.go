package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
	"github.com/guneyeroglu/wander-wheels-be/utils"
)

type Brand struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllBrands(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var brands []Brand

	query := `
		SELECT 
			id, 
			name
		FROM brands
		ORDER BY id ASC
	`

	rows, err := db.Query(query)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var brand Brand

		err := rows.Scan(&brand.Id, &brand.Name)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		brands = append(brands, brand)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    brands,
		"message": utils.GetTranslation(lang, "success"),
	})
}
