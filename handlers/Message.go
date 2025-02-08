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
		log.Println("someone read ", minutes)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Введите число минут цифрами."))
			return
		}

		params := db.CreateReadingLogParams{
			Userid:      strconv.FormatInt(userid, 10),
			Date:        time.Now(),
			MinutesRead: int32(minutes),
		}
		_, err = queries.CreateReadingLog(ctx, params)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					params := db.UpdateReadingLogParams{
						Userid:      strconv.FormatInt(userid, 10),
						Date:        time.Now(),
						MinutesRead: int32(minutes),
					}
					queries.UpdateReadingLog(ctx, params)
				}

				_, text := utils.GetTranslation(ctx, queries, updates, "read_3")
				msg := tgbotapi.NewMessage(chatID, text)
				msg.ParseMode = "HTML"
				_, err = bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			_, text := utils.GetTranslation(ctx, queries, updates, "read_1")
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ParseMode = "HTML"
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}

		err = queries.DeleteUserState(ctx, strconv.FormatInt(userid, 10))
		if err != nil {
			log.Println(err)
		}

		utils.LoadConfig()
		err = api.AddReadingMinutes(utils.GoogleApi, strconv.FormatInt(userid, 10), minutes)
		if err != nil {
			log.Fatalf("Error adding reading minutes: %v", err)
		}

	} else {
		_, text := utils.GetTranslation(ctx, queries, updates, "read_2")
		msg := tgbotapi.NewMessage(chatID, text)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
