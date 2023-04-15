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
)

type ConversationController interface {
	CreateConversation(c *fiber.Ctx) error
	GetByUserId(c *fiber.Ctx) error
}

type conversationController struct {
	conversationRepo repository.ConversationRepository
	profileRepo      repository.ProfileRepository
}

func NewConversationController(conversationRepo repository.ConversationRepository, profileRepo repository.ProfileRepository) ConversationController {
	return &conversationController{conversationRepo: conversationRepo, profileRepo: profileRepo}
}

func (co *conversationController) CreateConversation(c *fiber.Ctx) error {

	req := models.ConversationRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		log.Println(err)
		value, ok := err.(*fiber.Error)
		if ok {
			return util.ResponseNotSuccess(c, value.Code, value.Message)
		}

		return util.ResponseNotSuccess(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	senderId, err := util.CheckStringData(req.SenderId, "sender_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("sender id:", senderId)

	receiverId, err := util.CheckStringData(req.ReceiverId, "receiver_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("receiver id:", receiverId)

	if ok := primitive.IsValidObjectID(senderId); ok == false {
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrIdIsNotPrimitiveObjectID.Error())
	}

	soID, err := primitive.ObjectIDFromHex(senderId)
	if err != nil {
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}

	err = co.profileRepo.GetProfileByFilterForCheckExists(bson.M{"_id": soID})
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

	if ok := primitive.IsValidObjectID(receiverId); ok == false {
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrIdIsNotPrimitiveObjectID.Error())
	}

	roID, err := primitive.ObjectIDFromHex(receiverId)
	if err != nil {
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}

	err = co.profileRepo.GetProfileByFilterForCheckExists(bson.M{"_id": roID})
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

	conversationNew := &models.Conversation{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Members: []string{
			senderId,
			receiverId,
		},
	}

	re, err := co.conversationRepo.Insert(conversationNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create conversation success", map[string]interface{}{
		"conversation_id": re.InsertedID,
	})
}

func (co *conversationController) GetByUserId(c *fiber.Ctx) error {
	userId, err := util.CheckStringData(c.Query("user_id"), "user_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find by user id:", userId)

	conversations, err := co.conversationRepo.GetConversationAllByFilter(bson.M{"members": bson.M{
		"$in": bson.A{
			userId,
		},
	}})
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

	if len(conversations) == 0 {
		log.Println("conversation not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"conversation_list": conversations,
	})
}
