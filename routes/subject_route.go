package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type subjectRoutes struct {
	subjectController controller.SubjectController
}

func NewSubjectRoute(subjectController controller.SubjectController) Routes {
	return &subjectRoutes{subjectController: subjectController}
}

func (r *subjectRoutes) Install(app *fiber.App) {
	app.Get("/subject/all", r.subjectController.GetSubjectAll)
	// app.Get("/subject/id", r.subjectController.GetSubjectById)

	app.Post("/subject/create", r.subjectController.CreateSubject)
	app.Post("/subject/add-instructor", r.subjectController.AddInstructor)
	// app.Post("/subject/update", r.subjectController.UpdateSubject)
}
