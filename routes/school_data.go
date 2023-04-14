package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type schoolDataRoutes struct {
	schoolDataController controller.SchoolDataController
}

func NewSchoolDataRoute(schoolDataController controller.SchoolDataController) Routes {
	return &schoolDataRoutes{schoolDataController: schoolDataController}
}

func (r *schoolDataRoutes) Install(app *fiber.App) {
	// app.Get("/school-data/all", r.schoolDataController.)
	app.Get("/school-data/subject-category", r.schoolDataController.GetSubjectCategory)
	// app.Get("/school-data/term-year-data", r.schoolDataController.GetSubjectCategory)
	// app.Get("/school-data/id", r.schoolDataController.)

	app.Post("/school-data/add-year-term", r.schoolDataController.AddYearAndTerm)
	app.Post("/school-data/add-subject-category", r.schoolDataController.AddSubjectCategory)
	app.Post("/school-data/end-term", r.schoolDataController.EndTerm)
	// app.Post("/school-data/update", r.schoolDataController.)
}
