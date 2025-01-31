package handlers

import (
	db "PowerBook/db/sqlc"
	"PowerBook/utils"
	"context"
	"database/sql"
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
		callback_register(queries, updates, bot, chatid, strconv.FormatInt(userid, 10))
	case "callback_read":
		callback_read(queries, updates, bot, userid, chatid)
	case "callback_stat":
	case "callback_top":
	case "callback_lang":
		callback_lang(queries, updates, bot, chatid)
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
	default:
		log.Println("Unknown command: " + command)
	}
}

func changeLang(queries *db.Queries, userid string, lang string) error {
	ctx := context.Background()
	langStr := sql.NullString{
		String: lang,
		Valid:  true,
	}

	params := db.SetLanguageParams{Language: langStr, Userid: userid}
	err := queries.SetLanguage(ctx, params)
	if err != nil {
		log.Println(err)
	}
	return err
}
