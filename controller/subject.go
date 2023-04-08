package controller

import (
	"log"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"school-notification-backend/util"

	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubjectController interface {
	CreateSubject(c *fiber.Ctx) error
	UpdateSubject(c *fiber.Ctx) error
	GetSubjectAll(c *fiber.Ctx) error
	GetSubjectById(c *fiber.Ctx) error
}

type subjectController struct {
	subjectRepository repository.SubjectRepository
}

func NewSubjectController(subjectRepository repository.SubjectRepository) SubjectController {
	return &subjectController{subjectRepository: subjectRepository}
}

func (s *subjectController) CreateSubject(c *fiber.Ctx) error {

	req := models.SubjectRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	subjectId, err := util.CheckStringData(req.SubjectId, "subject_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject id:", subjectId)

	_, err = s.subjectRepository.GetSubjectByFilter(bson.M{"subject_id": subjectId})
	if err == nil {
		log.Println("subject id already exists")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "subject id"+util.ErrValueAlreadyExists.Error())
	}
	if err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	name, err := util.CheckStringData(req.Name, "name")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject name:", name)

	category, err := util.CheckStringData(req.Category, "category")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject category:", category)

	classYear, err := util.CheckStringData(req.ClassYear, "class_year")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject class year:", classYear)

	credit, err := util.CheckIntegerData(req.Credit, "credit")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject credit:", credit)

	instructorId := []string{}
	for _, v := range req.InstructorId {
		inId, err := util.CheckStringData(v, "instructor_id")
		if err != nil {
			continue
		}
		log.Println("instructor id:", inId)
		instructorId = append(instructorId, inId)
	}

	subjectNew := &models.Subject{
		Id:           primitive.NewObjectID(),
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
		SubjectId:    subjectId,
		Name:         name,
		Credit:       credit,
		Category:     category,
		ClassYear:    classYear,
		InstructorId: instructorId,
	}

	_, err = s.subjectRepository.Insert(subjectNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create subject success", map[string]interface{}{
		"subject_id": req.SubjectId,
	})
}

func (s *subjectController) UpdateSubject(c *fiber.Ctx) error {

	req := models.SubjectRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	subjectId, err := util.CheckStringData(req.SubjectId, "subject_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject id:", subjectId)

	subject, err := s.subjectRepository.GetSubjectByFilter(bson.M{"subject_id": subjectId})
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

	name, err := util.CheckStringData(req.Name, "name")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject name:", name)

	category, err := util.CheckStringData(req.Category, "category")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject category:", category)

	classYear, err := util.CheckStringData(req.ClassYear, "class_year")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject class year:", classYear)

	credit, err := util.CheckIntegerData(req.Credit, "credit")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject credit:", credit)

	instructorId := []string{}
	for _, v := range req.InstructorId {
		inId, err := util.CheckStringData(v, "instructor_id")
		if err != nil {
			continue
		}
		log.Println("instructor id:", inId)
		instructorId = append(instructorId, inId)
	}

	subject.UpdatedAt = time.Now().Format(time.RFC3339)
	subject.Name = name
	subject.Category = category
	subject.Credit = credit
	subject.ClassYear = classYear
	subject.InstructorId = instructorId

	result, err := s.subjectRepository.Update(subject)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create subject success", map[string]interface{}{
		"subject_id":   req.SubjectId,
		"update_count": result.ModifiedCount,
	})
}

func (s *subjectController) GetSubjectAll(c *fiber.Ctx) error {

	subjects, err := s.subjectRepository.GetAll()
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(subjects) == 0 {
		log.Println("Subject not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"subject_list": subjects,
	})
}

func (s *subjectController) GetSubjectById(c *fiber.Ctx) error {
	id, err := util.CheckStringData(c.Query("subject_id"), "subject_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find subject id:", id)

	subject, err := s.subjectRepository.GetSubjectByFilter(bson.M{"subject_id": id})
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

	if subject == nil {
		log.Println("Subject not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"subject": subject,
	})
}
