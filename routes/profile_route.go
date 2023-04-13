package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type profileRoutes struct {
	profileController controller.ProfileController
}

func NewProfileRoute(profileController controller.ProfileController) Routes {
	return &profileRoutes{profileController: profileController}
}

func (r *profileRoutes) Install(app *fiber.App) {
	app.Get("/profile/all", r.profileController.GetProfileAll)
	app.Get("/profile/id", r.profileController.GetProfileById)

	app.Post("/profile/create", r.profileController.CreateNewProfile)
	// app.Post("/profile/update", r.profileController.UpdateProfile)

	// app.Post("/profile/create-admin", r.profileController.CreateAdmin)
}
