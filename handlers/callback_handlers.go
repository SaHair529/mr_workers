package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

type CallbackHandler struct {
	bot *tgbotapi.BotAPI
}

func NewCallbackHandler(bot *tgbotapi.BotAPI) *CallbackHandler {
	return &CallbackHandler{bot: bot}
}

func (h *CallbackHandler) HandleCallback(callback *tgbotapi.CallbackQuery) {
	callbackParts := strings.Split(callback.Data, "__")

	externalCallback := callbackParts[0]
	internalCallback := callbackParts[1]

	switch externalCallback {
	case "accept_request":
		requestId, err := strconv.ParseInt(internalCallback, 10, 64)
		if err != nil {
			errPrintf("Failed to convert string to int64: %v", err)
			msg := tgbotapi.NewMessage(requestId, "Произошла ошибка. Пожалуйста, обратитесь к разработчику")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message: %v", err)
			return
		}

		h.handleAcceptRequest(requestId, callback)
	}
}

func (h *CallbackHandler) handleAcceptRequest(requestId int64, callback *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, strconv.FormatInt(requestId, 10))
	_, err := h.bot.Send(msg)
	errPrintf("Failed to send message: %v", err)

	// todo обработать принятие заявки
}
