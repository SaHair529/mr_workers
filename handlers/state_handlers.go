package handlers

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "shdbd/mr_workers/db"
    "strings"
)

type StateHandler struct {
	bot *tgbotapi.BotAPI
	db *db.Database
}

func NewStateHandler(bot *tgbotapi.BotAPI, db *db.Database) *StateHandler {
    return &StateHandler{
        bot: bot,
        db: db,
    }
}

func (h *StateHandler) HandleState(state string, message *tgbotapi.Message) {
    stateParts := strings.Split(state, "__")
    
    externalState := stateParts[0]
    internalState := stateParts[1]

    switch externalState {
    case "registration":
        h.handleRegistrationState(internalState, message)
    }
}

func (h *StateHandler) handleRegistrationState(internalState string, message *tgbotapi.Message) {
    switch internalState {
    case "pick_role":
        if message.Text == "Клиент" {
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

            for i, specialtity :=  range specialities {
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
        specialityDbRow, err :=h.db.GetSpecialityByTitle(pickedSpeciality)
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

        cities := []string{"Махачкала"}

        var rows [][]tgbotapi.KeyboardButton
        var row []tgbotapi.KeyboardButton

        for i, city := range cities {
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
    }
}