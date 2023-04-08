package main

import (
	"log"
	"os"
	"school-notification-backend/controller"
	"school-notification-backend/db"
	"school-notification-backend/repository"
	"school-notification-backend/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// ยังไม่ได้เช็ค instructor id in subject
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

	// information
	informationRepository := repository.NewInformationRepository(conn)
	informationController := controller.NewInformationController(informationRepository)
	informationRoutes := routes.NewInformationRoutes(informationController)

	// subject
	subjectRepository := repository.NewSubjectRepository(conn)
	subjectController := controller.NewSubjectController(subjectRepository)
	subjectRoutes := routes.NewSubjectRoute(subjectController)

	// class
	classRepository := repository.NewClassRepository(conn)
	classController := controller.NewClassController(classRepository)
	classRoutes := routes.NewClassRoute(classController)

	// location
	locationRepository := repository.NewLocationRepository(conn)
	locationController := controller.NewLocationController(locationRepository)
	locationRoutes := routes.NewLocationRoute(locationController)

	staticRoutes := routes.NewStaticRoutes()

	route := fiber.New()

	route.Use(logger.New())
	route.Use(cors.New())

	informationRoutes.Install(route)
	subjectRoutes.Install(route)
	classRoutes.Install(route)
	locationRoutes.Install(route)
	staticRoutes.Install(route)

	route.Listen(":" + os.Getenv("APP_PORT"))
}
