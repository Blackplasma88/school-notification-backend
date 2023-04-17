package controller

import (
	"fmt"
	"log"
	"os"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"school-notification-backend/security"
	"school-notification-backend/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InformationController interface {
	CreateInformation(c *fiber.Ctx) error
	UpdateInformation(c *fiber.Ctx) error
	GetInformationAll(c *fiber.Ctx) error
	GetInformationById(c *fiber.Ctx) error
}

type informationController struct {
	infoRepo repository.InformationRepository
	userRepo repository.UsersRepository
}

func NewInformationController(infoRepo repository.InformationRepository, userRepo repository.UsersRepository) InformationController {
	return &informationController{infoRepo: infoRepo, userRepo: userRepo}
}

func (i *informationController) CreateInformation(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], i.userRepo, []string{"admin"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	name, err := util.CheckStringData(c.FormValue("name"), "name")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("name:", name)

	description, err := util.CheckStringData(c.FormValue("description"), "description")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("description:", description)

	content, err := util.CheckStringData(c.FormValue("content"), "content")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("content:", content)

	category, err := util.CheckStringData(c.FormValue("category"), "category")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("category:", category)

	imageUrl := ""
	file, err := c.FormFile("file")
	if err == nil {
		log.Println("file type:", file.Header.Get("Content-Type"))

		if _, err := os.Stat("./storage"); os.IsNotExist(err) {
			err = os.Mkdir("./storage", 0777)
		}

		if _, err := os.Stat("./storage/information"); os.IsNotExist(err) {
			err = os.Mkdir("./storage/information", 0777)
		}

		filename := name + "-" + file.Filename

		err = c.SaveFile(file, fmt.Sprintf("./storage/information/%s", filename))
		if err != nil {
			log.Println("file save error --> ", err)
			value, ok := err.(*fiber.Error)
			if ok {
				return util.ResponseNotSuccess(c, value.Code, value.Message)
			}

			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, err.Error())
		}

		imageUrl = fmt.Sprintf("/files/information/%s", filename)
	}

	informationNew := &models.Information{
		Id:          primitive.NewObjectID(),
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
		Name:        name,
		Description: description,
		Content:     content,
		Category:    category,
		FilePath:    imageUrl,
	}

	information, err := i.infoRepo.Insert(informationNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create information success", map[string]interface{}{
		"information_id": information.InsertedID,
	})
}

func (i *informationController) UpdateInformation(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], i.userRepo, []string{"admin"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	id, err := util.CheckStringData(c.FormValue("id"), "id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("id:", id)

	information, err := i.infoRepo.GetInformationById(id)
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

	name, err := util.CheckStringData(c.FormValue("name"), "name")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("name:", name)

	description, err := util.CheckStringData(c.FormValue("description"), "description")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("description:", description)

	content, err := util.CheckStringData(c.FormValue("content"), "content")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("content:", content)

	category, err := util.CheckStringData(c.FormValue("category"), "category")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("category:", category)

	imageUrl := ""
	file, err := c.FormFile("file")
	if err == nil {
		log.Println("file type:", file.Header.Get("Content-Type"))

		if _, err := os.Stat("./storage"); os.IsNotExist(err) {
			err = os.Mkdir("./storage", 0777)
		}

		if _, err := os.Stat("./storage/information"); os.IsNotExist(err) {
			err = os.Mkdir("./storage/information", 0777)
		}

		filename := name + "-" + file.Filename

		err = c.SaveFile(file, fmt.Sprintf("./storage/information/%s", filename))
		if err != nil {
			log.Println("file save error --> ", err)
			value, ok := err.(*fiber.Error)
			if ok {
				return util.ResponseNotSuccess(c, value.Code, value.Message)
			}

			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, err.Error())
		}

		imageUrl = fmt.Sprintf("/files/information/%s", filename)
	}

	information.UpdatedAt = time.Now().Format(time.RFC3339)
	information.Name = name
	information.Description = description
	information.Content = content
	information.Category = category
	if imageUrl != "" {
		information.FilePath = imageUrl
	}

	informationUpdate, err := i.infoRepo.Update(information)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "update information success", map[string]interface{}{
		"information_id": information.Id,
		"update_count":   informationUpdate.ModifiedCount,
	})
}

func (i *informationController) GetInformationAll(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], i.userRepo, []string{"all"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	infos, err := i.infoRepo.GetAll()
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(infos) == 0 {
		log.Println("Information not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"information_list": infos,
	})
}

func (i *informationController) GetInformationById(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], i.userRepo, []string{"all"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	id, err := util.CheckStringData(c.Query("id"), "id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find information id:", id)

	information, err := i.infoRepo.GetInformationById(id)
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

	if information == nil {
		log.Println("information not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"information": information,
	})
}
