package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type CallbackHandler struct {
	bot *tgbotapi.BotAPI
}

func NewCallbackHandler(bot *tgbotapi.BotAPI) *CallbackHandler {
	return &CallbackHandler{bot: bot}
}

func (cbh *CallbackHandler) HandleCallback(callback tgbotapi.Update) {}
