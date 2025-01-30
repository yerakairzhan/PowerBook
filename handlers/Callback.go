package handlers

import (
	db "PowerBook/db/sqlc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func handleCallback(command string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64) {
	switch command {
	case "callback_register":
		userid := updates.CallbackQuery.From.ID
		callback_register(queries, updates, bot, chatid, strconv.FormatInt(userid, 10))
	case "callback_read":
		callback_read(queries, updates, bot)
	case "callback_stat":
	case "callback_top":
	case "callback_lang":
		callback_lang(queries, updates, bot, chatid)
	case "callback_ru":
	case "callback_en":
	case "callback_kz":
	default:
		log.Println("Unknown command: " + command)
	}
}
