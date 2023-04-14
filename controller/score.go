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

type ScoreController interface {
	CreateScore(c *fiber.Ctx) error
	// AddStudentScore(c *fiber.Ctx) error
	UpdateStudentScore(c *fiber.Ctx) error
	GetScoreByCourseId(c *fiber.Ctx) error
	GetScoreDataByCourseIdAndNameSore(c *fiber.Ctx) error
}

type scoreController struct {
	scoreRepository repository.ScoreRepository
	courseRepo      repository.CourseRepository
}

func NewScoreController(scoreRepository repository.ScoreRepository, courseRepo repository.CourseRepository) ScoreController {
	return &scoreController{scoreRepository: scoreRepository, courseRepo: courseRepo}
}

func (s *scoreController) CreateScore(c *fiber.Ctx) error {

	req := models.ScoreRequest{}
	err := c.BodyParser(&req)
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

	course, err := s.courseRepo.GetCourseById(courseId)
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

	name, err := util.CheckStringData(req.Name, "name")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("Score name:", name)

	scoreFull, err := util.CheckFloatData(req.ScoreFull, "score_full")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("score full:", scoreFull)

	_, err = s.scoreRepository.GetScoreByFilter(bson.M{"course_id": courseId, "name": name})
	if err == nil {
		log.Println("score name already exists")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "score name"+util.ErrValueAlreadyExists.Error())
	}
	if err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	typeScore, err := util.CheckStringData(req.Type, "type")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("Score type:", typeScore)
	if typeScore != "midterm" && typeScore != "final" && typeScore != "work" {
		log.Println(util.ErrTypeInvalid)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrTypeInvalid.Error())
	}
	if typeScore == "midterm" || typeScore == "final" {
		_, err = s.scoreRepository.GetScoreByFilter(bson.M{"course_id": courseId, "type": typeScore})
		if err == nil {
			log.Println("score type already exists")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "score type"+util.ErrValueAlreadyExists.Error())
		}
		if err.Error() != "mongo: no documents in result" {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
		}
	}

	t := time.Now().Format(time.RFC3339)
	scoreNew := &models.Score{
		Id:               primitive.NewObjectID(),
		CreatedAt:        t,
		UpdatedAt:        t,
		CourseId:         course.Id.Hex(),
		Name:             name,
		ScoreFull:        scoreFull,
		Type:             typeScore,
		ScoreInformation: createScoreinformation(course.StudentIdList, t),
	}

	result, err := s.scoreRepository.Insert(scoreNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create Score success", map[string]interface{}{
		"score_id": result.InsertedID,
	})
}

func (s *scoreController) UpdateStudentScore(c *fiber.Ctx) error {

	req := models.ScoreRequest{}
	err := c.BodyParser(&req)
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

	course, err := s.courseRepo.GetCourseById(courseId)
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

	name, err := util.CheckStringData(req.Name, "name")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("Score name:", name)

	scoreGet, err := util.CheckFloatData(req.ScoreGet, "score_get")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("score get:", scoreGet)

	score, err := s.scoreRepository.GetScoreByFilter(bson.M{"course_id": courseId, "name": name})
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	studentId, err := util.CheckStringData(req.StudentId, "student_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("student id:", studentId)

	status, err := util.CheckStringData(req.Status, "status")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("status:", status)

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

	// check status
	if status != "normal" && status != "late" {
		log.Println("status", util.ErrValueInvalid)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "status"+util.ErrValueInvalid.Error())
	}

	if scoreGet > score.ScoreFull {
		log.Println("score get", util.ErrValueInvalid)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "score get"+util.ErrValueInvalid.Error())
	}

	if len(score.ScoreInformation) != len(course.StudentIdList) {
		log.Println("score info", util.ErrValueInvalid)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "score info"+util.ErrValueInvalid.Error())
	}

	for i, v := range score.ScoreInformation {
		if v.StudentId == studentId {
			score.ScoreInformation[i].UpdatedAt = time.Now().Format(time.RFC3339)
			score.ScoreInformation[i].ScoreGet = &scoreGet
			score.ScoreInformation[i].Status = status
			break
		}
	}

	// score.ScoreInformation = append(score.ScoreInformation, models.ScoreInformation{
	// 	StudentId: studentId,
	// 	UpdatedAt: time.Now().Format(time.RFC3339),
	// 	ScoreGet:  scoreGet,
	// 	Status:    status,
	// })

	result, err := s.scoreRepository.Update(score)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "update subject success", map[string]interface{}{
		"score_id":     score.Id,
		"update_count": result.ModifiedCount,
	})
}

// func (s *scoreController) UpdateStudentScore(c *fiber.Ctx) error {

// 	req := models.ScoreRequest{}
// 	err := c.BodyParser(&req)
// 	if err != nil {
// 		log.Println(err)
// 		value, ok := err.(*fiber.Error)
// 		if ok {
// 			return util.ResponseNotSuccess(c, value.Code, value.Message)
// 		}

// 		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
// 	}

// 	courseId, err := util.CheckStringData(req.CourseId, "course_id")
// 	if err != nil {
// 		log.Println(err)
// 		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
// 	}
// 	log.Println("course id:", courseId)

// 	course, err := s.courseRepo.GetCourseById(courseId)
// 	if err != nil {
// 		log.Println(err)
// 		if err.Error() == "mongo: no documents in result" {
// 			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
// 		}
// 		if err.Error() == "Id is not primitive objectID" {
// 			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
// 		}
// 		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
// 	}

// 	if course.Status != "progress" {
// 		log.Println("course status does not progress")
// 		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "course status does not progress")
// 	}

// 	name, err := util.CheckStringData(req.Name, "name")
// 	if err != nil {
// 		log.Println(err)
// 		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
// 	}
// 	log.Println("Score name:", name)

// 	scoreGet, err := util.CheckFloatData(req.ScoreGet, "score_get")
// 	if err != nil {
// 		log.Println(err)
// 		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
// 	}
// 	log.Println("score get:", scoreGet)

// 	score, err := s.scoreRepository.GetScoreByFilter(bson.M{"course_id": courseId, "name": name})
// 	if err != nil {
// 		log.Println(err)
// 		if err == mongo.ErrNoDocuments {
// 			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
// 		}
// 		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
// 	}

// 	studentId, err := util.CheckStringData(req.StudentId, "student_id")
// 	if err != nil {
// 		log.Println(err)
// 		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
// 	}
// 	log.Println("student id:", studentId)

// 	status, err := util.CheckStringData(req.Status, "status")
// 	if err != nil {
// 		log.Println(err)
// 		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
// 	}
// 	log.Println("status:", status)

// 	check := true
// 	for _, v := range course.StudentIdList {
// 		if v == studentId {
// 			check = false
// 			break
// 		}
// 	}

// 	if check {
// 		log.Println("student id not found in course")
// 		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "student id not found in course")
// 	}

// 	index := -1
// 	for i, v := range score.ScoreInformation {
// 		if v.StudentId == studentId {
// 			index = i
// 			break
// 		}
// 	}

// 	if index == -1 {
// 		log.Println("student id", util.ErrValueNotAlreadyExists)
// 		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "student id"+util.ErrValueNotAlreadyExists.Error())
// 	}

// 	// check status
// 	if status != "normal" && status != "late" {
// 		log.Println("status", util.ErrValueInvalid)
// 		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "status"+util.ErrValueInvalid.Error())
// 	}

// 	score.ScoreInformation[index].ScoreGet = scoreGet
// 	score.ScoreInformation[index].Status = status
// 	score.ScoreInformation[index].UpdatedAt = time.Now().Format(time.RFC3339)

// 	result, err := s.scoreRepository.Update(score)
// 	if err != nil {
// 		log.Println(err)
// 		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
// 	}

// 	return util.ResponseSuccess(c, fiber.StatusCreated, "update subject success", map[string]interface{}{
// 		"score_id":     score.Id,
// 		"update_count": result.ModifiedCount,
// 	})
// }

func (s *scoreController) GetScoreByCourseId(c *fiber.Ctx) error {
	role := c.Query("role")
	profileId := c.Query("id")
	courseId, err := util.CheckStringData(c.Query("course_id"), "course_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find score of course id:", courseId)
	course, err := s.courseRepo.GetCourseById(courseId)
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

	scores, err := s.scoreRepository.GetByFilterAll(bson.M{"course_id": courseId})
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(scores) == 0 {
		log.Println("score not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	scoreRes := make(map[string]interface{})
	if role == "teacher" {
		if profileId != course.InstructorId {
			log.Println("not permiistion")
			return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, "not permission")
		}
		scoreNames := models.ScoreTeacherRes{}
		for _, v := range scores {
			scoreNames.Name = append(scoreNames.Name, v.Name)
		}
		scoreRes["score"] = scoreNames
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
		scoreList := []models.ScoreStudentRes{}
		for _, v := range scores {
			// data := models.ScoreStudentRes{
			// 	Name:      v.Name,
			// 	ScoreFull: v.ScoreFull,
			// }
			for _, info := range v.ScoreInformation {
				if info.StudentId == profileId {
					// data.UpdatedAt = info.UpdatedAt
					// data.ScoreGet = info.ScoreGet
					// data.Status = info.Status
					scoreList = append(scoreList, models.ScoreStudentRes{
						Name:      v.Name,
						UpdatedAt: info.UpdatedAt,
						ScoreFull: v.ScoreFull,
						ScoreGet:  info.ScoreGet,
						Status:    info.Status,
					})
					break
				}
			}
			// scoreList = append(scoreList, data)
		}
		if len(scoreList) == 0 {
			log.Println("student id", util.ErrNotFound)
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "score "+util.ErrNotFound.Error())
		}
		scoreRes["score_list"] = scoreList
	} else {
		log.Println("role", util.ErrValueInvalid)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "role"+util.ErrValueInvalid.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", scoreRes)
}

func (s *scoreController) GetScoreDataByCourseIdAndNameSore(c *fiber.Ctx) error {
	role := c.Query("role")
	profileId := c.Query("id")
	courseId, err := util.CheckStringData(c.Query("course_id"), "course_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find score of course id:", courseId)
	course, err := s.courseRepo.GetCourseById(courseId)
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

	name, err := util.CheckStringData(c.Query("name"), "name")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find score name:", name)

	score, err := s.scoreRepository.GetScoreByFilter(bson.M{"course_id": courseId, "name": name})
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	var scoreRes interface{}
	if role == "teacher" {
		if profileId != course.InstructorId {
			log.Println("not permiistion")
			return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, "not permission")
		}
		scoreRes = score
		// for _, v := range score.ScoreInformation {
		// 	scoreRes = append(scoreRes, v)
		// }
	} else {
		log.Println("role", util.ErrValueInvalid)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "role"+util.ErrValueInvalid.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"data_list": scoreRes,
	})
}

func createScoreinformation(studentIdList []string, t string) []models.ScoreInformation {

	var res []models.ScoreInformation
	for _, s := range studentIdList {
		res = append(res, models.ScoreInformation{
			StudentId: s,
			UpdatedAt: t,
			Status:    "not",
		})
	}

	return res
}
