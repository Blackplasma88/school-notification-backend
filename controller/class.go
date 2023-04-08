package controller

import (
	"fmt"
	"log"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"school-notification-backend/util"

	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassController interface {
	CreateClass(c *fiber.Ctx) error
	GetClassAll(c *fiber.Ctx) error
	GetClassById(c *fiber.Ctx) error
}

type classController struct {
	classRepo repository.ClassRepository
}

func NewClassController(classRepo repository.ClassRepository) ClassController {
	return &classController{classRepo: classRepo}
}

func (cl *classController) CreateClass(c *fiber.Ctx) error {

	num, err := cl.classRepo.GetCountOfClassYear("1")
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	classNew := &models.ClassData{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		ClassYear: "1",
		ClassRoom: fmt.Sprint(num + 1),
		Year:      time.Now().Format("%Y"),
		Term:      "1",
	}

	class, err := cl.classRepo.Insert(classNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create class success", map[string]interface{}{
		"class_id": class.InsertedID,
	})
}

func (cl *classController) GetClassAll(c *fiber.Ctx) error {

	classes, err := cl.classRepo.GetAll()
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(classes) == 0 {
		log.Println("class not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"class_list": classes,
	})
}

func (cl *classController) GetClassById(c *fiber.Ctx) error {

	id, err := util.CheckStringData(c.Query("id"), "id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find class id:", id)

	class, err := cl.classRepo.GetClassById(id)
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

	if class == nil {
		log.Println("class not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"class": class,
	})
}
