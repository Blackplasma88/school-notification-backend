package controller

import (
	"log"
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

type CourseSummaryController interface {
	SummaryCourse(c *fiber.Ctx) error
	GetSummaryCourse(c *fiber.Ctx) error
	StudentGetSummaryCourse(c *fiber.Ctx) error
}

type courseSummaryController struct {
	courseSummaryRepo   repository.CourseSummaryRepository
	courseRepo          repository.CourseRepository
	scoreRepository     repository.ScoreRepository
	checkNameRepository repository.CheckNameRepository
	userRepo            repository.UsersRepository
	profileRepo         repository.ProfileRepository
}

func NewCourseSummaryController(courseSummaryRepo repository.CourseSummaryRepository, courseRepo repository.CourseRepository, scoreRepository repository.ScoreRepository, checkNameRepository repository.CheckNameRepository, userRepo repository.UsersRepository, profileRepo repository.ProfileRepository) CourseSummaryController {
	return &courseSummaryController{courseSummaryRepo: courseSummaryRepo, courseRepo: courseRepo, scoreRepository: scoreRepository, checkNameRepository: checkNameRepository, userRepo: userRepo, profileRepo: profileRepo}
}

func (cs *courseSummaryController) GetSummaryCourse(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cs.userRepo, []string{"all"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	courseId, err := util.CheckStringData(c.Query("course_id"), "course_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find course summary by course id:", courseId)

	role, err := util.CheckStringData(c.Query("role"), "role")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("role:", role)

	courseSum, err := cs.courseSummaryRepo.GetByFilter(bson.M{"course_id": courseId})
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	var res interface{}
	if role == "teacher" {
		res = courseSum.StudentData
	} else if role == "student" {
		studentId, err := util.CheckStringData(c.Query("student_id"), "student_id")
		if err != nil {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		log.Println("find course summary by student id:", studentId)

		check := true
		for _, v := range courseSum.StudentData {
			if v.StudentId == studentId {
				res = v
				check = false
				break
			}
		}
		if check {
			log.Println("not found")
			return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, util.ErrNotFound.Error())
		}

	} else {
		log.Println("role", util.ErrValueInvalid)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "role"+util.ErrValueInvalid.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"course_summary": res,
	})
}

func (cs *courseSummaryController) StudentGetSummaryCourse(c *fiber.Ctx) error {
	user, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cs.userRepo, []string{"student"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	year, err := util.CheckStringData(c.Query("year"), "year")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("year:", year)

	term, err := util.CheckStringData(c.Query("term"), "term")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("term:", term)

	p, err := cs.profileRepo.GetProfileById(bson.M{"profile_id": user.ProfileId, "role": user.Role}, user.Role)
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

	profile, _ := p.(models.ProfileStudent)

	index := -1
	for i, c := range profile.TermScore {
		if c.Year == year && c.Term == term {
			index = i
			break
		}
	}

	if index == -1 {
		log.Println("year or term not found")
		return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
	}

	courses := []*models.Course{}
	for _, v := range profile.TermScore[index].CourseList {
		course, err := cs.courseRepo.GetCourseById(v.Id.Hex())
		if err != nil {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
		}

		courses = append(courses, course)
	}

	var courseSumList []*models.CourseSummary
	for _, v := range courses {
		courseSum, err := cs.courseSummaryRepo.GetByFilter(bson.M{"course_id": v.Id.Hex()})
		if err != nil {
			log.Println(err)
			continue
		}

		courseSumList = append(courseSumList, courseSum)

	}

	dataList := []models.StudentDataRes{}
	for _, d := range courseSumList {
		for _, data := range d.StudentData {
			if data.StudentId == user.ProfileId {
				course, err := cs.courseRepo.GetCourseById(d.CourseId)
				if err != nil {
					log.Println(err)
					if err == mongo.ErrNoDocuments {
						return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
					}
					return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
				}
				dataRes := models.StudentDataRes{
					CourseName:           course.Name,
					ScoreWorkGet:         data.ScoreWorkGet,
					ScoreWorkFull:        data.ScoreWorkFull,
					ScoreMidGet:          data.ScoreMidGet,
					ScoreMidFull:         data.ScoreMidFull,
					ScoreFinalGet:        data.ScoreFinalGet,
					ScoreFinaFull:        data.ScoreFinaFull,
					Grade:                data.Grade,
					AllDateCount:         data.AllDateCount,
					CheckNameAttendCount: data.CheckNameAttendCount,
					CheckNameAbsentCount: data.CheckNameAbsentCount,
					CheckNameLateCount:   data.CheckNameLateCount,
				}
				dataList = append(dataList, dataRes)
				break
			}
		}
	}

	if len(dataList) == 0 {
		log.Println(util.ErrNotFound)
		return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"course_summary": dataList,
	})
}

func (cs *courseSummaryController) SummaryCourse(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], cs.userRepo, []string{"admin", "teacher"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	req := models.CourseSummaryRequest{}
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

	course, err := cs.courseRepo.GetCourseById(courseId)
	if err != nil {
		log.Println(err)
		if err.Error() == "mongo: no documents in result" {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, "course_id "+util.ErrNotFound.Error())
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

	log.Println("get scores")
	scores, err := cs.scoreRepository.GetByFilterAll(bson.M{"course_id": courseId})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		// if err == mongo.ErrNoDocuments {
		// 	return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		// }
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	log.Println("get check name list")
	checkNameList, err := cs.checkNameRepository.GetByFilterAll(bson.M{"course_id": courseId})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		// if err == mongo.ErrNoDocuments {
		// 	return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		// }
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	log.Println("check summary")
	courseSum, err := cs.courseSummaryRepo.GetByFilter(bson.M{"course_id": courseId})
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
		// if err == mongo.ErrNoDocuments {
		// 	return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		// }
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	t := time.Now().Format(time.RFC3339)
	var courseSummary models.CourseSummary
	if courseSum == nil {
		courseSummary = models.CourseSummary{
			Id:        primitive.NewObjectID(),
			CreatedAt: t,
			UpdatedAt: t,
			CourseId:  courseId,
		}
	} else {
		courseSummary = *courseSum
		courseSummary.StudentData = nil
		courseSummary.UpdatedAt = t
	}

	for _, studentId := range course.StudentIdList {

		scoreWorkGet := 0.0
		scoreWorkFull := 0.0
		scoreMidGet := 0.0
		scoreMidFull := 0.0
		scoreFinalGet := 0.0
		scoreFinalFull := 0.0
		for _, s := range scores {
			if s.Type == "work" {
				scoreWorkFull += s.ScoreFull
			} else if s.Type == "midterm" {
				scoreMidFull += s.ScoreFull
			} else if s.Type == "final" {
				scoreFinalFull += s.ScoreFull
			} else {
				log.Println(util.ErrTypeInvalid)
				return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrTypeInvalid.Error())
			}
			for _, sData := range s.ScoreInformation {
				if sData.StudentId == studentId {
					if s.Type == "work" {
						if sData.ScoreGet != nil {
							scoreWorkGet += *sData.ScoreGet
						}
						break
					} else if s.Type == "midterm" {
						if sData.ScoreGet != nil {
							scoreMidGet += *sData.ScoreGet
						}
						break
					} else if s.Type == "final" {
						if sData.ScoreGet != nil {
							scoreFinalGet += *sData.ScoreGet
						}
						break
					} else {
						log.Println(util.ErrTypeInvalid)
						return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrTypeInvalid.Error())
					}
				}
			}
		}

		totalScoreWork := 100.0 - scoreMidFull - scoreFinalFull
		// scoreWorkGet = (totalScoreWork * scoreWorkGet) / scoreWorkFull
		scoreWorkGet = (totalScoreWork * scoreWorkGet) / scoreWorkFull

		totalScore := scoreWorkGet + scoreMidGet + scoreFinalGet
		grade := 0.0
		if totalScore >= 80 {
			grade = 4
		} else if totalScore >= 75 {
			grade = 3.5
		} else if totalScore >= 70 {
			grade = 3
		} else if totalScore >= 65 {
			grade = 2.5
		} else if totalScore >= 60 {
			grade = 2
		} else if totalScore >= 55 {
			grade = 1.5
		} else if totalScore >= 50 {
			grade = 1
		} else {
			grade = 0
		}

		totalDate := 0
		totalDateAttend := 0
		totalDateAbsent := 0
		totalDateLate := 0
		for _, checkName := range checkNameList {
			totalDate++
			for _, cData := range checkName.CheckNameData {
				if cData.StudentId == studentId {
					if cData.Status == "attend" {
						totalDateAttend++
						break
					} else if cData.Status == "absent" {
						totalDateAbsent++
						break
					} else if cData.Status == "late" || cData.Status == "" { // ต้องไม่ได้
						totalDateLate++
						break
					} else {
						log.Println(util.ErrTypeInvalid)
						return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrTypeInvalid.Error())
					}
				}
			}
		}

		studentData := models.StudentData{
			StudentId:            studentId,
			ScoreWorkGet:         scoreWorkGet,
			ScoreWorkFull:        totalScoreWork,
			ScoreMidGet:          scoreMidGet,
			ScoreMidFull:         scoreMidFull,
			ScoreFinalGet:        scoreFinalGet,
			ScoreFinaFull:        scoreFinalFull,
			Grade:                grade,
			AllDateCount:         totalDate,
			CheckNameAttendCount: totalDateAttend,
			CheckNameAbsentCount: totalDateAbsent,
			CheckNameLateCount:   totalDateLate,
		}

		courseSummary.StudentData = append(courseSummary.StudentData, studentData)
	}

	if courseSum == nil {
		_, err = cs.courseSummaryRepo.Insert(&courseSummary)
		if err != nil {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
		}
	} else {
		_, err = cs.courseSummaryRepo.Update(&courseSummary)
		if err != nil {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
		}
	}

	course.Status = "summary"
	course.UpdatedAt = time.Now().Format(time.RFC3339)
	_, err = cs.courseRepo.Update(course)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create courseSummary success", map[string]interface{}{
		"course_summary_id": courseSummary.Id,
	})
}
