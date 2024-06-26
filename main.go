package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
    "os"
    "os/signal"
    "shdbd/mr_workers/config"
    "shdbd/mr_workers/handlers"
    "syscall"
)

func main() {
	cfg, err := config.LoadConfig()
	onFail("Failed to load config %v", err)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	onFail("Failed to create bot %v", err)

	messageHandler := handlers.NewMessageHandler(bot)
	commandHanler := handlers.NewCommandHandler(bot)
	callbackHandler := handlers.NewCallbackHandler(bot)

	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = 60

	updatesChan, err := bot.GetUpdatesChan(updates)
	onFail("Failed to create updates channel %v", err)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case update := <- updatesChan:
			if update.Message != nil {
				if update.Message.IsCommand() {
					commandHanler.HandleCommand(update.Message)
				} else {
					messageHandler.HandleMessage(update.Message)
				}
			} else if update.CallbackQuery != nil {
				callbackHandler.HandleCallback(update)
			}
		}
	}
}

func onFail(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}