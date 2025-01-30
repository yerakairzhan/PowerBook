package handlers

import (
	db "PowerBook/db/sqlc"
	"PowerBook/utils"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func SetupHandlers(bot *tgbotapi.BotAPI, db *db.Queries) {
	ctx := context.Background()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var chatID int64
		var userID int64
		var command string

		if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
			userID = update.CallbackQuery.From.ID
			command = update.CallbackQuery.Data
		} else if update.Message != nil {
			chatID = update.Message.Chat.ID
			userID = update.Message.From.ID
			if update.Message.IsCommand() {
				command = update.Message.Command()
			}
		} else {
			continue
		}

		registered, err := checkRegistration(ctx, db, userID)
		if err != nil {
			continue
		}

		if registered || command == "start" || command == "callback_register" {
			if update.CallbackQuery != nil {
				handleCallback(command, db, update, bot, chatID)
			} else if update.Message != nil {
				if command != "" {
					handleCommand(command, db, update, bot, chatID)
				} else {
					handleMessage(db, update, bot)
				}
			}
		} else {
			_, text := utils.GetTranslation(ctx, db, update, "register_1")
			SendMessage(bot, chatID, text)
		}
	}
}

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}
}

func checkRegistration(ctx context.Context, db *db.Queries, userID int64) (bool, error) {
	reged, err := db.GetRegistered(ctx, strconv.FormatInt(userID, 10))
	if err != nil {
		log.Println("Error checking registration:", err)
		return false, err
	}
	return reged.Bool, nil
}
