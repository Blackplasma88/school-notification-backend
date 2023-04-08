package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type newsRoutes struct {
	informationController controller.InformationController
}

func NewInformationRoutes(informationController controller.InformationController) Routes {
	return &newsRoutes{informationController: informationController}
}

func (r *newsRoutes) Install(app *fiber.App) {
	app.Get("/information/all", r.informationController.GetInformationAll)
	app.Get("/information/id", r.informationController.GetInformationById)

	app.Post("/information/create", r.informationController.CreateInformation)
	app.Post("/information/update", r.informationController.UpdateInformation)
}
