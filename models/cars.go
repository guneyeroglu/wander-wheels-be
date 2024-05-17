package models

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/guneyeroglu/wander-wheels-be/database"
	"github.com/guneyeroglu/wander-wheels-be/utils"
	"github.com/lib/pq"
)

type Images struct {
	FeaturedImage string   `json:"featuredImage"`
	OtherImages   []string `json:"otherImages"`
}

type CarModel struct {
	Id                   uuid.UUID    `json:"id"`
	Model                Model        `json:"model"`
	Color                Color        `json:"color"`
	Transmission         Transmission `json:"transmission"`
	Fuel                 Fuel         `json:"fuel"`
	Year                 int          `json:"year"`
	DailyPrice           int          `json:"dailyPrice"`
	DiscountStatus       bool         `json:"discountStatus"`
	DiscountedDailyPrice *int         `json:"discountedDailyPrice"`
	Images               Images       `json:"images"`
	Seat                 int          `json:"seat"`
	City                 City         `json:"city"`
	CreatedDate          time.Time    `json:"createdDate"`
	UpdatedDate          time.Time    `json:"updatedDate"`
}

type Car struct {
	Id  uuid.UUID `json:"id"`
	Car CarModel  `json:"car"`
}

type CarParams struct {
	CityId    *int       `json:"cityId,omitempty"`
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
}
type CarData struct {
	ModelId        *int  `json:"modelId,omitempty"`
	BrandId        *int  `json:"brandId,omitempty"`
	ColorIds       []int `json:"colorIds,omitempty"`
	TransmissionId *int  `json:"transmissionId,omitempty"`
	FuelId         *int  `json:"fuelId,omitempty"`
	MinYear        *int  `json:"minYear,omitempty"`
	MaxYear        *int  `json:"maxYear,omitempty"`
	MinPrice       *int  `json:"minPrice,omitempty"`
	MaxPrice       *int  `json:"maxPrice,omitempty"`
	Seat           *int  `json:"seat,omitempty"`
	CarParams
}

type Translations map[string]map[string]string

func GetAllCars(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var cars []Car

	query := `
		SELECT 	
			CC.id AS car_and_city_id,
			CA.id AS car_id,
			M.id AS model_id, 
			M.name AS model_name,
			B.id AS brand_id, 
			B.name AS brand_name,
			CO.id AS color_id,
			CASE
				WHEN $1 = 'tr_TR' then CO.name_tr
				ELSE CO.name_en
			END AS color_name,
			CO.code AS color_code,
			T.id AS transmission_id, 
			CASE
				WHEN $1 = 'tr_TR' then T.name_tr
				ELSE T.name_en
			END AS transmission_name,
			F.id AS fuel_id,
			CASE
				WHEN $1 = 'tr_TR' then F.name_tr
				ELSE F.name_en
			END AS transmission_name,
			CA.year AS year, 
			CA.daily_price,
			CA.featured_image, 
			CA.other_images,
			CA.seat,
			CI.id,
			CI.name,
			CA.created_date,
			CA.updated_date,
			CA.discounted_daily_price
		FROM cars_and_cities AS CC
		JOIN cars AS CA ON CA.id = CC.car_id
		JOIN models AS M ON M.id = CA.model_id
		JOIN brands AS B ON B.id = M.brand_id
		JOIN fuels AS F ON F.id = CA.fuel_id
		JOIN transmissions AS T ON T.id = CA.transmission_id
		JOIN colors AS CO ON CO.id = CA.color_id
		JOIN cities AS CI ON CI.id = CC.city_id
		WHERE 
			(M.id = $2 OR $2 IS NULL) AND
			(B.id = $3 OR $3 IS NULL) AND
			(CO.id = ANY($4) OR $4 IS NULL) AND
			(T.id = $5 OR $5 IS NULL) AND
			(F.id = $6 OR $6 IS NULL) AND
			(CA.year >= $7 OR $7 IS NULL) AND
			(CA.year <= $8 OR $8 IS NULL) AND
			(CA.daily_price >= $9 OR $9 IS NULL) AND
			(CA.daily_price <= $10 OR $10 IS NULL) AND
			(CA.seat = $11 OR $11 IS NULL) AND
			(CI.id = $12 OR $12 IS NULL) AND
			(CC.id NOT IN (
				SELECT R.car_and_city_id
				FROM rentals AS R
				WHERE (
					(R.start_date BETWEEN $13 AND $14) OR
					(R.end_date BETWEEN $13 AND $14)
				)
			))
		ORDER BY M.name
	`

	var data CarData

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	formatDate := "2006-01-02"
	formattedStartDate := data.StartDate.Format(formatDate)
	formattedEndDate := data.EndDate.Format(formatDate)

	rows, err := db.Query(
		query,                   // query
		lang,                    // $1
		data.ModelId,            // $2
		data.BrandId,            // $3
		pq.Array(data.ColorIds), // $4
		data.TransmissionId,     // $5
		data.FuelId,             // $6
		data.MinYear,            // $7
		data.MaxYear,            // $8
		data.MinPrice,           // $9
		data.MaxPrice,           // $10
		data.Seat,               // $11
		data.CityId,             // $12
		formattedStartDate,      // $13
		formattedEndDate,        // $14
	)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var carAndCity Car
		var car CarModel
		var model Model
		var color Color
		var transmission Transmission
		var fuel Fuel
		var images Images
		var city City

		err := rows.Scan(
			&carAndCity.Id,
			&car.Id,
			&model.Id, &model.Name, &model.Brand.Id, &model.Brand.Name,
			&color.Id, &color.Name, &color.Code,
			&transmission.Id, &transmission.Name,
			&fuel.Id, &fuel.Name,
			&car.Year,
			&car.DailyPrice,
			&images.FeaturedImage, pq.Array(&images.OtherImages),
			&car.Seat,
			&city.Id, &city.Name,
			&car.CreatedDate,
			&car.UpdatedDate,
			&car.DiscountedDailyPrice,
		)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		car.Model = model
		car.Color = color
		car.Transmission = transmission
		car.Fuel = fuel
		car.Images = images
		car.City = city
		if car.DiscountedDailyPrice != nil {
			car.DiscountStatus = true
		} else {
			car.DiscountStatus = false
		}
		carAndCity.Car = car
		cars = append(cars, carAndCity)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    cars,
		"message": utils.GetTranslation(lang, "success"),
	})
}

func GetCarById(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	carAndCityId := c.Params("id")
	db := database.ConnectDb()
	defer db.Close()

	var cars []Car

	query := `
		SELECT 	
			CC.id AS car_and_city_id,
			CA.id AS car_id,
			M.id AS model_id, 
			M.name AS model_name,
			B.id AS brand_id, 
			B.name AS brand_name,
			CO.id AS color_id,
			CASE
				WHEN $1 = 'tr_TR' then CO.name_tr
				ELSE CO.name_en
			END AS color_name,
			CO.code AS color_code,
			T.id AS transmission_id, 
			CASE
				WHEN $1 = 'tr_TR' then T.name_tr
				ELSE T.name_en
			END AS transmission_name,
			F.id AS fuel_id,
			CASE
				WHEN $1 = 'tr_TR' then F.name_tr
				ELSE F.name_en
			END AS transmission_name,
			CA.year, 
			CA.daily_price,
			CA.featured_image, 
			CA.other_images,
			CA.seat,
			CI.id,
			CI.name,
			CA.created_date,
			CA.updated_date,
			CA.discounted_daily_price
		FROM cars_and_cities AS CC
		JOIN cars AS CA ON CA.id = CC.car_id
		JOIN models AS M ON M.id = CA.model_id
		JOIN brands AS B ON B.id = M.brand_id
		JOIN fuels AS F ON F.id = CA.fuel_id
		JOIN transmissions AS T ON T.id = CA.transmission_id
		JOIN colors AS CO ON CO.id = CA.color_id
		JOIN cities AS CI ON CI.id = CC.city_id
		WHERE (CC.id = $2)
		ORDER BY M.name
	`

	rows, err := db.Query(
		query,        // query
		lang,         // $1
		carAndCityId, // $2
	)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var carAndCity Car
		var car CarModel
		var model Model
		var color Color
		var transmission Transmission
		var fuel Fuel
		var images Images
		var city City

		err := rows.Scan(
			&carAndCity.Id,
			&car.Id,
			&model.Id, &model.Name, &model.Brand.Id, &model.Brand.Name,
			&color.Id, &color.Name, &color.Code,
			&transmission.Id, &transmission.Name,
			&fuel.Id, &fuel.Name,
			&car.Year,
			&car.DailyPrice,
			&images.FeaturedImage, pq.Array(&images.OtherImages),
			&car.Seat,
			&city.Id, &city.Name,
			&car.CreatedDate,
			&car.UpdatedDate,
			&car.DiscountedDailyPrice,
		)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		car.Model = model
		car.Color = color
		car.Transmission = transmission
		car.Fuel = fuel
		car.Images = images
		car.City = city
		if car.DiscountedDailyPrice != nil {
			car.DiscountStatus = true
		} else {
			car.DiscountStatus = false
		}
		carAndCity.Car = car

		cars = append(cars, carAndCity)
	}

	if cars == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"data":    nil,
			"message": utils.GetTranslation(lang, "notFound"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    cars[0],
		"message": utils.GetTranslation(lang, "success"),
	})
}

type Prices struct {
	MinPrice int `json:"minPrice"`
	MaxPrice int `json:"maxPrice"`
}

func GetPriceRange(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var prices Prices

	query := `
		SELECT
			MIN(C.daily_price) AS min_price,
			MAX(C.daily_price) AS max_price
		FROM cars_and_cities AS CC
		JOIN cars AS C ON C.id = CC.car_id
	`

	rows, err := db.Query(query)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		err := rows.Scan(
			&prices.MinPrice,
			&prices.MaxPrice,
		)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    prices,
		"message": utils.GetTranslation(lang, "success"),
	})
}

type Years struct {
	MinYear int `json:"minYear"`
	MaxYear int `json:"maxYear"`
}

func GetYearRange(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var years Years

	query := `
		SELECT
			MIN(C.year) AS min_year,
			MAX(C.year) AS max_year
		FROM cars_and_cities AS CC
		JOIN cars AS C ON C.id = CC.car_id
	`

	rows, err := db.Query(query)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		err := rows.Scan(
			&years.MinYear,
			&years.MaxYear,
		)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    years,
		"message": utils.GetTranslation(lang, "success"),
	})
}

func GetSeats(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var seats []int

	query := `
			SELECT DISTINCT
			C.seat
		FROM cars_and_cities AS CC
		JOIN cars AS C ON C.id = CC.car_id
		ORDER BY c.seat
	`

	rows, err := db.Query(query)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var seat int

		err := rows.Scan(
			&seat,
		)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		seats = append(seats, seat)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data": fiber.Map{
			"seats": seats,
		},
		"message": utils.GetTranslation(lang, "success"),
	})
}

type RentCar struct {
	UserId       uuid.UUID `json:"userId"`
	CarAndCityId uuid.UUID `json:"carAndCityId"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
}

type RentedCar struct {
	Id uuid.UUID `json:"id"`
}

func CreateRental(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	query := `
		INSERT INTO rentals (
			user_id,
			car_and_city_id,
			start_date,
			end_date
		) VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`
	queryForFindingCarAndCity := `
		SELECT
			car_and_city_id
		FROM rentals
		WHERE 
			(car_and_city_id = $1) AND
			((start_date BETWEEN $2 AND $3) OR
			(end_date BETWEEN $2 AND $3))
		ORDER BY id
	`

	var rent RentCar

	if err := c.BodyParser(&rent); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	formatDate := "2006-01-02"
	formattedStartDate := rent.StartDate.Format(formatDate)
	formattedEndDate := rent.EndDate.Format(formatDate)

	rowsForFindingCarAndCity, errForFindingCarAndCity := db.Query(
		queryForFindingCarAndCity, // query
		rent.CarAndCityId,         // $1
		formattedStartDate,        // $2
		formattedEndDate,          // $3

	)

	if errForFindingCarAndCity != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": errForFindingCarAndCity.Error(),
		})
	}

	var cars []RentedCar

	for rowsForFindingCarAndCity.Next() {
		var car RentedCar
		err := rowsForFindingCarAndCity.Scan(
			&car.Id,
		)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		cars = append(cars, car)
	}

	if cars != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  fiber.StatusConflict,
			"message": utils.GetTranslation(lang, "carAlreadyRented"),
		})
	}

	tokenType := os.Getenv("JWT_TOKEN_TYPE")
	jwtSecretCode := os.Getenv("JWT_SECRET_CODE")
	tokenString := strings.Join(strings.Split(c.Get("Authorization"), tokenType), "")
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretCode), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	_, err = db.Query(
		query,              // query
		rent.UserId,        // $1
		rent.CarAndCityId,  // $2
		formattedStartDate, // $3
		formattedEndDate,   // $4
	)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": utils.GetTranslation(lang, "rentalSuccess"),
	})
}
