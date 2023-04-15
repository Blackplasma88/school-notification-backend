package controller

import (
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

type CourseController interface {
	CreateCourse(c *fiber.Ctx) error
	ChangeCourseStatus(c *fiber.Ctx) error
	GetCourseByYearAndTerm(c *fiber.Ctx) error
	FinishCourse(c *fiber.Ctx) error
	GetCourseById(c *fiber.Ctx) error
}

type courseController struct {
	courseRepo           repository.CourseRepository
	subjectRepository    repository.SubjectRepository
	schoolDataRepository repository.SchoolDataRepository
	locationRepo         repository.LocationRepository
	classRepo            repository.ClassRepository
	profileRepo          repository.ProfileRepository
	courseSummaryRepo    repository.CourseSummaryRepository
}

func NewCourseController(courseRepo repository.CourseRepository, subjectRepository repository.SubjectRepository, schoolDataRepository repository.SchoolDataRepository, locationRepo repository.LocationRepository, classRepo repository.ClassRepository, profileRepo repository.ProfileRepository, courseSummaryRepo repository.CourseSummaryRepository) CourseController {
	return &courseController{courseRepo: courseRepo, subjectRepository: subjectRepository, schoolDataRepository: schoolDataRepository, locationRepo: locationRepo, classRepo: classRepo, profileRepo: profileRepo, courseSummaryRepo: courseSummaryRepo}
}

func (cc *courseController) CreateCourse(c *fiber.Ctx) error {

	req := models.CourseRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	dataList, err := cc.schoolDataRepository.GetByFilterAll(bson.M{"type": "YearAndTerm"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	sort.Slice(dataList, func(i, j int) bool {
		return dataList[i].CreatedAt > dataList[j].CreatedAt
	})

	if *dataList[len(dataList)-1].Status == true {
		log.Println("school data invalid")
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, "school data invalid")
	}

	year := *dataList[0].Year
	term := *dataList[0].Term

	subjectId, err := util.CheckStringData(req.SubjectId, "subject_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("subject id:", subjectId)
	instructorId, err := util.CheckStringData(req.InstructorId, "instructor_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("instructor id:", instructorId)

	subject, err := cc.subjectRepository.GetSubjectByFilter(bson.M{"subject_id": subjectId})
	if err != nil {
		log.Println(err)
		if err.Error() == "mongo: no documents in result" {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, "subject_id "+util.ErrNotFound.Error())
		}
		if err.Error() == "Id is not primitive objectID" {
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	check := true
	for _, v := range subject.InstructorId {
		if v == instructorId {
			check = false
			break
		}
	}

	if check {
		log.Println("instructor id not found in subject instructor")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "instructor id not found in subject instructor")
	}

	classId, err := util.CheckStringData(req.ClassId, "class_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("create profile in class id:", classId)

	class, err := cc.classRepo.GetClassById(classId)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		if err.Error() == "Id is not primitive objectID" {
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if class.Status == true {
		log.Println("class did finish")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "class did finish")
	}

	locationId, err := util.CheckStringData(req.LocationId, "location_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("location id:", locationId)

	location, err := cc.locationRepo.GetLocationByFilter(bson.M{"location_id": req.LocationId})
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

	// check create course again
	_, err = cc.courseRepo.GetCourseByFilter(bson.M{"subject_id": subjectId, "class_id": classId})
	if err == nil {
		log.Println("course data subject and class already exists")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "course data subject and class"+util.ErrValueAlreadyExists.Error())
	}
	if err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	courseNew := &models.Course{
		Id:              primitive.NewObjectID(),
		CreatedAt:       time.Now().Format(time.RFC3339),
		UpdatedAt:       time.Now().Format(time.RFC3339),
		Status:          "create",
		SubjectId:       subjectId,
		InstructorId:    instructorId,
		Year:            year,
		Term:            term,
		Name:            subject.Name + "-" + year + "-" + term,
		Credit:          subject.Credit,
		LocationId:      &location.Id,
		NumberOfStudent: class.NumberOfStudent,
		StudentIdList:   class.StudentIdList,
		ClassId:         &class.Id,
		ClassYear:       class.Year,
		ClassRoom:       class.ClassRoom,
	}

	if len(req.DateTime) == 0 {
		log.Println(util.ErrRequireParameter.Error() + "date_time")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"date_time")
	}

	filter := bson.M{
		"profile_id": instructorId,
		"role":       "teacher",
	}

	p, err := cc.profileRepo.GetProfileById(filter, "teacher")
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

	for i, v := range profile.CourseTeachesList {
		if v.Term == courseNew.Term && v.Year == courseNew.Year {
			profile.CourseTeachesList[i].CourseIdList = append(profile.CourseTeachesList[i].CourseIdList, courseNew.Id)
			break
		}
	}

	// check date time in class
	for di, dt := range req.DateTime {
		req.DateTime[di].Day, err = util.CheckStringData(dt.Day, "day")
		if err != nil {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		}
		dt.Day = req.DateTime[di].Day
		log.Println("day:", dt.Day)
		checkDay := true
		for i, slot := range class.Slot {
			if dt.Day == slot.Day {
				check := true
				if len(dt.Time) == 0 {
					log.Println(util.ErrRequireParameter.Error() + "time")
					return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"time")
				}
				for dti, t := range dt.Time {
					for j, ts := range slot.TimeSlot {
						req.DateTime[di].Time[dti], err = util.CheckStringData(t, "time")
						if err != nil {
							log.Println(err)
							return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
						}
						t = req.DateTime[di].Time[dti]
						log.Println("time:", t)
						if ts.Time == t {
							if ts.Status == true {
								log.Println("date time is used in class")
								return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "date time is used in this class")
							}
							class.Slot[i].TimeSlot[j].Status = true
							class.Slot[i].TimeSlot[j].CourseId = &courseNew.Id
							check = false
							break
						}
					}
					if check {
						log.Println("time", dt.Time, "is in valid")
						return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "time "+t+" is in valid")
					}
				}
				checkDay = false
				break
			}
		}
		if checkDay {
			log.Println("day", dt.Day, "is in valid")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "day "+dt.Day+" is in valid")
		}
	}

	// check date time in profile teacher
	for _, dt := range req.DateTime {
		for i, slot := range profile.Slot {
			if dt.Day == slot.Day {
				for _, t := range dt.Time {
					for j, ts := range slot.TimeSlot {
						if ts.Time == t {
							if ts.Status == true {
								log.Println("date time is used in teacher time lot")
								return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "date time is used in teacher time lot")
							}
							profile.Slot[i].TimeSlot[j].Status = true
							profile.Slot[i].TimeSlot[j].CourseId = &courseNew.Id
							break
						}
					}
				}
				break
			}
		}
	}

	// check date time in location and set
	for _, dt := range req.DateTime {
		// day, err := util.CheckStringData(dt.Day, "day")
		// if err != nil {
		// 	log.Println(err)
		// 	return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		// }
		// log.Println("day:", day)
		// checkDay := true
		for i, slot := range location.Slot {
			if dt.Day == slot.Day {
				// check := true
				// if len(dt.Time) == 0 {
				// 	log.Println(util.ErrRequireParameter.Error() + "time")
				// 	return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"time")
				// }
				for _, t := range dt.Time {
					for j, ts := range slot.TimeSlot {
						// ti, err := util.CheckStringData(t, "time")
						// if err != nil {
						// 	log.Println(err)
						// 	return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
						// }
						// log.Println("time:", ti)
						if ts.Time == t {
							if ts.Status == true {
								log.Println("date time is used in this location")
								return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "date time is used in this location")
							}
							location.Slot[i].TimeSlot[j].Status = true
							location.Slot[i].TimeSlot[j].CourseId = &courseNew.Id
							// check = false
							break
						}
					}
					// if check {
					// 	log.Println("time", dt.Time, "is in valid")
					// 	return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "time "+t+" is in valid")
					// }
				}
				// checkDay = false
				break
			}
		}
		// if checkDay {
		// 	log.Println("day", dt.Day, "is in valid")
		// 	return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "day "+dt.Day+" is in valid")
		// }
	}

	// update value
	courseNew.DateTime = req.DateTime

	_, err = cc.courseRepo.Insert(courseNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	_, err = cc.classRepo.Update(class)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	_, err = cc.profileRepo.Update(profile.Id, profile)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	for _, s := range class.StudentIdList {
		filter := bson.M{
			"profile_id": s,
			"role":       "student",
		}
		p, err := cc.profileRepo.GetProfileById(filter, "student")
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
		for i, v := range profile.TermScore {
			if v.Term == courseNew.Term && v.Year == courseNew.Year {
				profile.TermScore[i].CourseList = append(profile.TermScore[i].CourseList, models.CourseList{
					Id: courseNew.Id,
				})
				break
			}
		}

		_, err = cc.profileRepo.Update(profile.Id, profile)
		if err != nil {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
		}
	}

	_, err = cc.locationRepo.Update(location)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create course success", map[string]interface{}{
		"course_id": courseNew.Id,
	})
}

func (cc *courseController) ChangeCourseStatus(c *fiber.Ctx) error {
	req := models.CourseRequest{}
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

	course, err := cc.courseRepo.GetCourseById(id)
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

	event, err := util.CheckStringData(req.Event, "event")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("event:", event)

	if event == "ChangeToProgress" {
		if course.Status != "create" {
			log.Println("status invalid")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ReturnErrorStatusInvalid("course", "create").Error())
		}

		course.Status = "progress"
	} else if event == "ChangeReverseSummary" {
		if course.Status != "summary" {
			log.Println("status invalid")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ReturnErrorStatusInvalid("course", "summary").Error())
		}

		if course.SubjectId == "" {
			log.Println("subject id is nil")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "subject id is nil")
		}

		if course.InstructorId == "" {
			log.Println("instructor id is nil")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "instructor id is nil")
		}

		if course.LocationId == nil {
			log.Println("location id is nil")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "location id is nil")
		}

		if course.ClassId == nil {
			log.Println("class id is nil")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "class id is nil")
		}

		course.Status = "progress"

	} else if event == "ChangeToFinish" {
		if course.Status != "summary" {
			log.Println("status invalid")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ReturnErrorStatusInvalid("course", "final").Error())
		}

		// for _, v := range course.StudentSummary {
		// 	profile, err := cc.profileRepo.GetStudentProfileById(bson.M{"profile_id": v.StudentId, "role": "student"})
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 		// if err.Error() == "mongo: no documents in result" {
		// 		// 	return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		// 		// }
		// 		// if err.Error() == "Id is not primitive objectID" {
		// 		// 	return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		// 		// }
		// 		// return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
		// 	}

		// 	for i, t := range profile.TermScore {
		// 		if t.Year == course.CourseData.Year && t.Term == course.CourseData.Term {

		// 			totalTermGrade := 0.0
		// 			for j, courseList := range t.CourseList {
		// 				if courseList.CourseId == course.CourseId {
		// 					profile.TermScore[i].CourseList[j].CreatedAt = time.Now().Format(time.RFC3339)
		// 					profile.TermScore[i].CourseList[j].Grade = v.Grade
		// 					profile.TermScore[i].CourseList[j].ScoreGet = v.ScoreGet
		// 					profile.TermScore[i].CourseList[j].ScoreFull = v.ScoreFull
		// 					profile.TermScore[i].CourseList[j].Status = v.Status
		// 					profile.TermScore[i].CourseList[j].Credit = course.CourseData.Credit
		// 					profile.TermScore[i].TermCredit += course.CourseData.Credit
		// 					profile.AllCredit += course.CourseData.Credit
		// 					break
		// 				}
		// 			}

		// 			for _, courseList := range profile.TermScore[i].CourseList {
		// 				totalTermGrade += courseList.Grade * float64(courseList.Credit)
		// 			}
		// 			// เกรด * หน่วยกิต นำมารวมกัน หารด้วยหน่วยกิตทั้งหมด
		// 			profile.TermScore[i].GPA = totalTermGrade / float64(profile.TermScore[i].TermCredit)
		// 		}

		// 		totalGrade := 0.0
		// 		for _, t := range profile.TermScore {
		// 			totalGrade += t.GPA * float64(t.TermCredit)
		// 		}

		// 		profile.GPA = totalGrade / float64(profile.AllCredit)

		// 		_, err = cc.profileRepo.Update(profile.Id, profile)
		// 		if err != nil {
		// 			log.Println(err)
		// 			continue
		// 			// return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
		// 		}

		// 		break
		// 	}
		// }

		course.Status = "finish"
	} else {
		log.Println("event is invalid")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrEventInvalid.Error())
	}

	course.UpdatedAt = time.Now().Format(time.RFC3339)
	courseUpdate, err := cc.courseRepo.Update(course)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "update success", map[string]interface{}{
		"course_id":           course.Id,
		"course_update_count": courseUpdate.ModifiedCount,
	})
}

func (cc *courseController) GetCourseByYearAndTerm(c *fiber.Ctx) error {
	profileId, err := util.CheckStringData(c.Query("profile_id"), "profile_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find course of profile id", profileId)
	role, err := util.CheckStringData(c.Query("role"), "role")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find course of role", role)
	year, err := util.CheckStringData(c.Query("year"), "year")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find course of year", year)
	term, err := util.CheckStringData(c.Query("term"), "term")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find course of term", term)

	var coursesRes interface{}
	if role == "admin" {
		courses, err := cc.courseRepo.GetCourseAllByFilter(bson.M{"year": year, "term": term})
		if err != nil {
			log.Println(err)
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
		}

		coursesRes = courses
	} else if role == "teacher" {
		p, err := cc.profileRepo.GetProfileById(bson.M{"profile_id": profileId, "role": role}, role)
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

		index := -1
		for i, c := range profile.CourseTeachesList {
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
		for _, v := range profile.CourseTeachesList[index].CourseIdList {
			course, err := cc.courseRepo.GetCourseById(v.Hex())
			if err != nil {
				log.Println(err)
				return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
			}

			courses = append(courses, course)
		}

		coursesRes = courses
	} else if role == "student" {
		p, err := cc.profileRepo.GetProfileById(bson.M{"profile_id": profileId, "role": role}, role)
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
			course, err := cc.courseRepo.GetCourseById(v.Id.String())
			if err != nil {
				log.Println(err)
				return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
			}

			courses = append(courses, course)
		}

		coursesRes = courses
	} else {
		log.Println("role is invalid")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "role"+util.ErrValueInvalid.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"profile_id":  profileId,
		"role":        role,
		"course_list": coursesRes,
	})
}

func (cc *courseController) GetCourseById(c *fiber.Ctx) error {
	id, err := util.CheckStringData(c.Query("course_id"), "course_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find course id:", id)

	course, err := cc.courseRepo.GetCourseById(id)
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

	if course == nil {
		log.Println("course not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"course": course,
	})
}

func (cc *courseController) FinishCourse(c *fiber.Ctx) error {
	req := models.CourseRequest{}
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

	course, err := cc.courseRepo.GetCourseById(id)
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

	if course.Status != "summary" {
		log.Println("status invalid")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ReturnErrorStatusInvalid("course", "summary").Error())
	}

	log.Println("get course summary")
	courseSum, err := cc.courseSummaryRepo.GetByFilter(bson.M{"course_id": id})
	if err != nil && err.Error() != "mongo: no documents in result" {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	for _, sData := range courseSum.StudentData {
		p, err := cc.profileRepo.GetProfileById(bson.M{"profile_id": sData.StudentId, "role": "student"}, "student")
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
			if t.Year == course.Year && t.Term == course.Term {
				totalTermGrade := 0.0
				for j, courseList := range t.CourseList {
					if courseList.Id == course.Id {
						if profile.TermScore[i].CourseList[j].Credit != 0 {
							profile.TermScore[i].TermCredit -= profile.TermScore[i].CourseList[j].Credit
							profile.AllCredit -= profile.TermScore[i].CourseList[j].Credit
						}
						// profile.TermScore[i].CourseList[j].CreatedAt = time.Now().Format(time.RFC3339)
						profile.TermScore[i].CourseList[j].Grade = sData.Grade
						profile.TermScore[i].CourseList[j].ScoreWorkGet = sData.ScoreWorkGet
						profile.TermScore[i].CourseList[j].ScoreWorkFull = sData.ScoreWorkFull
						profile.TermScore[i].CourseList[j].ScoreMidGet = sData.ScoreMidGet
						profile.TermScore[i].CourseList[j].ScoreMidFull = sData.ScoreMidFull
						profile.TermScore[i].CourseList[j].ScoreFinalGet = sData.ScoreFinalGet
						profile.TermScore[i].CourseList[j].ScoreFinaFull = sData.ScoreFinaFull
						profile.TermScore[i].CourseList[j].Credit = course.Credit
						profile.TermScore[i].CourseList[j].AllDateCount = sData.AllDateCount
						profile.TermScore[i].CourseList[j].CheckNameAttendCount = sData.CheckNameAttendCount
						profile.TermScore[i].CourseList[j].CheckNameAbsentCount = sData.CheckNameAbsentCount
						profile.TermScore[i].CourseList[j].CheckNameLateCount = sData.CheckNameLateCount

						profile.TermScore[i].TermCredit += course.Credit
						profile.AllCredit += course.Credit
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

			_, err = cc.profileRepo.Update(profile.Id, profile)
			if err != nil {
				log.Println(err)
				return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
			}

			break
		}
	}

	course.Status = "finish"
	course.UpdatedAt = time.Now().Format(time.RFC3339)
	courseUpdate, err := cc.courseRepo.Update(course)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	// update location
	log.Println("get location")
	location, err := cc.locationRepo.GetLocationById(course.LocationId.Hex())
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

	for _, dt := range course.DateTime {
		for i, slot := range location.Slot {
			if slot.Day == dt.Day {
				for _, t := range dt.Time {
					for j, ts := range slot.TimeSlot {
						if ts.Time == t && ts.CourseId != nil && *ts.CourseId == course.Id {
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

	_, err = cc.locationRepo.Update(location)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "update success", map[string]interface{}{
		"course_id":           course.Id,
		"course_update_count": courseUpdate.ModifiedCount,
	})
}

// func createDateTime(dateTime []models.DateTime) {
// 	for _, dt := range dateTime {
// 		dayTime := ""
// 		times := strings.Split(dt.Time, "-")
// 		for _, t := range times {
// 			if t != "" {
// 				if len(dayTime) == 0 {
// 					dayTime += t + "-"
// 				} else {
// 					tmp := strings.Split(t, ":")
// 					if tmp[1] == "30" {
// 						tmp[1] = "00"

// 					}
// 				}
// 			}
// 		}

// 	}
// }
