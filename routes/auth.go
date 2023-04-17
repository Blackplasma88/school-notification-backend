package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type authRoutes struct {
	authController controller.AuthController
}

func NewAuthRoutes(authController controller.AuthController) Routes {
	return &authRoutes{authController: authController}
}

func (r *authRoutes) Install(app *fiber.App) {
	// app.Get("/user", r.authController.GetUserWithId)

	// app.Post("/sign-up", r.authController.SignUp)
	app.Post("/sign-in", r.authController.SignIn)
}
