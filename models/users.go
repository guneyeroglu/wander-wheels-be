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
)

type Role struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Mail     string    `json:"mail"`
	Password string    `json:"password"`
	Role     Role      `json:"role"`
}

type LoginData struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Role     Role      `json:"role"`
}

func Login(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var data LoginData

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	var userData []LoginData

	query := `
		SELECT
			U.id,
			U.username,
			U.password,
			R.id,
			R.name
		FROM users AS U
		JOIN roles as R ON r.id = u.role_id
		WHERE 
			(U.username = $1) AND
			(U.password = $2)
		ORDER BY U.username ASC
	`

	rows, err := db.Query(
		query,         //query
		data.Username, // $1
		data.Password, // $2
	)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var user LoginData
		var role Role

		err := rows.Scan(&user.Id, &user.Username, &user.Password, &role.Id, &role.Name)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		user.Role = role

		userData = append(userData, user)
	}

	if userData == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"data":    nil,
			"message": utils.GetTranslation(lang, "notFound"),
		})
	}

	claims := jwt.MapClaims{
		"id":       userData[0].Id,
		"username": userData[0].Username,
		"roleName": userData[0].Role.Name,
		"exp":      time.Now().Add(time.Hour * 24 * 1).Unix(), // 1 day
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecretCode := os.Getenv("JWT_SECRET_CODE")

	t, err := token.SignedString([]byte(jwtSecretCode))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	loginResponse := map[string]interface{}{
		"username": userData[0].Username,
		"token":    t,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    loginResponse,
		"message": utils.GetTranslation(lang, "success"),
	})
}

func GetUserInfo(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	tokenType := os.Getenv("JWT_TOKEN_TYPE")
	jwtSecretCode := os.Getenv("JWT_SECRET_CODE")
	tokenString := strings.Join(strings.Split(c.Get("Authorization"), tokenType), "")
	claims := jwt.MapClaims{}

	parseToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretCode), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	id := parseToken.Claims.(jwt.MapClaims)["id"].(string)

	var users []User

	query := `
		SELECT 
			U.id, 
			U.username,
			U.mail, 
			R.id,
			R.name 
		FROM users AS U
		JOIN roles AS R ON r.id = u.role_id
		WHERE 
			(U.id = $1)
		ORDER BY U.username ASC
	`

	rows, err := db.Query(
		query, // query
		id,    // $1
	)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var user User
		var role Role

		err := rows.Scan(&user.Id, &user.Username, &user.Mail, &role.Id, &role.Name)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		user.Role = role

		users = append(users, user)
	}

	if users == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"data":    nil,
			"message": utils.GetTranslation(lang, "notFound"),
		})
	}

	userInfo := jwt.MapClaims{
		"id":       users[0].Id,
		"username": users[0].Username,
		"mail":     users[0].Mail,
		"role":     users[0].Role,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    userInfo,
		"message": utils.GetTranslation(lang, "success"),
	})
}

type SignUpData struct {
	Username *string `json:"username"`
	Mail     *string `json:"mail"`
	Password *string `json:"password"`
}

func SignUp(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	roleId := 2 // customer
	db := database.ConnectDb()
	defer db.Close()

	queryForFindingUser := `
		SELECT 
			U.id, 
			U.username,
			U.mail, 
			U.password,
			R.id,
			R.name 
		FROM users AS U
		JOIN roles AS R ON r.id = u.role_id
		WHERE 
			(U.username = $1) AND
			(U.password = $2)
		ORDER BY U.username ASC
	`

	var data SignUpData

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	var users []User

	rowsForFindingUser, errForFindingUser := db.Query(
		queryForFindingUser, // query
		data.Username,       // $1
		data.Password,       // $2
	)

	if errForFindingUser != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": errForFindingUser.Error(),
		})
	}

	for rowsForFindingUser.Next() {
		var user User
		var role Role

		err := rowsForFindingUser.Scan(&user.Id, &user.Username, &user.Mail, &user.Password, &role.Id, &role.Name)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		user.Role = role

		users = append(users, user)
	}

	if users != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  fiber.StatusConflict,
			"message": utils.GetTranslation(lang, "userAlreadyExists"),
		})
	}

	query := `
		INSERT INTO users (
		username,
		mail, 
		password,
		role_id
	)
		VALUES (
			$1, 
			$2, 
			$3,
			$4
		);
	`

	_, err := db.Query(
		query,         // query
		data.Username, // $1
		data.Mail,     // $2
		data.Password, // $3
		roleId,        // $4
	)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": utils.GetTranslation(lang, "newUserCreated"),
	})
}
