package controller

import (
	"log"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"school-notification-backend/security"
	"school-notification-backend/util"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CheckNameController interface {
	AddDateForCheck(c *fiber.Ctx) error
	GetDateByCourseId(c *fiber.Ctx) error
	CheckNameStudent(c *fiber.Ctx) error
	GetCheckNameDataByCourseIdAndDate(c *fiber.Ctx) error
	EndDateCheckName(c *fiber.Ctx) error
}

type checkNameController struct {
	checkNameRepository repository.CheckNameRepository
	courseRepo          repository.CourseRepository
	userRepo            repository.UsersRepository
}

func NewCheckNameController(checkNameRepository repository.CheckNameRepository, courseRepo repository.CourseRepository, userRepo repository.UsersRepository) CheckNameController {
	return &checkNameController{checkNameRepository: checkNameRepository, courseRepo: courseRepo, userRepo: userRepo}
}

func (cn *checkNameController) AddDateForCheck(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cn.userRepo, []string{"admin", "teacher"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	req := models.CheckNameRequest{}
	err = c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	courseId, err := util.CheckStringData(req.CourseId, "course_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("course id:", courseId)

	course, err := cn.courseRepo.GetCourseById(courseId)
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

	if course.Status != "progress" {
		log.Println("course status does not progress")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "course status does not progress")
	}

	date, err := util.CheckStringData(req.Date, "date")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("check name date:", date)
	timeLate, err := util.CheckIntegerData(req.TimeLate, "time_late")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("check name time late:", timeLate)

	tDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	weekDay := strings.ToLower(tDate.Weekday().String())
	log.Println(weekDay)

	checkDay := true
	for _, dt := range course.DateTime {
		if dt.Day == weekDay {
			checkDay = false
			break
		}
	}
	if checkDay {
		log.Println("day not found in course")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "day not found in course")
	}

	_, err = cn.checkNameRepository.GetByFilter(bson.M{"course_id": courseId, "date": date})
	if err == nil {
		log.Println("check name date already exists")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "check name date"+util.ErrValueAlreadyExists.Error())
	}
	if err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	t := time.Now()
	checkNameNew := &models.CheckName{
		Id:            primitive.NewObjectID(),
		CreatedAt:     t.Format(time.RFC3339),
		UpdatedAt:     t.Format(time.RFC3339),
		CourseId:      course.Id.Hex(),
		Date:          date,
		TimeLate:      t.Add(time.Minute * time.Duration(timeLate)).Format(time.RFC3339),
		Status:        "progress",
		CheckNameData: createCheckNameData(course.StudentIdList, t.Format(time.RFC3339)),
	}

	result, err := cn.checkNameRepository.Insert(checkNameNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create date check name success", map[string]interface{}{
		"id": result.InsertedID,
	})
}

func (cn *checkNameController) GetDateByCourseId(c *fiber.Ctx) error {
	user, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cn.userRepo, []string{"all"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}
	role := user.Role
	profileId := user.ProfileId

	courseId, err := util.CheckStringData(c.Query("course_id"), "course_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find score of course id:", courseId)
	course, err := cn.courseRepo.GetCourseById(courseId)
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

	checkNameList, err := cn.checkNameRepository.GetByFilterAll(bson.M{"course_id": courseId})
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(checkNameList) == 0 {
		log.Println("check name not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	if role == "teacher" {
		if profileId != course.InstructorId {
			log.Println("not permiistion")
			return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, "not permission")
		}

	} else if role == "student" {
		check := true
		for _, v := range course.StudentIdList {
			if v == profileId {
				check = false
				break
			}
		}
		if check {
			log.Println("not permiistion")
			return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, "not permission")
		}
	} else {
		if role != "server" {
			log.Println("role", util.ErrValueInvalid)
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "role"+util.ErrValueInvalid.Error())
		}
	}

	checkNamel := []string{}
	for _, v := range checkNameList {
		checkNamel = append(checkNamel, v.Date)
	}

	if len(checkNamel) == 0 {
		log.Println("date ", util.ErrNotFound)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "date "+util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"date_list": checkNamel,
	})
}

func (cn *checkNameController) CheckNameStudent(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cn.userRepo, []string{"admin", "teacher"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	req := models.CheckNameRequest{}
	err = c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	courseId, err := util.CheckStringData(req.CourseId, "course_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("course id:", courseId)

	course, err := cn.courseRepo.GetCourseById(courseId)
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

	if course.Status != "progress" {
		log.Println("course status does not progress")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "course status does not progress")
	}

	date, err := util.CheckStringData(req.Date, "date")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("check name date:", date)

	chcekName, err := cn.checkNameRepository.GetByFilter(bson.M{"course_id": courseId, "date": date})
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, "date "+util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	studentId, err := util.CheckStringData(req.StudentId, "student_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("student id:", studentId)

	checkBy, err := util.CheckStringData(req.CheckBy, "check_by")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("check by:", checkBy)

	if checkBy != "teacher" && checkBy != "server" {
		log.Println("check by", util.ErrValueInvalid)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "check by"+util.ErrValueInvalid.Error())
	}

	check := true
	for _, v := range course.StudentIdList {
		if v == studentId {
			check = false
			break
		}
	}

	if check {
		log.Println("student id not found in course")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "student id not found in course")
	}

	t := time.Now()
	for i, v := range chcekName.CheckNameData {
		if v.StudentId == studentId {
			chcekName.CheckNameData[i].UpdatedAt = t.Format(time.RFC3339)
			chcekName.CheckNameData[i].Time = strings.Split(strings.Split(t.Format(time.RFC3339), "T")[1], "+")[0]
			chcekName.CheckNameData[i].CheckBy = checkBy

			tl, err := time.Parse(time.RFC3339, chcekName.TimeLate)
			if err != nil {
				log.Println(err)
				return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, err.Error())
			}

			if !(t.After(tl)) {
				chcekName.CheckNameData[i].Status = "attend"
			} else {
				chcekName.CheckNameData[i].Status = "late"
			}
		}
	}

	result, err := cn.checkNameRepository.Update(chcekName)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "check name success", map[string]interface{}{
		"check_name_id": chcekName.Id,
		"update_count":  result.ModifiedCount,
	})
}

func (cn *checkNameController) GetCheckNameDataByCourseIdAndDate(c *fiber.Ctx) error {
	user, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cn.userRepo, []string{"all"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}
	// role := c.Query("role")
	role := user.Role
	// profileId := c.Query("id")
	profileId := user.ProfileId
	courseId, err := util.CheckStringData(c.Query("course_id"), "course_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find check name of course id:", courseId)
	course, err := cn.courseRepo.GetCourseById(courseId)
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

	var dataRes interface{}
	if role == "teacher" {
		date, err := util.CheckStringData(c.Query("date"), "date")
		if err != nil {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		log.Println("check name date:", date)

		if profileId != course.InstructorId {
			log.Println("not permiistion")
			return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, "not permission")
		}

		data, err := cn.checkNameRepository.GetByFilter(bson.M{"course_id": courseId, "date": date})
		if err != nil {
			log.Println(err)
			if err == mongo.ErrNoDocuments {
				return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
			}
			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
		}
		dataRes = data
	} else if role == "student" {

		check := true
		for _, v := range course.StudentIdList {
			if v == profileId {
				check = false
				break
			}
		}
		if check {
			log.Println("not permiistion")
			return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, "not permission")
		}

		checkNameList, err := cn.checkNameRepository.GetByFilterAll(bson.M{"course_id": courseId})
		if err != nil {
			log.Println(err)
			if err == mongo.ErrNoDocuments {
				return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
			}
			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
		}

		checkNameListRes := []models.CheckNameStudentRes{}
		for _, v := range checkNameList {
			for _, check := range v.CheckNameData {
				if check.StudentId == profileId {
					checkNameListRes = append(checkNameListRes, models.CheckNameStudentRes{
						Date:      v.Date,
						UpdatedAt: check.UpdatedAt,
						Time:      check.Time,
						Status:    check.Status,
						CheckBy:   check.CheckBy,
					})
					break
				}
			}
		}
		if len(checkNameListRes) == 0 {
			log.Println("student id not have checked")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "student id not have checked")
		}
		dataRes = checkNameListRes
	} else {
		if role != "server" {
			log.Println("role", util.ErrValueInvalid)
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "role"+util.ErrValueInvalid.Error())
		}
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"date_data": dataRes,
	})
}

func (cn *checkNameController) EndDateCheckName(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cn.userRepo, []string{"admin", "teacher"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	req := models.CheckNameRequest{}
	err = c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	courseId, err := util.CheckStringData(req.CourseId, "course_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("course id:", courseId)

	course, err := cn.courseRepo.GetCourseById(courseId)
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

	if course.Status != "progress" {
		log.Println("course status does not progress")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "course status does not progress")
	}

	date, err := util.CheckStringData(req.Date, "date")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("check name date:", date)

	chcekName, err := cn.checkNameRepository.GetByFilter(bson.M{"course_id": courseId, "date": date})
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if chcekName.Status != "progress" {
		log.Println("this date not in progress")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "this date not in progress")
	}

	t := time.Now()
	// for _, v := range course.StudentIdList {
	// 	check := true
	for i, d := range chcekName.CheckNameData {
		if d.Status == "" {
			chcekName.CheckNameData[i].UpdatedAt = t.Format(time.RFC3339)
			chcekName.CheckNameData[i].Time = strings.Split(strings.Split(t.Format(time.RFC3339), "T")[1], "+")[0]
			chcekName.CheckNameData[i].Status = "absent"
			chcekName.CheckNameData[i].CheckBy = "server"
		}
		// if v == d.StudentId {
		// 	check = false
		// 	break
		// }
	}
	// if check {
	// 	chcekName.CheckNameData = append(chcekName.CheckNameData, models.CheckNameData{
	// 		StudentId: v,
	// 		UpdatedAt: t.Format(time.RFC3339),
	// 		Time:      strings.Split(strings.Split(t.Format(time.RFC3339), "T")[1], "+")[0],
	// 		Status:    "absent",
	// 		CheckBy:   "server",
	// 	})
	// }
	// }

	chcekName.Status = "end"

	result, err := cn.checkNameRepository.Update(chcekName)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "check name success", map[string]interface{}{
		"check_name_id": chcekName.Id,
		"update_count":  result.ModifiedCount,
	})
}

func createCheckNameData(studentIdList []string, t string) []models.CheckNameData {

	var res []models.CheckNameData
	for _, s := range studentIdList {
		res = append(res, models.CheckNameData{
			StudentId: s,
			UpdatedAt: t,
		})
	}

	return res
}
