package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type locationRoutes struct {
	locationController controller.LocationController
}

func NewLocationRoute(locationController controller.LocationController) Routes {
	return &locationRoutes{locationController: locationController}
}

func (r *locationRoutes) Install(app *fiber.App) {
	app.Get("/location/all", r.locationController.GetLocationAll)
	// app.Get("/location/id", r.locationController.GetLocationById)

	app.Post("/location/create", r.locationController.CreateLocation)
	// app.Post("/location/update", r.locationController.UpdateLocationData)
}
