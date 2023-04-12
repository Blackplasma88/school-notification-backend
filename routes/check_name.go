package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type checkNameRoutes struct {
	checkNameController controller.CheckNameController
}

func NewCheckNameRoute(checkNameController controller.CheckNameController) Routes {
	return &checkNameRoutes{checkNameController: checkNameController}
}

func (r *checkNameRoutes) Install(app *fiber.App) {
	app.Get("/check-name/check-name-in-course", r.checkNameController.GetDateByCourseId)
	app.Get("/check-name/check-name-data", r.checkNameController.GetCheckNameDataByCourseIdAndDate)
	// app.Get("/CheckName/CheckName-data", r.CheckNameController.GetCheckNameDataByCourseIdAndNameSore)

	app.Post("/check-name/add-date", r.checkNameController.AddDateForCheck)
	app.Post("/check-name/student-check", r.checkNameController.CheckNameStudent)
	app.Post("/check-name/end-date", r.checkNameController.EndDateCheckName)
	// app.Post("/CheckName/add-student-CheckName", r.CheckNameController.AddStudentCheckName)
	// app.Post("/CheckName/update-student-CheckName", r.CheckNameController.UpdateStudentCheckName)

}
