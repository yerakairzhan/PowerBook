package handlers

import (
	db "PowerBook/db/sqlc"
	"PowerBook/utils"
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"time"
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

func changeTimer(queries *db.Queries, userid string, bot *tgbotapi.BotAPI, updates tgbotapi.Update, chatid int64, hour int) error {
	ctx := context.Background()

	now := time.Now()
	timerValue := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
	params := db.SetTimerParams{
		Userid: userid,
		Timer:  timerValue,
	}
	err := queries.SetTimer(ctx, params)
	if err != nil {
		log.Println(err.Error())
		_, text := utils.GetTranslation(ctx, queries, updates, "timer_1")
		SendMessage(bot, chatid, text)
	} else {
		_, text := utils.GetTranslation(ctx, queries, updates, "timer")
		SendMessage(bot, chatid, text)
	}
	return err
}
