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

type MessageController interface {
	CreateMessage(c *fiber.Ctx) error
	GetByConversationId(c *fiber.Ctx) error
}

type messageController struct {
	messageRepo      repository.MessageRepository
	conversationRepo repository.ConversationRepository
}

func NewMessageController(messageRepo repository.MessageRepository, conversationRepo repository.ConversationRepository) MessageController {
	return &messageController{messageRepo: messageRepo, conversationRepo: conversationRepo}
}

func (m *messageController) CreateMessage(c *fiber.Ctx) error {

	req := models.MessageRequest{}
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

	conversationId, err := util.CheckStringData(req.ConversationId, "conversation_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("conversation id:", conversationId)

	if req.Text == "" {
		log.Println(util.ReturnError(util.ErrRequireParameter.Error() + "text"))
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrRequireParameter.Error()+"text")
	}
	log.Println("conversation text:", req.Text)

	con, err := m.conversationRepo.GetConversationById(conversationId)
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

	chcek := true
	for _, v := range con.Members {
		if v == senderId {
			chcek = false
			break
		}
	}

	if chcek {
		log.Println("not permiistion")
		return util.ResponseNotSuccess(c, fiber.StatusUnauthorized, "not permission")
	}

	messageNew := &models.Message{
		Id:             primitive.NewObjectID(),
		CreatedAt:      time.Now().Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
		ConversationId: conversationId,
		Sender:         senderId,
		Text:           req.Text,
	}

	re, err := m.messageRepo.Insert(messageNew)
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusInternalServerError, util.ErrInternalServerError.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusCreated, "create message success", map[string]interface{}{
		"message_id": re.InsertedID,
	})
}

func (m *messageController) GetByConversationId(c *fiber.Ctx) error {
	conversationId, err := util.CheckStringData(c.Query("conversation_id"), "conversation_id")
	if err != nil {
		log.Println(err)
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, err.Error())
	}
	log.Println("find by conversation id:", conversationId)

	messages, err := m.messageRepo.GetConversationAllByFilter(bson.M{"conversation_id": conversationId})
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

	if len(messages) == 0 {
		log.Println("message not found")
		return util.ResponseNotSuccess(c, fiber.StatusBadRequest, util.ErrNotFound.Error())
	}

	return util.ResponseSuccess(c, fiber.StatusOK, "success", map[string]interface{}{
		"message_list": messages,
	})
}
