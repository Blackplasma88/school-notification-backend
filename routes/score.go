package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type scoreRoutes struct {
	scoreController controller.ScoreController
}

func NewScoreRoute(scoreController controller.ScoreController) Routes {
	return &scoreRoutes{scoreController: scoreController}
}

func (r *scoreRoutes) Install(app *fiber.App) {
	app.Get("/score/score-in-course", r.scoreController.GetScoreByCourseId)
	app.Get("/score/score-data", r.scoreController.GetScoreDataByCourseIdAndNameSore)

	app.Post("/score/create", r.scoreController.CreateScore)
	app.Post("/score/add-student-score", r.scoreController.AddStudentScore)
	app.Post("/score/update-student-score", r.scoreController.UpdateStudentScore)

}
