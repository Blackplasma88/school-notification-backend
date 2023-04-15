package routes

import (
	"school-notification-backend/controller"

	"github.com/gofiber/fiber/v2"
)

type conversationRoutes struct {
	conversationController controller.ConversationController
}

func NewConversationRoute(conversationController controller.ConversationController) Routes {
	return &conversationRoutes{conversationController: conversationController}
}

func (r *conversationRoutes) Install(app *fiber.App) {
	app.Get("/conversation/user-id", r.conversationController.GetByUserId)

	app.Post("/conversation/create", r.conversationController.CreateConversation)
}
