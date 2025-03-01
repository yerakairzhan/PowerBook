package handlers

import (
	"PowerBook/db/sqlc"
	"PowerBook/utils"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

var scheduledMessages = make(map[int64]*time.Timer) // Tracks scheduled notifications

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

			err, text := utils.GetTranslation(ctx, db, update, "timer_2")
			if err != nil {
				log.Println(err)
			}
			scheduleReminder(bot, chatID, userID, text)

		} else {
			_, text := utils.GetTranslation(ctx, db, update, "register_1")
			SendMessage(bot, chatID, text)
		}
	}
}

func scheduleReminder(bot *tgbotapi.BotAPI, chatID int64, userID int64, text string) {
	if timer, exists := scheduledMessages[userID]; exists {
		timer.Stop()
		delete(scheduledMessages, userID)
	}

	now := time.Now()
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 30, 0, now.Location())

	if now.After(targetTime) {
		targetTime = targetTime.Add(24 * time.Hour)
	}

	durationUntilReminder := time.Until(targetTime)

	// Schedule a message
	timer := time.AfterFunc(durationUntilReminder, func() {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Ошибка отправки напоминания:", err)
		}
		delete(scheduledMessages, userID)
	})

	scheduledMessages[userID] = timer
}
