package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type CommandHandler struct {
	bot *tgbotapi.BotAPI
}

func NewCommandHandler(bot *tgbotapi.BotAPI) *CommandHandler {
	return &CommandHandler{
		bot: bot,
	}
}

func (h *CommandHandler) HandleCommand(message *tgbotapi.Message) {
	responseMsg := tgbotapi.NewMessage(message.Chat.ID, "Command "+message.Command())

	_, err := h.bot.Send(responseMsg)
	onFail("Failed to send message %v", err)
}
