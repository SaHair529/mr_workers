package handlers

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"shdbd/mr_workers/constants"
	"shdbd/mr_workers/db"
)

type CommandHandler struct {
	bot      *tgbotapi.BotAPI
	messages CommandHandlerMessages
	db       *db.Database
}

type CommandHandlerMessages struct {
	Registration                 string `json:"registration"`
	ContactReceived              string `json:"contact_received"`
	ContactReceivedFail          string `json:"contact_received_fail"`
	ContactReceivedAlreadyExists string `json:"contact_received_already_exists"`
	Default                      string `json:"default_unregistered"`
}

func NewCommandHandler(bot *tgbotapi.BotAPI, db *db.Database) *CommandHandler {
	messagesJson, err := ioutil.ReadFile("messages.json")
	errPrintf("Failed to read file %v", err)
	var messages CommandHandlerMessages
	err = json.Unmarshal(messagesJson, &messages)
	errPrintf("Failed to unmarshal json %v", err)

	return &CommandHandler{
		bot:      bot,
		messages: messages,
		db:       db,
	}
}

func (h *CommandHandler) HandleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "registration":
		h.handleRegistrationCommand(message)
	case "reset":
		h.handleResetCommand(message)
	case "main":
		h.handleMainCommand(message)
	case "create_request":
		h.handleCreateRequestCommand(message)
	default:
		h.handleDefault(message)
	}
}

func (h *CommandHandler) HandleContact(message *tgbotapi.Message) {
	user, err := h.db.GetUserByTgId(message.Chat.ID)
	errPrintf("Failed to get user %v", err)

	if user == (db.User{}) {
		msg := tgbotapi.NewMessage(message.Chat.ID, h.messages.ContactReceivedAlreadyExists)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		_, err = h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)

		return
	}

	fullname := message.Contact.FirstName + " " + message.Contact.LastName
	err = h.db.AddUser(message.Chat.ID, fullname, message.Contact.PhoneNumber)
	if err != nil {
		log.Printf("Failed to add user %v", err)

		msg := tgbotapi.NewMessage(message.Chat.ID, h.messages.ContactReceivedFail)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		_, err = h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)

		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, h.messages.ContactReceived)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	_, err = h.bot.Send(msg)
	errPrintf("Failed to send message %v", err)
}

func (h *CommandHandler) handleDefault(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, constants.MainMenuMessage)

	_, err := h.bot.Send(msg)
	errPrintf("Failed to send message %v", err)
}

func (h *CommandHandler) handleRegistrationCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы клиент или рабочий?")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Клиент"),
			tgbotapi.NewKeyboardButton("Рабочий"),
		),
	)

	_, err := h.bot.Send(msg)
	errPrintf("Failed to send message %v", err)

	err = h.db.SetUserState(message.Chat.ID, "registration__pick_role")
	errPrintf("Failed to set user state %v", err)
}

func (h *CommandHandler) handleResetCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Действие отменено")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	_, err := h.bot.Send(msg)
	errPrintf("Failed to send message %v", err)

	err = h.db.SetUserState(message.Chat.ID, "")
	errPrintf("Failed to send message %v", err)
}

func (h *CommandHandler) handleCreateRequestCommand(message *tgbotapi.Message) {
	specialities, err := h.db.GetAllSpecialities()
	if err != nil {
		errPrintf("Failed to get specialities %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Повторите попытку позже")

		_, err := h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)
		return
	}

	var rows [][]tgbotapi.KeyboardButton
	var row []tgbotapi.KeyboardButton

	for i, speciality := range specialities {
		button := tgbotapi.NewKeyboardButton(speciality.Speciality)
		row = append(row, button)

		// Добавить строку кнопок каждые 2 кнопки
		if (i+1)%2 == 0 {
			rows = append(rows, row)
			row = []tgbotapi.KeyboardButton{}
		}
	}

	// Добавить оставшиеся кнопки
	if len(row) > 0 {
		rows = append(rows, row)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Какой специалист Вам нужен? (Выберите нажав на кнопку ниже)")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows...)

	_, err = h.bot.Send(msg)
	errPrintf("Failed to send message %v", err)

	err = h.db.SetUserState(message.Chat.ID, "createrequest__pick_specialist")
	errPrintf("Failed to set user state %v", err)
}

func (h *CommandHandler) handleMainCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, constants.MainMenuMessage)

	_, err := h.bot.Send(msg)
	errPrintf("Failed to send message %v", err)
}
