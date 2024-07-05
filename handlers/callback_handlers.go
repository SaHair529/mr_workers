package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"shdbd/mr_workers/db"
	"strconv"
	"strings"
)

type CallbackHandler struct {
	bot *tgbotapi.BotAPI
	db  *db.Database
}

func NewCallbackHandler(bot *tgbotapi.BotAPI, db *db.Database) *CallbackHandler {
	return &CallbackHandler{
		bot: bot,
		db:  db,
	}
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
	request, err := h.db.GetRequestById(requestId)
	errPrintf("Failed to get request: %v", err)
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка. Пожалуйста, обратитесь к разработчику")
		_, err := h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)
		return
	}

	if request == (db.Request{}) {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Заявка не найдена. Пожалуйста, обратитесь к разработчику")
		_, err := h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)
		return
	}

	if !request.Free {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Заявка уже недоступна, так как принята другим человеком")
		_, err := h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)

		newInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Заявка недоступна", "ignore")))
		editMessage := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, newInlineKeyboard)
		_, err = h.bot.Send(editMessage)
		errPrintf("Failed to edit markup %v", err)
		return
	}

	user, err := h.db.GetUserByTgId(request.TelegramID)
	errPrintf("Failed to get request: %v", err)
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка. Пожалуйста, обратитесь к разработчику")
		_, err := h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)
		return
	}

	if user == (db.User{}) {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Пользователь не найден. Пожалуйста, обратитесь к разработчику")
		_, err := h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)
		return
	}

	msgText := fmt.Sprintf("Заявка принята ✅\nСвяжитесь с работодателем для обсуждения деталей по номеру %s", user.Phone)
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, msgText)
	_, err = h.bot.Send(msg)
	errPrintf("Failed to send message %v", err)

	newInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Заявка принята", "ignore")))
	editMessage := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, newInlineKeyboard)
	_, err = h.bot.Send(editMessage)
	errPrintf("Failed to edit markup %v", err)

	err = h.db.SetUnfreeRequest(request.TelegramID)
	errPrintf("Failed to get request: %v", err)
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка. Пожалуйста, обратитесь к разработчику")
		_, err := h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)
		return
	}
}
