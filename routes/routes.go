package routes

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

var jwtSecretKey = []byte("")
var jwtSingingMethod = jwt.SigningMethodHS256.Name

type Routes interface {
	Install(app *fiber.App)
}

func authenticath(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:    jwtSecretKey,
		SigningMethod: jwtSingingMethod,
		TokenLookup:   "header:Authorization",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(map[string]interface{}{
				"error": err.Error(),
			})
		},
	})(c)
}
