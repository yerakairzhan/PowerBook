package handlers

import (
	db "PowerBook/db/sqlc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func SetupHandlers(bot *tgbotapi.BotAPI, db *db.Queries) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			// handleCallback(bot, update, db) // Uncomment and implement when needed
		} else if update.Message != nil {
			if update.Message.IsCommand() {
				command := update.Message.Command()
				handleCommand(command, db, update, bot) // Now passing the correct arguments (command string and db object)
			} else if update.Message.Photo != nil {
				log.Println("Photo sent")
				// handlePhoto(bot, update, db) // Uncomment and implement when needed
			}
		}
	}
}
