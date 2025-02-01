package handlers

import (
	"PowerBook/db/sqlc"
	"PowerBook/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"time"
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

		timer, err := db.GetTimer(ctx, strconv.FormatInt(userID, 10))
		if err != nil {
			log.Println("Error fetching timer:", err)
		} else {
			go ScheduleDaily(timer.Hour(), bot, chatID, db, update)
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

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) int {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	Sent, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
		return 0
	} else {
		return Sent.MessageID
	}
	return 0
}

func checkRegistration(ctx context.Context, db *db.Queries, userID int64) (bool, error) {
	reged, err := db.GetRegistered(ctx, strconv.FormatInt(userID, 10))
	if err != nil {
		log.Println("Error checking registration:", err)
		if err.Error() == "sql: no rows in result set" {
			return true, nil
		}
		return false, err
	}
	return reged.Bool, nil
}

func ScheduleDaily(hour int, bot *tgbotapi.BotAPI, chatid int64, queries *db.Queries, update tgbotapi.Update) {
	for {
		ctx := context.Background()

		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		fmt.Println("Следующая отправка:", next)
		time.Sleep(time.Until(next))

		_, text := utils.GetTranslation(ctx, queries, update, "timer_2")
		SendMessage(bot, chatid, text)
		fmt.Println("Сообщение отправлено:", time.Now().Format("15:04:05"))
	}
}
