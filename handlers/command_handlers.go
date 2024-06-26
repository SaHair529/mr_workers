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
		responseMsg := tgbotapi.NewMessage(message.Chat.ID, h.messages.Registration)
		_, err := h.bot.Send(responseMsg)
		onFail("Failed to send message %v", err)
    default:

    }
}
