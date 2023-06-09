package main

import (
	"log"
	"os"
	"school-notification-backend/controller"
	"school-notification-backend/db"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"school-notification-backend/routes"
	"school-notification-backend/security"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	// load env
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Some error occured. Err: %s", err)
	}

	// set local location time
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}

func main() {
	conn := db.NewConnection()
	defer conn.Close()

	// school data
	schoolDataRepository := repository.NewSchoolDataRepository(conn)
	userRepository := repository.NewUsersRepository(conn)

	// information
	informationRepository := repository.NewInformationRepository(conn)
	informationController := controller.NewInformationController(informationRepository, userRepository)
	informationRoutes := routes.NewInformationRoutes(informationController)

	// class
	classRepository := repository.NewClassRepository(conn)

	// location
	locationRepository := repository.NewLocationRepository(conn)
	locationController := controller.NewLocationController(locationRepository, userRepository)
	locationRoutes := routes.NewLocationRoute(locationController)

	// profile
	profileRepository := repository.NewProfileRepository(conn)

	// subject
	subjectRepository := repository.NewSubjectRepository(conn)
	subjectController := controller.NewSubjectController(subjectRepository, schoolDataRepository, profileRepository, userRepository)
	subjectRoutes := routes.NewSubjectRoute(subjectController)

	// course
	courseRepository := repository.NewCoursesRepository(conn)
	courseSummaryRepository := repository.NewCourseSummaryRepository(conn)
	courseController := controller.NewCourseController(courseRepository, subjectRepository, schoolDataRepository, locationRepository, classRepository, profileRepository, courseSummaryRepository, userRepository)
	courseRoutes := routes.NewCourseRoute(courseController)

	// score
	scoreRepository := repository.NewScoreRepository(conn)
	scoreController := controller.NewScoreController(scoreRepository, courseRepository, userRepository)
	scoreRoutes := routes.NewScoreRoute(scoreController)

	// check name
	checkNameRepository := repository.NewCheckNameRepository(conn)
	checkNameController := controller.NewCheckNameController(checkNameRepository, courseRepository, userRepository)
	checkNameRoutes := routes.NewCheckNameRoute(checkNameController)

	// course summary
	courseSummaryController := controller.NewCourseSummaryController(courseSummaryRepository, courseRepository, scoreRepository, checkNameRepository, userRepository, profileRepository)
	courseSummaryRoutes := routes.NewCourseSummaryRoute(courseSummaryController)

	schoolDataController := controller.NewSchoolDataController(schoolDataRepository, courseRepository, courseSummaryRepository, profileRepository, classRepository, locationRepository, userRepository)
	schoolDataRoutes := routes.NewSchoolDataRoute(schoolDataController)

	// auth
	authController := controller.NewAuthController(userRepository, profileRepository)
	authRoutes := routes.NewAuthRoutes(authController)

	// conversation
	conversationRepository := repository.NewConversationRepository(conn)
	conversationController := controller.NewConversationController(conversationRepository, profileRepository, userRepository)
	conversationRoutes := routes.NewConversationRoute(conversationController)

	// message
	messageRepository := repository.NewMessageRepository(conn)
	messageController := controller.NewMessageController(messageRepository, conversationRepository, userRepository)
	messageRoutes := routes.NewMessageRoute(messageController)

	faceDetectionRepository := repository.NewFaceDetectionRepository(conn)
	faceDetectionController := controller.NewFaceDetectionController(faceDetectionRepository, classRepository, userRepository)
	faceDetectionRoutes := routes.NewFaceDetectionRoute(faceDetectionController)

	profileController := controller.NewProfileController(profileRepository, classRepository, schoolDataRepository, userRepository, faceDetectionRepository)
	profileRoutes := routes.NewProfileRoute(profileController)

	classController := controller.NewClassController(classRepository, schoolDataRepository, profileRepository, userRepository, faceDetectionRepository)
	classRoutes := routes.NewClassRoute(classController)

	staticRoutes := routes.NewStaticRoutes()

	initFirstData(profileRepository, userRepository)
	route := fiber.New()

	route.Use(logger.New())
	route.Use(cors.New())

	schoolDataRoutes.Install(route)
	informationRoutes.Install(route)
	subjectRoutes.Install(route)
	classRoutes.Install(route)
	locationRoutes.Install(route)
	profileRoutes.Install(route)
	courseRoutes.Install(route)
	scoreRoutes.Install(route)
	checkNameRoutes.Install(route)
	courseSummaryRoutes.Install(route)
	authRoutes.Install(route)
	conversationRoutes.Install(route)
	messageRoutes.Install(route)
	faceDetectionRoutes.Install(route)
	staticRoutes.Install(route)

	route.Listen(":" + os.Getenv("APP_PORT"))
}

func initFirstData(profileRepo repository.ProfileRepository, userRepo repository.UsersRepository) {
	err := profileRepo.GetProfileByFilterForCheckExists(bson.M{
		"profile_id": "admin1",
		"role":       "admin",
	})
	if err != nil && err.Error() == "mongo: no documents in result" {
		profileAdmin := models.ProfileAdmin{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			ProfileId: "admin1",
			Name:      "admin1",
			Role:      "admin",
		}
		_, err = profileRepo.Insert(profileAdmin)
		if err != nil {
			panic(err)
		}

		password, err := security.EncryptPassword("admin1")
		if err != nil {
			panic(err)
		}

		// sign up
		id := profileAdmin.Id
		user := models.User{
			Id:        primitive.NewObjectID(),
			CreatedAt: time.Now().Format(time.RFC3339),
			Username:  "admin1",
			Password:  password,
			ProfileId: "admin1",
			Role:      "admin",
			UserId:    id.Hex(),
		}

		_, err = userRepo.InsertUser(&user)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
}
