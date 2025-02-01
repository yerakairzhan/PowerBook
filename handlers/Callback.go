package handlers

import (
	db "PowerBook/db/sqlc"
	"PowerBook/utils"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func handleCallback(command string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64) {
	userid := strconv.FormatInt(updates.CallbackQuery.From.ID, 10)

	ctx := context.Background()
	switch command {
	case "callback_register":
		userid := updates.CallbackQuery.From.ID
		callbackRegister(queries, updates, bot, chatid, strconv.FormatInt(userid, 10))
	case "callback_read":
		callbackRead(queries, updates, bot, userid, chatid)
	case "callback_stat":
	case "callback_top":
		callbackTop(bot, chatid, queries, updates)
	case "callback_lang":
		callbackLang(queries, updates, bot, chatid)
	case "callback_ru":
		err := changeLang(queries, userid, "ru")
		if err != nil {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_2")
			SendMessage(bot, chatid, text)
		} else {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_1")
			SendMessage(bot, chatid, text)
		}
	case "callback_en":
		err := changeLang(queries, userid, "en")
		if err != nil {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_2")
			SendMessage(bot, chatid, text)
		} else {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_1")
			SendMessage(bot, chatid, text)
		}
	case "callback_kz":
		err := changeLang(queries, userid, "kz")
		if err != nil {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_2")
			SendMessage(bot, chatid, text)
		} else {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_1")
			SendMessage(bot, chatid, text)
		}

	case "timer_15_00":
		changeTimer(queries, userid, bot, updates, chatid, 15)
	case "timer_18_00":
		changeTimer(queries, userid, bot, updates, chatid, 18)
	case "timer_21_00":
		changeTimer(queries, userid, bot, updates, chatid, 21)
	default:
		log.Println("Unknown command: " + command)
	}
}
