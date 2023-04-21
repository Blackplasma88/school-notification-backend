package controller

import (
	"fmt"
	"log"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"school-notification-backend/security"
	"school-notification-backend/util"
	"sort"

	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassController interface {
	CreateClass(c *fiber.Ctx) error
	// GetClassAll(c *fiber.Ctx) error
	GetClassAllByClassYear(c *fiber.Ctx) error
	GetClassById(c *fiber.Ctx) error
	GetClassByClassYearAndRoom(c *fiber.Ctx) error
	SetAdvisor(c *fiber.Ctx) error
}

type classController struct {
	classRepo            repository.ClassRepository
	schoolDataRepository repository.SchoolDataRepository
	profileRepo          repository.ProfileRepository
	userRepo             repository.UsersRepository
	faceDetectionRepo    repository.FaceDetectionRepository
}

func NewClassController(classRepo repository.ClassRepository, schoolDataRepository repository.SchoolDataRepository, profileRepo repository.ProfileRepository, userRepo repository.UsersRepository, faceDetectionRepo repository.FaceDetectionRepository) ClassController {
	return &classController{classRepo: classRepo, schoolDataRepository: schoolDataRepository, profileRepo: profileRepo, userRepo: userRepo, faceDetectionRepo: faceDetectionRepo}
}

func (cl *classController) CreateClass(c *fiber.Ctx) error {

	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cl.userRepo, []string{"admin"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	num, err := cl.classRepo.GetCountOfClassYear("1")
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	dataList, err := cl.schoolDataRepository.GetByFilterAll(bson.M{"type": "YearAndTerm"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	sort.Slice(dataList, func(i, j int) bool {
		return dataList[i].CreatedAt > dataList[j].CreatedAt
	})

	if *dataList[0].Status == true {
		log.Println("school data invalid")
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, "school data invalid")
	}

	year := *dataList[0].Year
	term := *dataList[0].Term

	classNew := &models.ClassData{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		ClassYear: "1",
		ClassRoom: fmt.Sprint(num + 1),
		Status:    false,
		Year:      year,
		Term:      term,
		Slot:      createTimeSlot(),
	}

	class, err := cl.classRepo.Insert(classNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	dataNew := &models.FaceDetectData{
		Id:                   primitive.NewObjectID(),
		CreatedAt:            time.Now().Format(time.RFC3339),
		UpdatedAt:            time.Now().Format(time.RFC3339),
		Status:               "not",
		Name:                 classNew.ClassYear + "/" + classNew.ClassRoom,
		ClassId:              classNew.Id.Hex(),
		NumberOfStudent:      classNew.NumberOfStudent,
		StudentIdList:        classNew.StudentIdList,
		NumberOfImage:        0,
		ImageStudentPathList: createEmptyImageList(classNew.NumberOfStudent),
	}

	_, err = cl.faceDetectionRepo.Insert(dataNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create class success", map[string]interface{}{
		"class_id": class.InsertedID,
	})
}

func (cl *classController) GetClassAllByClassYear(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cl.userRepo, []string{"all"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	classYear, err := util.CheckStringData(c.Query("class_year"), "class_year")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("class year:", classYear)

	classes, err := cl.classRepo.GetClassByFilterAll(bson.M{"class_year": classYear})
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

	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cl.userRepo, []string{"all"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	id, err := util.CheckStringData(c.Query("class_id"), "class_id")
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

func (cl *classController) GetClassByClassYearAndRoom(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cl.userRepo, []string{"all"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	classYear, err := util.CheckStringData(c.Query("class_year"), "class_year")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find class year:", classYear)

	classRoom, err := util.CheckStringData(c.Query("class_room"), "class_room")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find class room:", classRoom)

	class, err := cl.classRepo.GetClassByFilter(bson.M{"class_year": classYear, "class_room": classRoom, "status": false})
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

func (cl *classController) SetAdvisor(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cl.userRepo, []string{"admin"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	req := models.ClassRequest{}
	err = c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	id, err := util.CheckStringData(req.ClassId, "class_id")
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

	if class.AdvisorId != "" {
		log.Println("class has advisor")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "class has advisor")
	}

	advisorId, err := util.CheckStringData(req.AdvisorId, "advisor_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find advisord id:", advisorId)

	filter := bson.M{
		"profile_id": advisorId,
		"role":       "teacher",
	}

	p, err := cl.profileRepo.GetProfileById(filter, "teacher")
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

	profile, _ := p.(models.ProfileTeacher)

	if profile.ClassInCounseling != "" {
		log.Println("teacher has class in counseling")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "teacher has class in counseling")
	}

	class.AdvisorId = advisorId
	profile.ClassInCounseling = class.Id.Hex()

	_, err = cl.profileRepo.Update(profile.Id, profile)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	_, err = cl.classRepo.Update(class)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "add success", map[string]interface{}{
		"class_id":   class.Id,
		"advisor_id": profile.Id,
	})
}
