package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	//    "os"
	//    "os/signal"
	"shdbd/mr_workers/config"
	db2 "shdbd/mr_workers/db"
	"shdbd/mr_workers/handlers"
	//    "syscall"
)

func main() {
	cfg, err := config.LoadConfig()
	errPrintf("Failed to load config %v", err)

	db, err := db2.ConnectDB(cfg.DatabaseURL)
	errPrintf("Failed to connect db %v", err)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	errPrintf("Failed to create bot %v", err)

	stateHandler := handlers.NewStateHandler(bot, db)
	messageHandler := handlers.NewMessageHandler(bot)
	commandHanler := handlers.NewCommandHandler(bot, db)
	callbackHandler := handlers.NewCallbackHandler(bot, db)

	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = 60

	updatesChan, err := bot.GetUpdatesChan(updates)
	errPrintf("Failed to create updates channel %v", err)

	//	sigCh := make(chan os.Signal, 1)
	//	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case update := <-updatesChan:
			var chatId int64
			if update.Message != nil {
				chatId = update.Message.Chat.ID
			} else {
				chatId = update.CallbackQuery.Message.Chat.ID
			}

			userState, err := db.GetUserState(chatId)
			errPrintf("Failed to get user state %v", err)

			if userState != "" && update.Message != nil && !update.Message.IsCommand() {
				stateHandler.HandleState(userState, update.Message)
			} else if update.Message != nil && update.Message.Contact != nil {
				commandHanler.HandleContact(update.Message)
			} else if update.Message != nil {
				if update.Message != nil && update.Message.IsCommand() {
					commandHanler.HandleCommand(update.Message)
				} else {
					messageHandler.HandleMessage(update.Message)
				}
			} else if update.CallbackQuery != nil {
				callbackHandler.HandleCallback(update.CallbackQuery)
			}
		}
	}
}

func errPrintf(message string, err error) {
	if err != nil {
		log.Printf(message, err)
	}
}
