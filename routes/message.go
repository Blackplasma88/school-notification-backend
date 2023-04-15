package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type messageRoutes struct {
	messageController controller.MessageController
}

func NewMessageRoute(messageController controller.MessageController) Routes {
	return &messageRoutes{messageController: messageController}
}

func (r *messageRoutes) Install(app *fiber.App) {
	app.Get("/message/conversation-id", r.messageController.GetByConversationId)

	app.Post("/message/create", r.messageController.CreateMessage)
}
