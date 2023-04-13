package controller

import (
	"errors"
	"log"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"school-notification-backend/util"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileController interface {
	CreateNewProfile(c *fiber.Ctx) error
	GetProfileById(c *fiber.Ctx) error
	GetProfileAll(c *fiber.Ctx) error
	GetProfileTeacherByCategory(c *fiber.Ctx) error
}

type profileController struct {
	profileRepo          repository.ProfileRepository
	classRepo            repository.ClassRepository
	schoolDataRepository repository.SchoolDataRepository
}

func NewProfileController(profileRepo repository.ProfileRepository, classRepo repository.ClassRepository, schoolDataRepository repository.SchoolDataRepository) ProfileController {
	return &profileController{profileRepo: profileRepo, classRepo: classRepo, schoolDataRepository: schoolDataRepository}
}

func (p *profileController) GetProfileAll(c *fiber.Ctx) error {

	role, err := util.CheckStringData(c.Query("role"), "role")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find profile role is:", role)

	profiles, err := p.profileRepo.GetAll(role)
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"role":         role,
		"profile_list": profiles,
	})
}

func (p *profileController) GetProfileById(c *fiber.Ctx) error {

	id, err := util.CheckStringData(c.Query("id"), "id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find profile id is:", id)

	role, err := util.CheckStringData(c.Query("role"), "role")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find profile role is:", role)

	filter := bson.M{
		"profile_id": id,
		"role":       role,
	}

	profile, err := p.profileRepo.GetProfileById(filter, role)
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

	if profile == nil {
		log.Println("Profile not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"profile": profile,
	})
}

func (p *profileController) GetProfileTeacherByCategory(c *fiber.Ctx) error {

	category, err := util.CheckStringData(c.Query("category"), "category")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find profile category:", category)

	filter := bson.M{
		"category": category,
		"role":     "teacher",
	}

	profile, err := p.profileRepo.GetProfileById(filter, "teacher")
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

	if profile == nil {
		log.Println("Profile not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"profile": profile,
	})
}

func (p *profileController) CreateNewProfile(c *fiber.Ctx) error {

	req := models.ProfileRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	req.Role, err = util.CheckStringData(req.Role, "role")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("create profile role:", req.Role)

	req.ProfileId, err = util.CheckStringData(req.ProfileId, "profile_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("profile id:", req.ProfileId)

	// validate and create profile
	var profile interface{}
	if req.Role == "teacher" {
		profile, err = newTeacherProfile(req, p.profileRepo, p.schoolDataRepository)
	} else if req.Role == "student" {
		profile, err = newStudentProfile(req, p.profileRepo, p.classRepo)
	} else if req.Role == "parent" {
		// profile, err = newParentProfile(req, p.profileRepo)
	} else {
		log.Println("role is invalid")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "role"+util.ErrValueInvalid.Error())
	}
	if err != nil {
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}

	profileInsert, err := p.profileRepo.Insert(profile)
	// insert
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	log.Println("create profile success")
	log.Println("profile id:", profileInsert.InsertedID)

	return util.ResponseSuccess(c, fiber.StatusCreated, "create profile success", map[string]interface{}{
		"profile_id": req.ProfileId,
	})
}

func newTeacherProfile(req models.ProfileRequest, profileRepo repository.ProfileRepository, schoolDataRepository repository.SchoolDataRepository) (*models.ProfileTeacher, error) {

	err := profileRepo.GetProfileByFilterForCheckExists(bson.M{
		"profile_id": req.ProfileId,
		"role":       req.Role})
	if err == nil {
		log.Println(util.ErrProfileIdAlreadyExists)
		return nil, util.ErrProfileIdAlreadyExists
	}
	if err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return nil, util.ErrInternalServerError
	}

	name, err := util.CheckStringData(req.Name, "name")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("create profile name:", name)

	category, err := util.CheckStringData(req.Category, "category")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("create profile category:", category)

	dataListSub, err := schoolDataRepository.GetByFilterAll(bson.M{"type": "SubjectCategory"})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(dataListSub) == 0 {
		log.Println("data SubjectCategory not found")
		return nil, util.ErrNotFound
	}

	check := true
	for _, v := range dataListSub {
		if *v.SubjectCategory == category {
			check = false
			break
		}
	}

	if check {
		log.Println("category", util.ErrValueNotAlreadyExists)
		return nil, util.ReturnError("category" + util.ErrValueNotAlreadyExists.Error())
	}

	dataList, err := schoolDataRepository.GetByFilterAll(bson.M{"type": "YearAndTerm"})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(dataList) == 0 {
		log.Println("data YearAndTerm not found")
		return nil, util.ErrNotFound
	}

	sort.Slice(dataList, func(i, j int) bool {
		return dataList[i].CreatedAt > dataList[j].CreatedAt
	})

	if *dataList[len(dataList)-1].Status == true {
		log.Println("school data invalid")
		return nil, errors.New("school data invalid")
	}

	year := *dataList[len(dataList)-1].Year
	term := *dataList[len(dataList)-1].Term

	p := models.ProfileTeacher{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		ProfileId: req.ProfileId,
		Name:      name,
		Role:      req.Role,
		Category:  category,
		Slot:      createTimeSlot(),
		CourseTeachesList: []models.CourseTeachesList{
			{
				Year: year,
				Term: term,
			},
		},
	}

	return &p, nil
}

func newStudentProfile(req models.ProfileRequest, profileRepo repository.ProfileRepository, classRepo repository.ClassRepository) (*models.ProfileStudent, error) {

	err := profileRepo.GetProfileByFilterForCheckExists(bson.M{
		"profile_id": req.ProfileId,
		"role":       req.Role})
	if err == nil {
		log.Println(util.ErrProfileIdAlreadyExists)
		return nil, util.ErrProfileIdAlreadyExists
	}
	if err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return nil, util.ErrInternalServerError
	}

	name, err := util.CheckStringData(req.Name, "name")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("create profile name:", name)

	classId, err := util.CheckStringData(req.ClassId, "class_id")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("create profile in class id:", classId)

	class, err := classRepo.GetClassById(req.ClassId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	check := true
	for _, s := range class.StudentIdList {
		if s == req.ProfileId {
			check = false
			break
		}
	}
	if check {
		class.StudentIdList = append(class.StudentIdList, req.ProfileId)
		class.NumberOfStudent = len(class.StudentIdList)
		_, err = classRepo.Update(class)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	p := models.ProfileStudent{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Name:      req.Name,
		Role:      req.Role,
		ProfileId: req.ProfileId,
		ClassId:   req.ClassId,
		TermScore: []models.TermScore{
			{
				Year: class.Year,
				Term: class.Term,
			},
		},
	}

	return &p, nil
}
