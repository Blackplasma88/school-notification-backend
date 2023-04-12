package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type courseRoutes struct {
	courseController controller.CourseController
}

func NewCourseRoute(courseController controller.CourseController) Routes {
	return &courseRoutes{courseController: courseController}
}

func (r *courseRoutes) Install(app *fiber.App) {
	app.Get("/course/year-term", r.courseController.GetCourseByYearAndTerm)
	// app.Get("/course-id", r.courseController.GetCoursesById)
	// app.Get("/course/score", r.courseController.GetScoreById)
	// app.Get("/course/check-name", r.courseController.GetCheckNameById)

	app.Post("/course/create", r.courseController.CreateCourse)
	// app.Post("/course/update-data", r.courseController.UpdateCoursesData)
	// app.Post("/course/check-name", r.courseController.ManageCheckName)
	// app.Post("/course/score", r.courseController.ManageScore)
	app.Post("/course/change-status", r.courseController.ChangeCourseStatus)
	app.Post("/course/finish-course", r.courseController.FinishCourse)
}
