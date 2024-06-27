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
    
}