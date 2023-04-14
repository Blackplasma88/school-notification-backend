package controller

import (
	"log"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"school-notification-backend/util"
	"sort"
	"strconv"

	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SchoolDataController interface {
	AddYearAndTerm(c *fiber.Ctx) error
	AddSubjectCategory(c *fiber.Ctx) error
	UpdateSchoolData(c *fiber.Ctx) error
	GetSchoolDataAll(c *fiber.Ctx) error
	GetSubjectCategory(c *fiber.Ctx) error
	GetSchoolDataById(c *fiber.Ctx) error
	EndTerm(c *fiber.Ctx) error
}

type schoolDataController struct {
	schoolDataRepository repository.SchoolDataRepository
	courseRepo           repository.CourseRepository
	courseSummaryRepo    repository.CourseSummaryRepository
	profileRepo          repository.ProfileRepository
	classRepo            repository.ClassRepository
	locationRepo         repository.LocationRepository
}

func NewSchoolDataController(schoolDataRepository repository.SchoolDataRepository, courseRepo repository.CourseRepository, courseSummaryRepo repository.CourseSummaryRepository, profileRepo repository.ProfileRepository, classRepo repository.ClassRepository, locationRepo repository.LocationRepository) SchoolDataController {
	return &schoolDataController{schoolDataRepository: schoolDataRepository, courseRepo: courseRepo, courseSummaryRepo: courseSummaryRepo, profileRepo: profileRepo, classRepo: classRepo, locationRepo: locationRepo}
}

func (s *schoolDataController) AddYearAndTerm(c *fiber.Ctx) error {

	req := models.SchoolDataRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	year, err := util.CheckStringData(req.Year, "year")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("year:", year)

	term, err := util.CheckStringData(req.Term, "term")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("term:", term)

	dataList, err := s.schoolDataRepository.GetByFilterAll(bson.M{"type": "YearAndTerm"})
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(dataList) != 0 {
		if *dataList[len(dataList)-1].Status == false {
			log.Println("school year and term does not finish")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "school year and term does not finish")
		}
	}

	status := false
	dataNew := &models.SchoolData{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Type:      "YearAndTerm",
		Year:      &year,
		Term:      &term,
		Status:    &status,
	}

	_, err = s.schoolDataRepository.Insert(dataNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create data success", map[string]interface{}{
		"school_data_id": dataNew.Id,
	})
}

func (s *schoolDataController) AddSubjectCategory(c *fiber.Ctx) error {

	req := models.SchoolDataRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	category, err := util.CheckStringData(req.Category, "category")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("category:", category)

	dataList, err := s.schoolDataRepository.GetByFilterAll(bson.M{"type": "SubjectCategory"})
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(dataList) != 0 {
		for _, v := range dataList {
			if *v.SubjectCategory == category {
				log.Println("category has been addded")
				return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "category has been addded")
			}
		}
	}
	dataNew := &models.SchoolData{
		Id:              primitive.NewObjectID(),
		CreatedAt:       time.Now().Format(time.RFC3339),
		UpdatedAt:       time.Now().Format(time.RFC3339),
		Type:            "SubjectCategory",
		SubjectCategory: &category,
	}

	_, err = s.schoolDataRepository.Insert(dataNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create data success", map[string]interface{}{
		"school_data_id": dataNew.Id,
	})
}

func (s *schoolDataController) UpdateSchoolData(c *fiber.Ctx) error {

	req := models.SchoolDataRequest{}
	err := c.BodyParser(&req)
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

	data, err := s.schoolDataRepository.GetById(id)
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

	if *data.Status == true {
		log.Println("term does ended")
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, "term does ended")
	}

	data.UpdatedAt = time.Now().Format(time.RFC3339)
	if req.Status == nil {
		log.Println(util.ErrRequireParameter, "status")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"status")
	}

	// data.YearAndTerm.Status = *req.Status
	// if data.Status == true {

	// }

	result, err := s.schoolDataRepository.Update(data)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "update data success", map[string]interface{}{
		"school_data_id": data.Id,
		"update_count":   result.ModifiedCount,
	})
}

func (s *schoolDataController) GetSchoolDataAll(c *fiber.Ctx) error {

	data, err := s.schoolDataRepository.GetAll()
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(data) == 0 {
		log.Println("data not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"school_data": data,
	})
}

func (s *schoolDataController) GetSubjectCategory(c *fiber.Ctx) error {

	data, err := s.schoolDataRepository.GetByFilterAll(bson.M{"type": "SubjectCategory"})
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(data) == 0 {
		log.Println("data not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"school_data": data,
	})
}

func (s *schoolDataController) GetSchoolDataById(c *fiber.Ctx) error {
	id, err := util.CheckStringData(c.Query("id"), "id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find subject id:", id)

	data, err := s.schoolDataRepository.GetById(id)
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
		"school_data": data,
	})
}

func (s *schoolDataController) EndTerm(c *fiber.Ctx) error {

	dataList, err := s.schoolDataRepository.GetByFilterAll(bson.M{"type": "YearAndTerm"})
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

	sort.Slice(dataList, func(i, j int) bool {
		return dataList[i].CreatedAt > dataList[j].CreatedAt
	})

	// data := dataList[len(dataList)-1]
	data := dataList[0]

	if data.Status == nil || *data.Status == true {
		log.Println("school data invalid")
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, "school data invalid")
	}

	*data.Status = true

	if data.Term == nil || data.Year == nil {
		log.Println("school data invalid")
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, "school data invalid")
	}

	// check and finish course
	log.Println("get course list in term")
	courseList, err := s.courseRepo.GetCourseAllByFilter(bson.M{"year": *data.Year, "term": *data.Term})
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
		// if err.Error() == "mongo: no documents in result" {
		// 	return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		// }
		if err.Error() == "Id is not primitive objectID" {
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	for _, v := range courseList {
		if v.Status != "summary" && v.Status != "finish" {
			log.Println("have course does not finish")
			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, "have course does not finish")
		}
	}

	for _, cl := range courseList {
		if cl.Status == "summary" {
			log.Println("finish course:", cl.Id)
			courseSum, err := s.courseSummaryRepo.GetByFilter(bson.M{"course_id": cl.Id.Hex()})
			if err != nil && err.Error() != "mongo: no documents in result" {
				log.Println(err)
				if err == mongo.ErrNoDocuments {
					return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
				}
				return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
			}

			for _, sData := range courseSum.StudentData {
				p, err := s.profileRepo.GetProfileById(bson.M{"profile_id": sData.StudentId, "role": "student"}, "student")
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

				for i, t := range profile.TermScore {
					if t.Year == cl.Year && t.Term == cl.Term {
						totalTermGrade := 0.0
						for j, courseList := range t.CourseList {
							if courseList.Id == cl.Id {
								// profile.TermScore[i].CourseList[j].CreatedAt = time.Now().Format(time.RFC3339)
								profile.TermScore[i].CourseList[j].Grade = sData.Grade
								profile.TermScore[i].CourseList[j].ScoreWorkGet = sData.ScoreWorkGet
								profile.TermScore[i].CourseList[j].ScoreWorkFull = sData.ScoreWorkFull
								profile.TermScore[i].CourseList[j].ScoreMidGet = sData.ScoreMidGet
								profile.TermScore[i].CourseList[j].ScoreMidFull = sData.ScoreMidFull
								profile.TermScore[i].CourseList[j].ScoreFinalGet = sData.ScoreFinalGet
								profile.TermScore[i].CourseList[j].ScoreFinaFull = sData.ScoreFinaFull
								profile.TermScore[i].CourseList[j].Credit = cl.Credit
								profile.TermScore[i].CourseList[j].AllDateCount = sData.AllDateCount
								profile.TermScore[i].CourseList[j].CheckNameAttendCount = sData.CheckNameAttendCount
								profile.TermScore[i].CourseList[j].CheckNameAbsentCount = sData.CheckNameAbsentCount
								profile.TermScore[i].CourseList[j].CheckNameLateCount = sData.CheckNameLateCount

								profile.TermScore[i].TermCredit += cl.Credit
								profile.AllCredit += cl.Credit
								break
							}
						}

						for _, courseList := range profile.TermScore[i].CourseList {
							totalTermGrade += courseList.Grade * float64(courseList.Credit)
						}
						// เกรด * หน่วยกิต นำมารวมกัน หารด้วยหน่วยกิตทั้งหมด
						profile.TermScore[i].GPA = totalTermGrade / float64(profile.TermScore[i].TermCredit)
					}

					totalGrade := 0.0
					for _, t := range profile.TermScore {
						totalGrade += t.GPA * float64(t.TermCredit)
					}

					profile.GPA = totalGrade / float64(profile.AllCredit)

					_, err = s.profileRepo.Update(profile.Id, profile)
					if err != nil {
						log.Println(err)
						return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
					}

					break
				}
			}

			cl.Status = "finish"
			cl.UpdatedAt = time.Now().Format(time.RFC3339)
			_, err = s.courseRepo.Update(cl)
			if err != nil {
				log.Println(err)
				return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
			}

			// clear location
			log.Println("get location in course:", cl.Name)
			location, err := s.locationRepo.GetLocationById(cl.LocationId.Hex())
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

			for _, dt := range cl.DateTime {
				for i, slot := range location.Slot {
					if slot.Day == dt.Day {
						for _, t := range dt.Time {
							for j, ts := range slot.TimeSlot {
								if ts.Time == t && ts.CourseId != nil && *ts.CourseId == cl.Id {
									location.Slot[i].TimeSlot[j].Status = false
									location.Slot[i].TimeSlot[j].CourseId = nil
									break
								}
							}
						}
						break
					}
				}
			}

			_, err = s.locationRepo.Update(location)
			if err != nil {
				log.Println(err)
				return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
			}
		}
	}

	// add school data new term
	status := false
	var dataNew models.SchoolData
	if *data.Term == "2" {
		year, _ := strconv.Atoi(*data.Year)
		year++
		yearStr := strconv.Itoa(year)
		term := "1"
		dataNew = models.SchoolData{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			Type:      "YearAndTerm",
			Year:      &yearStr,
			Term:      &term,
			Status:    &status,
		}
	} else if *data.Term == "1" {
		term := "2"
		dataNew = models.SchoolData{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			Type:      "YearAndTerm",
			Year:      data.Year,
			Term:      &term,
			Status:    &status,
		}
	} else {
		log.Println("school data invalid")
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, "school data invalid")
	}

	// update class data
	log.Println("get all class")
	classes, err := s.classRepo.GetClassByFilterAll(bson.M{"status": false})
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	for _, class := range classes {
		if class.Year == *dataNew.Year && class.Term == *dataNew.Term {
			continue
		}
		if class.Status != true {
			class.Year = *dataNew.Year
			class.Term = *dataNew.Term
		}
		if class.Term == "1" {
			if class.ClassYear == "1" {
				class.ClassYear = "2"
			} else if class.ClassYear == "2" {
				class.ClassYear = "3"
			} else if class.ClassYear == "3" {
				class.ClassYear = "4"
			} else if class.ClassYear == "4" {
				class.ClassYear = "5"
			} else if class.ClassYear == "5" {
				class.ClassYear = "6"
			} else if class.ClassYear == "6" {
				class.Status = true
			}
		}

		class.Slot = createTimeSlot()

		_, err = s.classRepo.Update(class)
		if err != nil {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
		}

		// update profile student term
		log.Println("get student in class:", class.Id)
		if class.Status == false {
			for _, sid := range class.StudentIdList {
				filter := bson.M{
					"profile_id": sid,
					"role":       "student",
				}
				p, err := s.profileRepo.GetProfileById(filter, "student")
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

				check := false
				for _, ts := range profile.TermScore {
					if ts.Year == class.Year && ts.Term == class.Term {
						check = true
						break
					}
				}
				if check {
					continue
				}

				profile.TermScore = append(profile.TermScore, models.TermScore{
					Year: class.Year,
					Term: class.Term,
				})

				_, err = s.profileRepo.Update(profile.Id, profile)
				if err != nil {
					log.Println(err)
					return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
				}
			}
		}
	}

	// update teacher student term
	log.Println("get all teacher")
	p, err := s.profileRepo.GetAll("teacher")
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	for _, v := range p {
		profileTeacher, _ := v.(*models.ProfileTeacher)
		check := false
		for _, ctl := range profileTeacher.CourseTeachesList {
			if ctl.Year == *dataNew.Year && ctl.Term == *dataNew.Term {
				check = true
				break
			}
		}
		if check {
			continue
		}

		profileTeacher.Slot = createTimeSlot()

		profileTeacher.CourseTeachesList = append(profileTeacher.CourseTeachesList, models.CourseTeachesList{
			Year: *dataNew.Year,
			Term: *dataNew.Term,
		})

		_, err = s.profileRepo.Update(profileTeacher.Id, profileTeacher)
		if err != nil {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
		}
	}

	_, err = s.schoolDataRepository.Insert(dataNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "end term success", map[string]interface{}{
		"school_data_id": data.Id,
	})
}
