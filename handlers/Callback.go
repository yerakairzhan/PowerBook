package handlers

import (
	db "PowerBook/db/sqlc"
	"PowerBook/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

func handleCallback(command string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64) {
	userid := strconv.FormatInt(updates.CallbackQuery.From.ID, 10)

	ctx := context.Background()
	switch {
	case "callback_register" == command:
		userid := updates.CallbackQuery.From.ID
		callbackRegister(queries, updates, bot, chatid, strconv.FormatInt(userid, 10))
	case "callback_read" == command:
		callbackRead(queries, updates, bot, userid, chatid)
	case "callback_stat" == command:
	case "callback_top" == command:
		messageID := updates.CallbackQuery.Message.MessageID
		callbackTop(bot, chatid, userid, messageID, queries, updates)
	case "callback_lang" == command:
		callbackLang(queries, updates, bot, chatid)
	case "callback_ru" == command:
		err := changeLang(queries, userid, "ru")
		if err != nil {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_2")
			SendMessage(bot, chatid, text)
		} else {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_1")
			SendMessage(bot, chatid, text)
			callbackMenu(queries, updates, bot, chatid)
		}
	case "callback_en" == command:
		err := changeLang(queries, userid, "en")
		if err != nil {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_2")
			SendMessage(bot, chatid, text)
		} else {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_1")
			SendMessage(bot, chatid, text)
			callbackMenu(queries, updates, bot, chatid)
		}
	case "callback_kz" == command:
		err := changeLang(queries, userid, "kz")
		if err != nil {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_2")
			SendMessage(bot, chatid, text)
		} else {
			_, text := utils.GetTranslation(ctx, queries, updates, "lang_1")
			SendMessage(bot, chatid, text)
			callbackMenu(queries, updates, bot, chatid)
		}

	//case "timer_15_00" == command:
	//	changeTimer(queries, userid, bot, updates, chatid)
	//	timerTime, _ := time.Parse("15:04", "15:00")
	//	params := db.SetTimerParams{Timer: timerTime, Userid: userid}
	//	queries.SetTimer(ctx, params)
	//case "timer_18_00" == command:
	//	changeTimer(queries, userid, bot, updates, chatid)
	//	timerTime, _ := time.Parse("15:04", "18:00")
	//	params := db.SetTimerParams{Timer: timerTime, Userid: userid}
	//	queries.SetTimer(ctx, params)
	//case "timer_21_00" == command:
	//	changeTimer(queries, userid, bot, updates, chatid)
	//	timerTime, _ := time.Parse("15:04", "21:00")
	//	params := db.SetTimerParams{Timer: timerTime, Userid: userid}
	//	queries.SetTimer(ctx, params)
	case strings.HasPrefix(command, "calendar_"):
		var year, month int
		_, err := fmt.Sscanf(command, "calendar_%d_%d", &year, &month)
		if err != nil {
			log.Println("Error parsing command:", err, "Command:", command)
			return
		}

		messageID := updates.CallbackQuery.Message.MessageID
		sendCalendar(chatid, userid, year, month, queries, bot, updates, true, messageID)
	case strings.HasPrefix(command, "day_"):
		var text string
		_, err := fmt.Sscanf(command, "day_%v", &text)
		if err != nil {
			log.Println("Error parsing command:", err, "Command:", command)
			return
		}
		callbackQueryID := updates.CallbackQuery.ID
		callback := tgbotapi.NewCallback(callbackQueryID, text)
		bot.Request(callback)

	default:
		log.Println("Unknown command:", command)
	}
}
