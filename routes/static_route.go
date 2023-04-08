package routes

import (
	"github.com/gofiber/fiber/v2"
)

type staticRoutes struct{}

func NewStaticRoutes() Routes {
	return &staticRoutes{}
}

func (sr *staticRoutes) Install(app *fiber.App) {
	app.Static("/files/information", "./storage/information")
	app.Static("/files/profile", "./storage/profile")
}
