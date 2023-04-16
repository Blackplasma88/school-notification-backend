package controller

import (
	"errors"
	"fmt"
	"log"
	"school-notification-backend/models"
	"school-notification-backend/repository"
	"school-notification-backend/security"
	"school-notification-backend/util"

	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthController interface {
	SignUp(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	GetUserWithId(c *fiber.Ctx) error
}

type authController struct {
	userRepo    repository.UsersRepository
	profileRepo repository.ProfileRepository
}

func NewAuthController(userRepo repository.UsersRepository, profileRepo repository.ProfileRepository) AuthController {
	return &authController{userRepo: userRepo, profileRepo: profileRepo}
}

func (a *authController) SignUp(c *fiber.Ctx) error {

	var input models.UserRequest
	err := c.BodyParser(&input)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	if len(strings.TrimSpace(input.Username)) == 0 {
		log.Println("do not have parameter username")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"username")
	}
	input.Username = strings.TrimSpace(input.Username)
	log.Println("username:", input.Username)

	_, err = a.userRepo.GetByUsername(input.Username)
	if err == nil {
		log.Println("username already exists")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, "username"+util.ErrValueAlreadyExists.Error())
	}
	if err.Error() != "mongo: no documents in result" {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	if len(strings.TrimSpace(input.Password)) == 0 {
		log.Println("do not have parameter password")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"password")
	}
	input.Password = strings.TrimSpace(input.Password)
	log.Println("password:", input.Password)

	input.Password, err = security.EncryptPassword(input.Password)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}

	if len(strings.TrimSpace(input.Role)) == 0 {
		log.Println("do not have parameter role")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"role")
	}
	input.Role = strings.TrimSpace(input.Role)
	log.Println("role:", input.Role)

	if len(strings.TrimSpace(input.ProfileId)) == 0 {
		log.Println("do not have parameter profile id")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"profile_id")
	}
	input.ProfileId = strings.TrimSpace(input.ProfileId)
	log.Println("profile id:", input.ProfileId)

	if len(strings.TrimSpace(input.UserId)) == 0 {
		log.Println("do not have parameter user id")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"user_id")
	}
	input.UserId = strings.TrimSpace(input.UserId)
	log.Println("user id:", input.UserId)

	// err = a.profileRepo.GetProfileByFilterForCheckExists(bson.M{"profile_id": input.ProfileId, "role": input.Role})
	// if err != nil {
	// 	log.Println(err)
	// 	if err.Error() == "mongo: no documents in result" {
	// 		return util.ResponseNotSuccess(c, fiber.StatusNotFound, "profile_id"+util.ErrValueNotAlreadyExists.Error())
	// 	}
	// 	return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, err.Error())
	// }

	user := models.User{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now().Format(time.RFC3339),
		Username:  input.Username,
		Password:  input.Password,
		ProfileId: input.ProfileId,
		Role:      input.Role,
		UserId:    input.UserId,
	}

	result, err := a.userRepo.InsertUser(&user)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}

	log.Println("result:", result)

	return util.ResponseSuccess(c, fiber.StatusCreated, "create user success", map[string]interface{}{
		"user_id":  result.InsertedID,
		"username": user.Username,
	})
}

func (a *authController) SignIn(c *fiber.Ctx) error {
	var input models.UserRequest
	err := c.BodyParser(&input)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	if len(strings.TrimSpace(input.Username)) == 0 {
		log.Println("do not have parameter username")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"username")
	}
	input.Username = strings.TrimSpace(input.Username)
	log.Println("username:", input.Username)

	exists, err := a.userRepo.GetByUsername(input.Username)
	if err != nil {
		log.Println(input.Username, "signin failed")
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, errors.New("invalid credentials").Error())
	}

	if len(strings.TrimSpace(input.Password)) == 0 {
		log.Println("do not have parameter password")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"password")
	}
	input.Password = strings.TrimSpace(input.Password)
	log.Println("password:", input.Password)

	err = security.VerifyPassword(exists.Password, input.Password)
	if err != nil {
		log.Println(input.Username, "signin failed")
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, errors.New("invalid credentials").Error())
	}

	tokenStr, err := security.NewToken(exists.Id.Hex())
	if err != nil {
		log.Println(input.Username, "signin failed")
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, err.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "signin success", map[string]interface{}{
		"token":      fmt.Sprintf("Bearer %s", tokenStr),
		"user_id":    exists.UserId,
		"profile_id": exists.ProfileId,
		"role":       exists.Role,
	})
}

func (a *authController) GetUserWithId(c *fiber.Ctx) error {
	log.Println(c.GetReqHeaders()["Authorization"])

	payload, err := security.ParseToken(c.GetReqHeaders()["Authorization"])
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, err.Error())
	}

	log.Println("payload Id:", payload.Id)

	// user, err := a.userRepo.GetById(payload.Id)
	// if err != nil {
	// 	log.Println(err)
	// 	return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, err.Error())
	// }

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		// "user": user,
	})
}
