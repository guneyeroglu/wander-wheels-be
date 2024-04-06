package models

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guneyeroglu/wander-wheels-be/database"
	"github.com/lib/pq"
)

type Images struct {
	FeaturedImage string   `json:"featuredImage"`
	OtherImages   []string `json:"otherImages"`
}

type Car struct {
	Id           uuid.UUID    `json:"id"`
	Model        Model        `json:"model"`
	Color        Color        `json:"color"`
	Transmission Transmission `json:"transmission"`
	Fuel         Fuel         `json:"fuel"`
	Year         int          `json:"year"`
	DailyPrice   int          `json:"dailyPrice"`
	Images       Images       `json:"images"`
	Seat         int          `json:"seat"`
	CreatedDate  time.Time    `json:"createdDate"`
	UpdatedDate  time.Time    `json:"updatedDate"`
}

func GetAllCars(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var cars []Car

	query := `
		SELECT 	
			CA.id AS car_id, 
			M.id AS model_id, 
			M.name AS model_name,
			B.id AS brand_id, 
			B.name AS brand_name,
			CO.id AS color_id,
			CASE
				WHEN $1 = 'tr_TR' then CO.name_tr
				ELSE CO.name_en
			END as color_name,
			T.id AS transmission_id, 
			CASE
				WHEN $1 = 'tr_TR' then T.name_tr
				ELSE T.name_en
			END as transmission_name,
			F.id AS fuel_id,
			CASE
				WHEN $1 = 'tr_TR' then F.name_tr
				ELSE F.name_en
			END as transmission_name,
			CA.year, 
			CA.daily_price,
			CA.featured_image, 
			CA.other_images,
			CA.seat,
			CA.created_date,
			CA.updated_date
		FROM cars AS CA
		JOIN models AS M ON M.id = CA.model_id
		JOIN brands AS B ON B.id = M.brand_id
		JOIN fuels AS F ON F.id = CA.fuel_id
		JOIN transmissions AS T ON T.id = CA.transmission_id
		JOIN colors AS CO ON CO.id = CA.color_id
		ORDER BY M.name
`

	rows, err := db.Query(query, lang)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for rows.Next() {
		var car Car
		var model Model
		var color Color
		var transmission Transmission
		var fuel Fuel
		var images Images

		err := rows.Scan(
			&car.Id,
			&model.Id, &model.Name, &model.Brand.Id, &model.Brand.Name,
			&color.Id, &color.Name,
			&transmission.Id, &transmission.Name,
			&fuel.Id, &fuel.Name,
			&car.Year,
			&car.DailyPrice,
			&images.FeaturedImage, pq.Array(&images.OtherImages),
			&car.Seat,
			&car.CreatedDate,
			&car.UpdatedDate,
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		car.Model = model
		car.Color = color
		car.Transmission = transmission
		car.Fuel = fuel
		car.Images = images

		cars = append(cars, car)
	}

	return c.Status(fiber.StatusOK).JSON(cars)
}
