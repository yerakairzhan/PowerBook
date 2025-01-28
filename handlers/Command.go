package handlers

import (
	db "PowerBook/db/sqlc"
	"PowerBook/utils"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func handleCommand(command string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI) {

	userid := strconv.FormatInt(updates.Message.From.ID, 10)
	username := updates.Message.From.UserName
	chatid := updates.Message.Chat.ID
	name := updates.Message.From.FirstName

	ctx := context.Background()
	switch command {
	case "start":
		_, err := queries.GetUser(ctx, userid)
		if err != nil {
			params := db.CreateUserParams{Userid: userid, Username: username}
			_, err := queries.CreateUser(ctx, params)
			if err != nil {
				log.Println("Error creating user", err)
			} else {
				_, text := utils.GetTranslation(ctx, queries, updates, "start")
				msg := tgbotapi.NewMessage(chatid, text+name)
				bot.Send(msg)
			}
		} else {
			_, text := utils.GetTranslation(ctx, queries, updates, "start_1")
			msg := tgbotapi.NewMessage(chatid, text)
			bot.Send(msg)
		}

	}
}
