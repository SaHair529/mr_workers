package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type MessageHandler struct {
	bot *tgbotapi.BotAPI
}

func NewMessageHandler(bot *tgbotapi.BotAPI) *MessageHandler {
	return &MessageHandler{
		bot: bot,
	}
}

func (h *MessageHandler) HandleMessage(message *tgbotapi.Message) {
	responseMsg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	_, err := h.bot.Send(responseMsg)
	onFail("Failed to send message %v", err)

}

func onFail(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
