package handlers

import (
	"PowerBook/api"
	db "PowerBook/db/sqlc"
	"PowerBook/utils"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lib/pq"
	"log"
	"strconv"
	"time"
)

func handleMessage(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI) {
	ctx := context.Background()
	chatID := updates.Message.Chat.ID
	userid := updates.Message.From.ID
	userMessage := updates.Message.Text

	state, _ := queries.GetUserState(ctx, strconv.FormatInt(userid, 10))

	if state.String == "waiting_for_reading_time" {
		minutes, err := strconv.Atoi(userMessage)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Введите число минут цифрами."))
			return
		}

		log.Printf("User %d read %d minutes", userid, minutes)
		userIDStr := strconv.FormatInt(userid, 10)
		now := time.Now()

		params := db.CreateReadingLogParams{
			Userid:      userIDStr,
			Date:        now,
			MinutesRead: int32(minutes),
		}

		if _, err = queries.CreateReadingLog(ctx, params); err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				updateParams := db.UpdateReadingLogParams{
					Userid:      userIDStr,
					Date:        now,
					MinutesRead: int32(minutes),
				}
				queries.UpdateReadingLog(ctx, updateParams)
			}
		}

		translationKey := "read_1"
		if err != nil {
			translationKey = "read_3"
		}

		_, text := utils.GetTranslation(ctx, queries, updates, translationKey)
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		if _, err = bot.Send(msg); err != nil {
			log.Println("Error sending message:", err)
		}

		if err = queries.DeleteUserState(ctx, userIDStr); err != nil {
			log.Println("Error deleting user state:", err)
		}

		utils.LoadConfig()
		if err = api.AddReadingMinutes(utils.GoogleApi, userIDStr, minutes); err != nil {
			log.Fatalf("Error adding reading minutes: %v", err)
		}
	} else {
		_, text := utils.GetTranslation(ctx, queries, updates, "read_2")
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		if _, err := bot.Send(msg); err != nil {
			log.Println("Error sending fallback message:", err)
		}
	}
}
