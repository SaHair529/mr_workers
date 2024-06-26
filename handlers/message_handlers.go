package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type MessageHandler struct {
	bot *tgbotapi.BotAPI
}

func NewMessageHandler(bot *tgbotapi.BotAPI) *MessageHandler {
	return &MessageHandler{
		bot: bot,
	}
}

func (h *MessageHandler) HandleMessage(command *tgbotapi.Message) {
	// todo handle message
}
