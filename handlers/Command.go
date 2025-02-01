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

func handleCommand(command string, queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64) {
	userid := updates.Message.From.ID
	username := updates.Message.From.UserName
	name := updates.Message.From.FirstName

	ctx := context.Background()
	switch command {
	case "start":
		_, err := queries.GetUser(ctx, strconv.FormatInt(userid, 10))
		if err != nil {
			params := db.CreateUserParams{Userid: strconv.FormatInt(userid, 10), Username: username}
			_, err := queries.CreateUser(ctx, params)
			if err != nil {
				log.Println("Error creating user", err)
			} else {
				_, text := utils.GetTranslation(ctx, queries, updates, "start")
				msg := tgbotapi.NewMessage(chatid, text+name)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Error sending message", err)
				}
			}
		} else {
			_, text := utils.GetTranslation(ctx, queries, updates, "start_1")
			msg := tgbotapi.NewMessage(chatid, text)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Error sending message", err)
			}
		}

		time.Sleep(1 * time.Second)
		_, text := utils.GetTranslation(ctx, queries, updates, "start_2")
		msg := tgbotapi.NewMessage(chatid, text)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = utils.InlineRegister()
		_, err = bot.Send(msg)
		if err != nil {
			log.Println("Error sending message", err)
		}

		time.Sleep(1 * time.Second) //TODO на старте теперь надо добавить сет таймер

	case "menu":
		callback_menu(queries, updates, bot)

	case "stat":

	case "top":

	case "language":
		callback_lang(queries, updates, bot, chatid)

	case "read":
		callback_read(queries, updates, bot, strconv.FormatInt(userid, 10), chatid)

	case "timer":
		callback_timer(queries, updates, bot, chatid)
	}
}

func callback_menu(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatid := updates.Message.Chat.ID
	ctx := context.Background()
	_, text := utils.GetTranslation(ctx, queries, updates, "menu")
	msg := tgbotapi.NewMessage(chatid, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = utils.InlineMenu()
	bot.Send(msg)
}

func callback_lang(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64) {
	ctx := context.Background()
	err, text := utils.GetTranslation(ctx, queries, updates, "lang")
	if err != nil {
		log.Println("Error getting translation", err)
	}
	msg := tgbotapi.NewMessage(chatid, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = utils.InlineLang()
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Error sending message", err)
	}

	//TODO: change the language logic
}

func callback_read(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, userid string, chatid int64) {
	log.Println(chatid, userid)
	ctx := context.Background()
	err, text := utils.GetTranslation(ctx, queries, updates, "read")
	if err != nil {
		log.Println("Error getting translation", err)
	}
	msg := tgbotapi.NewMessage(chatid, text)
	msg.ParseMode = "HTML"
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Error sending message", err)
	}

	state := sql.NullString{
		String: "waiting_for_reading_time",
		Valid:  true,
	}

	params := db.SetUserStateParams{
		State:  state,
		Userid: userid,
	}

	err = queries.SetUserState(ctx, params)
	if err != nil {
		log.Println("Error setting user state:", err)
	}
}

func callback_register(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid string) {
	ctx := context.Background()
	err := queries.SetRegistered(ctx, userid)
	if err != nil {
		log.Println("Error setting user registered:", err)
	} else {
		err, text := utils.GetTranslation(ctx, queries, updates, "register")
		if err != nil {
			log.Println("Error getting translation", err)
		}
		msg := tgbotapi.NewMessage(chatid, text)
		msg.ParseMode = "HTML"
		_, err = bot.Send(msg)
		if err != nil {
			log.Println("Error sending message", err)
		}
		callback_timer(queries, updates, bot, chatid)
	}
}

func callback_timer(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64) {
	ctx := context.Background()
	time.Sleep(1 * time.Second)
	_, text := utils.GetTranslation(ctx, queries, updates, "register_2")
	msg := tgbotapi.NewMessage(chatid, text)
	msg.ReplyMarkup = utils.InlineTimer()
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message", err)
	}
}
