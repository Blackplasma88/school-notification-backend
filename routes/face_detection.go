package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type faceDetectionRoutes struct {
	faceDetectionController controller.FaceDetectionController
}

func NewFaceDetectionRoute(faceDetectionController controller.FaceDetectionController) Routes {
	return &faceDetectionRoutes{faceDetectionController: faceDetectionController}
}

func (r *faceDetectionRoutes) Install(app *fiber.App) {
	app.Get("/face-detection/open-camera", r.faceDetectionController.OpenCamera)
	app.Get("/face-detection/id", r.faceDetectionController.GetById)
	app.Get("/face-detection/all", r.faceDetectionController.GetAll)
	app.Get("/face-detection/class_id", r.faceDetectionController.GetByClassId)

	app.Post("/face-detection/create-data", r.faceDetectionController.CreatFaceDetectionData)
	app.Post("/face-detection/upload-image-data", r.faceDetectionController.UploadImageData)
	app.Post("/face-detection/trained-model", r.faceDetectionController.ModelTrained)

}
