package controller

import (
	"log"
	"net/http"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"school-notification-backend/security"
	"school-notification-backend/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FaceDetectionController interface {
	OpenCamera(c *fiber.Ctx) error
	CreatFaceDetectionData(c *fiber.Ctx) error
	UploadImageData(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	ModelTrained(c *fiber.Ctx) error
	GetByClassId(c *fiber.Ctx) error
}

type faceDetectionController struct {
	faceDetectionRepo repository.FaceDetectionRepository
	classRepo         repository.ClassRepository
	userRepo          repository.UsersRepository
}

func NewFaceDetectionController(faceDetectionRepo repository.FaceDetectionRepository, classRepo repository.ClassRepository, userRepo repository.UsersRepository) FaceDetectionController {
	return &faceDetectionController{faceDetectionRepo: faceDetectionRepo, classRepo: classRepo, userRepo: userRepo}
}

func (f *faceDetectionController) OpenCamera(c *fiber.Ctx) error {
	classId, err := util.CheckStringData(c.Query("class_id"), "class_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("class id:", classId)

	courseId, err := util.CheckStringData(c.Query("course_id"), "course_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("course id:", courseId)

	log.Println("run camera")
	go func() {
		log.Println("call")
		_, err := http.Get("http://localhost:8000/cv?class_id=" + classId + "&course_id=" + courseId)
		if err != nil {
			return
		}
		log.Println("close camera")
	}()

	log.Println("call success")
	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"success": true,
	})
}

func (f *faceDetectionController) CreatFaceDetectionData(c *fiber.Ctx) error {
	err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], f.userRepo, []string{"admin"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	req := models.FaceDetectDataRequest{}
	err = c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	classId, err := util.CheckStringData(req.ClassId, "class_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("class id:", classId)

	_, err = f.faceDetectionRepo.GetByFilter(bson.M{"class_id": classId})
	if err == nil {
		log.Println("data for class id already exists")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "data for class"+util.ErrValueAlreadyExists.Error())
	}
	if err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	class, err := f.classRepo.GetClassById(classId)
	if err != nil {
		log.Println(err)
		if err.Error() == "mongo: no documents in result" {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		if err.Error() == "Id is not primitive objectID" {
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	dataNew := &models.FaceDetectData{
		Id:                   primitive.NewObjectID(),
		CreatedAt:            time.Now().Format(time.RFC3339),
		UpdatedAt:            time.Now().Format(time.RFC3339),
		Status:               "not",
		Name:                 class.ClassYear + "/" + class.ClassRoom,
		ClassId:              classId,
		NumberOfStudent:      class.NumberOfStudent,
		StudentIdList:        class.StudentIdList,
		NumberOfImage:        0,
		ImageStudentPathList: createEmptyImageList(class.NumberOfStudent),
	}

	data, err := f.faceDetectionRepo.Insert(dataNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create face detection data success", map[string]interface{}{
		"data_id": data.InsertedID,
	})
}

func (f *faceDetectionController) UploadImageData(c *fiber.Ctx) error {
	err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], f.userRepo, []string{"admin"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	req := models.FaceDetectDataRequest{}
	err = c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	id, err := util.CheckStringData(req.Id, "id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("id:", id)

	data, err := f.faceDetectionRepo.GetById(id)
	if err != nil {
		log.Println(err)
		if err.Error() == "mongo: no documents in result" {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		if err.Error() == "Id is not primitive objectID" {
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	studentId, err := util.CheckStringData(req.StudentId, "student_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("student id:", studentId)

	index := -1
	for i, v := range data.StudentIdList {
		if v == studentId {
			index = i
			break
		}
	}

	if index == -1 {
		log.Println("not found student")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error()+" student")
	}

	if len(req.ImagePathList) == 0 {
		log.Println(util.ReturnError(util.ErrRequireParameter.Error() + "image_path_list"))
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"image_path_list")
	}

	for _, v := range req.ImagePathList {
		data.ImageStudentPathList[index] = append(data.ImageStudentPathList[index], v)
	}

	data.Status = "not"

	_, err = f.faceDetectionRepo.Update(data)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "update success", map[string]interface{}{
		"data_id": data.Id,
	})
}

func (f *faceDetectionController) ModelTrained(c *fiber.Ctx) error {
	err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], f.userRepo, []string{"admin"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	req := models.FaceDetectDataRequest{}
	err = c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	id, err := util.CheckStringData(req.Id, "id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("id:", id)

	data, err := f.faceDetectionRepo.GetById(id)
	if err != nil {
		log.Println(err)
		if err.Error() == "mongo: no documents in result" {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		if err.Error() == "Id is not primitive objectID" {
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	classId, err := util.CheckStringData(req.ClassId, "class_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("class id:", classId)

	if classId != data.ClassId {
		log.Println("class id not match")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "class id not match")
	}

	data.Status = "progress"

	go func() {
		log.Println("call")
		_, err := http.Post("http://localhost:8000/train-model?class_id="+classId, "application/json", nil)
		if err != nil {
			return
		}
		data.Status = "yes"
		data.UpdatedAt = time.Now().Format(time.RFC3339)
		log.Println("finish")
		_, err = f.faceDetectionRepo.Update(data)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	data.UpdatedAt = time.Now().Format(time.RFC3339)
	_, err = f.faceDetectionRepo.Update(data)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "update success", map[string]interface{}{
		"data_id": data.Id,
	})
}

func (f *faceDetectionController) GetAll(c *fiber.Ctx) error {

	datas, err := f.faceDetectionRepo.GetAll()
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(datas) == 0 {
		log.Println("data not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"data_list": datas,
	})
}

func (f *faceDetectionController) GetById(c *fiber.Ctx) error {
	id, err := util.CheckStringData(c.Query("id"), "id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find id:", id)

	data, err := f.faceDetectionRepo.GetById(id)
	if err != nil {
		log.Println(err)
		if err.Error() == "mongo: no documents in result" {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		if err.Error() == "Id is not primitive objectID" {
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if data == nil {
		log.Println("data not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"data": data,
	})
}

func (f *faceDetectionController) GetByClassId(c *fiber.Ctx) error {
	classId, err := util.CheckStringData(c.Query("class_id"), "class_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("class id:", classId)

	data, err := f.faceDetectionRepo.GetByFilter(bson.M{"class_id": classId})
	if err != nil {
		log.Println(err)
		if err.Error() == "mongo: no documents in result" {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		if err.Error() == "Id is not primitive objectID" {
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if data == nil {
		log.Println("data not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"data": data,
	})
}

func createEmptyImageList(num int) [][]string {
	var res [][]string
	for i := 0; i < num; i++ {
		res = append(res, []string{})
	}

	return res
}
