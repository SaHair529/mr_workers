package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
    "shdbd/mr_workers/constants"
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
	msg := tgbotapi.NewMessage(message.Chat.ID, constants.MainMenuMessage)

	_, err := h.bot.Send(msg)
	errPrintf("Failed to send message %v", err)

}

func errPrintf(message string, err error) {
	if err != nil {
		log.Printf(message, err)
	}
}
