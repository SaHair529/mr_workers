package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"shdbd/mr_workers/constants"
	"shdbd/mr_workers/db"
	"strconv"
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
		freeRequest, err := h.db.GetFreeRequest(message.Chat.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка. Пожалуйста, обратитесь к разработчику")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			err = h.db.SetUserState(message.Chat.ID, "")
			errPrintf("Failed to set user state %v", err)
			return
		}

		if freeRequest != (db.Request{}) {
			msg := tgbotapi.NewMessage(message.Chat.ID, "У вас уже есть одна свободная заявка, которая еще никем не принята. Пожалуйста закройте её или подождите пока её примут")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			err = h.db.SetUserState(message.Chat.ID, "")
			errPrintf("Failed to set user state %v", err)
			return
		}

		pickedSpecialist := message.Text
		specialities, err := h.db.GetAllSpecialities()
		errPrintf("Failed to get specialities %v", err)

		if !isSpecialityInList(specialities, pickedSpecialist) {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Введенный Вами специалист несуществует. Выберите подходящего, нажав на одну из кнопок ниже")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}

		freeRequest, err = h.db.GetFreeRequest(message.Chat.ID)
		errPrintf("Failed to get free request %v", err)
		if freeRequest == (db.Request{}) {
			err = h.db.CreateFreeRequest(message.Chat.ID)
			errPrintf("Failed to create free request %v", err)
		}
		err = h.db.SetFreeRequestField(message.Chat.ID, "specialist", pickedSpecialist)
		errPrintf("Failed to set free request field %v", err)

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

		err := h.db.SetFreeRequestField(message.Chat.ID, "city", pickedCity)
		errPrintf("Failed to set row field %v", err)

		msg := tgbotapi.NewMessage(message.Chat.ID, "Опишите работу, которую необходимо проделать")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		_, err = h.bot.Send(msg)
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

		err := h.db.SetFreeRequestField(message.Chat.ID, "description", description)
		errPrintf("Failed to set row field %v", err)

		err = h.db.SetUserState(message.Chat.ID, "createrequest__submit_request")
		errPrintf("Failed to set user state %v", err)

		request, err := h.db.GetFreeRequest(message.Chat.ID)
		if err != nil {
			errPrintf("Failed to get free request %v", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка. Повторите попытку позже")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}
		if request == (db.Request{}) {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка. Повторите попытку позже")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}

		msgText := fmt.Sprintf(`Вот как выглядет Ваша заявка:

Ваш ID в телеграме: %s
Необходимый специалист: %s
Город: %s
Описание работ: %s

Подтверждаете создание заявки? (Нажмите кнопку ниже)
`, request.TelegramID, request.Specialist, request.City, request.Description)

		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Подтвердить"),
				tgbotapi.NewKeyboardButton("Отменить"),
			),
		)
		_, err = h.bot.Send(msg)
		errPrintf("Failed to send message %v", err)

	case "submit_request":
		if message.Text != "Подтвердить" && message.Text != "Отменить" {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Некорректный ответ. Нажмите на подходящую кнопку снизу. \"Подтвердить\", если хотите создать заявку, \"Отменить\", если хотите отменить создание заявки")
			_, err := h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
			return
		}

		if message.Text == "Подтвердить" {
			request, err := h.db.GetFreeRequest(message.Chat.ID)
			if err != nil {
				errPrintf("Failed to get free request %v", err)

				msg := tgbotapi.NewMessage(message.Chat.ID, "Не удалось найти заявку в базе данных. Пожалуйста, сообщите об ошибке разработчику")
				_, err := h.bot.Send(msg)
				errPrintf("Failed to send message %v", err)
				return
			}

			workers, err := h.db.GetFreeWorkersByCityAndSpeciality(request.City, request.Specialist)
			if err != nil {
				errPrintf("Failed to get workers %v", err)
				msg := tgbotapi.NewMessage(message.Chat.ID, "Возникла ошибка при рассылке заявки рабочим. Пожалуйста, сообщите об этом разработчику")

				_, err := h.bot.Send(msg)
				errPrintf("Failed to send message %v", err)
				return
			}

			for _, worker := range workers {
				messageText := fmt.Sprintf("⚡️⚡️⚡️ Новая заявка!\nОписание работ: %s", request.Description)
				msg := tgbotapi.NewMessage(worker.TelegramID, messageText)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Принять✅", "accept_"+strconv.FormatInt(request.ID, 10)),
						tgbotapi.NewInlineKeyboardButtonData("Отклонить❌", "decline_"+strconv.FormatInt(request.ID, 10)),
					),
				)
				_, err = h.bot.Send(msg)
				errPrintf("Failed to send message %v", err)
			}
			errPrintf("Failed to set unfree request %v", err)

			msg := tgbotapi.NewMessage(message.Chat.ID, "✅ Заявка создана! Ждите откликов рабочих")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err = h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
		} else if message.Text == "Отменить" {
			err := h.db.DeleteFreerequest(message.Chat.ID)
			errPrintf("Failed to delete free request %v", err)

			msg := tgbotapi.NewMessage(message.Chat.ID, "Заявка успешно удалена")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err = h.bot.Send(msg)
			errPrintf("Failed to send message %v", err)
		}

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
