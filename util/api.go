package util

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ResponseNotSuccess(c *fiber.Ctx, code int, errMsg string) error {
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": errMsg,
	})
}

func ResponseSuccess(c *fiber.Ctx, code int, msg string, data interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

func CheckStringData(data string, name string) (string, error) {
	data = strings.TrimSpace(data)
	if len(data) == 0 {
		return "", ReturnError(ErrRequireParameter.Error() + name)
	}

	return data, nil
}

func CheckIntegerData(data *int, name string) (int, error) {

	if data == nil {
		return 0, ReturnError(ErrRequireParameter.Error() + name)
	}

	if *data <= 0 {
		return 0, ReturnError(name + ErrValueInvalid.Error())
	}

	return *data, nil
}
