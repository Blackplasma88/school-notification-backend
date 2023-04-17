package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type courseSummaryRoutes struct {
	courseSummaryController controller.CourseSummaryController
}

func NewCourseSummaryRoute(courseSummaryController controller.CourseSummaryController) Routes {
	return &courseSummaryRoutes{courseSummaryController: courseSummaryController}
}

func (r *courseSummaryRoutes) Install(app *fiber.App) {
	app.Get("/course-summary", r.courseSummaryController.GetSummaryCourse)
	app.Get("/course-summary/student", r.courseSummaryController.StudentGetSummaryCourse)

	app.Post("/summary/course-id", r.courseSummaryController.SummaryCourse)
}
