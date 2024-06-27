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
	errPrintf("Failed to send message %v", err)

}

func errPrintf(message string, err error) {
	if err != nil {
		log.Printf(message, err)
	}
}
