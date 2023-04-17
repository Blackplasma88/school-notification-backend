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

type LocationController interface {
	CreateLocation(c *fiber.Ctx) error
	UpdateLocationData(c *fiber.Ctx) error
	GetLocationAll(c *fiber.Ctx) error
	GetLocationById(c *fiber.Ctx) error
}

type locationController struct {
	locationRepo repository.LocationRepository
	userRepo     repository.UsersRepository
}

func NewLocationController(locationRepo repository.LocationRepository, userRepo repository.UsersRepository) LocationController {
	return &locationController{locationRepo: locationRepo, userRepo: userRepo}
}

func (l *locationController) CreateLocation(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], l.userRepo, []string{"admin"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	req := models.LocationRequest{}
	err = c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	buildingName, err := util.CheckStringData(req.BuildingName, "building_name")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("building name:", buildingName)

	floor, err := util.CheckStringData(req.Floor, "floor")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("floor:", floor)

	room, err := util.CheckStringData(req.Room, "room")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("room:", room)

	_, err = l.locationRepo.GetLocationByFilter(bson.M{"building_name": buildingName, "floor": floor, "room": room})
	if err == nil {
		log.Println("location already exists")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "location"+util.ErrValueAlreadyExists.Error())
	}
	if err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	locationId := buildingName + "-" + floor + "-" + room
	log.Println("location id:", locationId)

	locationNew := &models.Location{
		Id:           primitive.NewObjectID(),
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
		LocationId:   locationId,
		BuildingName: buildingName,
		Floor:        floor,
		Room:         room,
		Status:       true,
		Slot:         createTimeSlot(),
	}

	_, err = l.locationRepo.Insert(locationNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create location success", map[string]interface{}{
		"location_id": locationId,
	})
}

func (l *locationController) UpdateLocationData(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], l.userRepo, []string{"admin"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	req := models.LocationRequest{}
	err = c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	locationId, err := util.CheckStringData(req.LocationId, "location_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("location id:", locationId)

	location, err := l.locationRepo.GetLocationByFilter(bson.M{"location_id": req.LocationId})
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

	buildingName, err := util.CheckStringData(req.BuildingName, "building_name")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("building name:", buildingName)

	floor, err := util.CheckStringData(req.Floor, "floor")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("floor:", floor)

	room, err := util.CheckStringData(req.Room, "room")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("room:", room)

	locationIdNew := location.BuildingName + "-" + location.Floor + "-" + location.Room

	if locationId != locationIdNew {
		_, err := l.locationRepo.GetLocationByFilter(bson.M{"location_id": location.LocationId})
		if err == nil {
			log.Println("location new already exists")
			return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "location"+util.ErrValueAlreadyExists.Error())
		}
		if err.Error() != "mongo: no documents in result" {
			log.Println(err)
			return util.ErrInternalServerError
		}
	}

	location.LocationId = locationIdNew
	location.BuildingName = buildingName
	location.Floor = floor
	location.Room = room
	if req.Status != nil {
		location.Status = *req.Status
	}
	location.UpdatedAt = time.Now().Format(time.RFC3339)

	locationUpdate, err := l.locationRepo.Update(location)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "update location success", map[string]interface{}{
		"location_id":  location.LocationId,
		"update_count": locationUpdate.ModifiedCount,
	})
}

func (l *locationController) GetLocationAll(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], l.userRepo, []string{"all"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	locations, err := l.locationRepo.GetAll()
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return util.ResponseNotSuccess(c, fiber.StatusNotFound, util.ErrNotFound.Error())
		}
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(locations) == 0 {
		log.Println("location not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"location_list": locations,
	})
}

func (l *locationController) GetLocationById(c *fiber.Ctx) error {
	_, err := security.CheckRoleFromToken(c.GetReqHeaders()["Authorization"], l.userRepo, []string{"all"})
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.ErrUnauthorized.Code, err.Error())
	}

	id, err := util.CheckStringData(c.Query("location_id"), "location_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find location id:", id)

	location, err := l.locationRepo.GetLocationById(id)
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

	if location == nil {
		log.Println("location not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"location": location,
	})
}

func createTimeSlot() []models.Slot {
	var slot []models.Slot
	day := []string{"monday", "tuesday", "wednesday", "thursday", "friday"}
	time := []string{"08:30", "09:00", "09:30", "10:00", "10:30", "11:00", "12:30", "13:00", "13:30", "14:00", "14:30", "15:00", "15:30", "16:00"}

	for _, d := range day {
		s := models.Slot{
			Day: d,
		}
		for _, t := range time {
			s.TimeSlot = append(s.TimeSlot, models.TimeSlot{
				Time:   t,
				Status: false,
			})
		}
		slot = append(slot, s)
	}

	return slot
}
