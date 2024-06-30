package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"shdbd/mr_workers/constants"
	"shdbd/mr_workers/db"
	"strings"
)

type StateHandler struct {
	bot *tgbotapi.BotAPI
	db  *db.Database
}

func NewStateHandler(bot *tgbotapi.BotAPI, db *db.Database) *StateHandler {
	return &StateHandler{
		bot: bot,
		db:  db,
	}
}

func (h *StateHandler) HandleState(state string, message *tgbotapi.Message) {
	stateParts := strings.Split(state, "__")

	externalState := stateParts[0]
	internalState := stateParts[1]

	switch externalState {
	case "registration":
		h.handleRegistrationState(internalState, message)
	case "createrequest":
		h.handleCreateRequestState(internalState, message)
	}
}

func (h *StateHandler) handleRegistrationState(internalState string, message *tgbotapi.Message) {
	switch internalState {
	case "pick_role":
		if message.Text == "Клиент" {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Поделитесь вашим номером телефона, нажав на кнопку снизу")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButtonContact("Поделиться контактом"),
				),
			)
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)

			h.db.SetUserState(message.Chat.ID, "registration__client_share_contact")
		} else if message.Text == "Рабочий" {
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

			for i, specialtity := range specialities {
				button := tgbotapi.NewKeyboardButton(specialtity.Speciality)
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

			keyboard := tgbotapi.NewReplyKeyboard(rows...)

			msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите вашу специальность")
			msg.ReplyMarkup = keyboard

			_, err = h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)

			h.db.SetUserState(message.Chat.ID, "registration__worker_pick_speciality")
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите роль, нажав на подходящую кнопку ниже")

			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
		}

	case "worker_pick_speciality":
		pickedSpeciality := message.Text
		specialityDbRow, err := h.db.GetSpecialityByTitle(pickedSpeciality)
		if err != nil {
			errPrintf("Failed to get speciality %v", err)

			msg := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка, повторите позже")
			_, err = h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}
		if specialityDbRow == (db.Speciality{}) {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Введенная Вами специальность некорректна. Выберите ее из предоставленных ниже вариантов")
			_, err = h.bot.Send(msg)
			return
		}

		h.db.SetWorkerSpeciality(message.Chat.ID, pickedSpeciality)

		var rows [][]tgbotapi.KeyboardButton
		var row []tgbotapi.KeyboardButton

		for i, city := range constants.Cities {
			button := tgbotapi.NewKeyboardButton(city)
			row = append(row, button)

			if (i+1)%2 == 0 {
				rows = append(rows, row)
				row = []tgbotapi.KeyboardButton{}
			}
		}

		if len(row) > 0 {
			rows = append(rows, row)
		}

		keyboard := tgbotapi.NewReplyKeyboard(rows...)

		msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите город")
		msg.ReplyMarkup = keyboard

		_, err = h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)

		h.db.SetUserState(message.Chat.ID, "registration__worker_pick_city")

	case "worker_pick_city":
		pickedCity := message.Text

		if !cityExists(constants.Cities, pickedCity) {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Введенный Вами город некорректен. Выберите подходящий город, нажав на одну из кнопок ниже")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}

		err := h.db.SetRowField(message.Chat.ID, "workers", "city", pickedCity)
		errPrintf("Failed to set worker city %v", err)

		msg := tgbotapi.NewMessage(message.Chat.ID, "Поделитесь вашим номером телефона, нажав на кнопку снизу")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonContact("Поделиться контактом"),
			),
		)
		_, err = h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)

		h.db.SetUserState(message.Chat.ID, "registration__worker_share_contact")

	case "worker_share_contact":
		if message.Contact == nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Некорректный ответ. Предоставьте свои контактные данные, нажав на кнопку снизу")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}

		fullname := message.Contact.FirstName + " " + message.Contact.LastName
		err := h.db.SetWorkerContactData(message.Chat.ID, fullname, message.Contact.PhoneNumber)
		errPrintf("Failed to set worker contact data %v", err)

		msg := tgbotapi.NewMessage(message.Chat.ID, "Регистрация прошла успешно ✅")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		_, err = h.bot.Send(msg)

		h.db.SetUserState(message.Chat.ID, "")

	case "client_share_contact":
		if message.Contact == nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Некорректный ответ. Предоставьте свои контактные данные, нажав на кнопку снизу")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}

		fullname := message.Contact.FirstName + " " + message.Contact.LastName
		err := h.db.AddUser(message.Chat.ID, fullname, message.Contact.PhoneNumber)
		errPrintf("Failed to add user %v", err)

		msg := tgbotapi.NewMessage(message.Chat.ID, "Регистрация прошла успешно ✅")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		_, err = h.bot.Send(msg)

		h.db.SetUserState(message.Chat.ID, "")
	}
}

func (h *StateHandler) handleCreateRequestState(internalState string, message *tgbotapi.Message) {
	switch internalState {
	case "pick_specialist":
		pickedSpecialist := message.Text
		specialities, err := h.db.GetAllSpecialities()
		errPrintf("Failed to get specialities %v", err)

		if !isSpecialityInList(specialities, pickedSpecialist) {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Введенный Вами специалист несуществует. Выберите подходящего, нажав на одну из кнопок ниже")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}

		//todo err = h.db.SetRowField(message.Chat.ID, "requests", "specialist", pickedSpecialist)

		var rows [][]tgbotapi.KeyboardButton
		var row []tgbotapi.KeyboardButton

		for i, city := range constants.Cities {
			button := tgbotapi.NewKeyboardButton(city)
			row = append(row, button)

			if (i+1)%2 == 0 {
				rows = append(rows, row)
				row = []tgbotapi.KeyboardButton{}
			}
		}

		if len(row) > 0 {
			rows = append(rows, row)
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите город")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows...)

		_, err = h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)

		err = h.db.SetUserState(message.Chat.ID, "createrequest__pick_city")
		errPrintf("Failed to set user state %v", err)

	case "pick_city":
		pickedCity := message.Text

		if !cityExists(constants.Cities, pickedCity) {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Введенный Вами город некорректен. Выберите подходящий город, нажав на одну из кнопок ниже")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}

		//todo err := h.db.SetRowField(message.Chat.ID, "requests", "city", pickedCity)

		msg := tgbotapi.NewMessage(message.Chat.ID, "Опишите работу, которую необходимо проделать")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		_, err := h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)

		err = h.db.SetUserState(message.Chat.ID, "createrequest__describe_work")
		errPrintf("Failed to set user state %v", err)

	case "describe_work":
		description := message.Text
		if description == "" {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Введите сообщение с текстом описания работ, которые нужно провести")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}

		//todo err := h.db.SetRowField(message.Chat.ID, "requests", "description", description)

		err := h.db.SetUserState(message.Chat.ID, "")
		errPrintf("Failed to set user state %v", err)
	}
}

func isSpecialityInList(specialities []db.Speciality, pickedSpeciality string) bool {
	for _, speciality := range specialities {
		if speciality.Speciality == pickedSpeciality {
			return true
		}
	}
	return false
}

func cityExists(cities []string, city string) bool {
	for _, c := range cities {
		if c == city {
			return true
		}
	}
	return false
}
