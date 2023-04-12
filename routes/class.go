package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type classRoutes struct {
	classController controller.ClassController
}

func NewClassRoute(classController controller.ClassController) Routes {
	return &classRoutes{classController: classController}
}

func (r *classRoutes) Install(app *fiber.App) {
	app.Get("/class/all", r.classController.GetClassAllByClassYear)
	// app.Get("/class-id", r.classController.GetClassById)

	app.Post("/class/create", r.classController.CreateClass)
	app.Post("/class/set-advisor", r.classController.SetAdvisor)
	// app.Post("/class/update", r.classController.UpdateClassData)
}
