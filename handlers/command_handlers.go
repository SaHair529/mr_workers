package handlers

import (
    "encoding/json"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "io/ioutil"
)

type CommandHandler struct {
	bot *tgbotapi.BotAPI
	messages CommandHandlerMessages
}

type CommandHandlerMessages struct {
	Registration string `json:"registration"`
	ContactReceived string `json:"contact_received"`
	Default string `json:"default_unregistered"`
}

func NewCommandHandler(bot *tgbotapi.BotAPI) *CommandHandler {
	messagesJson, err := ioutil.ReadFile("messages.json")
	onFail("Failed to read file %v", err)
	var messages CommandHandlerMessages
	err = json.Unmarshal(messagesJson, &messages)
	onFail("Failed to unmarshal json %v", err)

	return &CommandHandler{
		bot: bot,
		messages: messages,
	}
}

func (h *CommandHandler) HandleCommand(message *tgbotapi.Message) {
    switch message.Command() {
    case "registration":
		h.handleRegistrationCommand(message)
    default:
		h.handleDefault(message)
    }
}

func (h *CommandHandler) HandleContact(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, h.messages.ContactReceived)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	_, err := h.bot.Send(msg)
	onFail("Failed to send message %v", err)
}

func (h *CommandHandler) handleDefault(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, h.messages.Default)

	_, err := h.bot.Send(msg)
	onFail("Failed to send message %v", err)
}

func (h *CommandHandler) handleRegistrationCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, h.messages.Registration)
	button := tgbotapi.NewKeyboardButtonContact("Поделиться контактными данными")
	keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(button))
	msg.ReplyMarkup = keyboard

	_, err := h.bot.Send(msg)

	onFail("Failed to send message %v", err)
}
