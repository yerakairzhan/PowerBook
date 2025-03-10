package handlers

import (
	"PowerBook/api"
	db "PowerBook/db/sqlc"
	"PowerBook/utils"
	"context"
	"database/sql"
	"fmt"
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
				msg.ParseMode = "HTML"
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Error sending message", err)
				}
			}
		} else {
			_, text := utils.GetTranslation(ctx, queries, updates, "start_1")
			msg := tgbotapi.NewMessage(chatid, text)
			msg.ParseMode = "HTML"
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
		callbackMenu(queries, updates, bot, chatid)

	case "stat":

	case "top":
		callbackTop(bot, chatid, strconv.FormatInt(userid, 10), queries, updates)

	case "language":
		callbackLang(queries, updates, bot, chatid)

	case "read":
		callbackRead(queries, updates, bot, strconv.FormatInt(userid, 10), chatid)

	case "timer":
		callbackTimer(queries, updates, bot, chatid)
	}
}

func callbackMenu(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64) {
	ctx := context.Background()
	_, text := utils.GetTranslation(ctx, queries, updates, "menu")
	msg := tgbotapi.NewMessage(chatid, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = utils.InlineMenu()
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message", err)
	}
}

func callbackLang(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64) {
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

func callbackRead(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, userid string, chatid int64) {
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

func callbackRegister(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64, userid string) {
	ctx := context.Background()
	yes, _ := queries.GetRegistered(ctx, userid)
	if yes.Bool == false {
		utils.LoadConfig()
		if err := api.AddUserToSheet(utils.GoogleApi, userid, updates.CallbackQuery.From.UserName); err != nil {
			log.Fatalf("Error adding user to sheet: %v", err)
		}
	}
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
		callbackTimer(queries, updates, bot, chatid)
	}
}

func callbackTimer(queries *db.Queries, updates tgbotapi.Update, bot *tgbotapi.BotAPI, chatid int64) {
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

func callbackTop(bot *tgbotapi.BotAPI, chatID int64, userid string, queries *db.Queries, updates tgbotapi.Update) {
	ctx := context.Background()

	topReaders, err := queries.GetTopReadersThisMonth(ctx)
	if err != nil {
		log.Println("Error getting top readers:", err)
		return
	}

	var inlineButtons [][]tgbotapi.InlineKeyboardButton
	for i, reader := range topReaders {
		var medal string
		switch i {
		case 0:
			medal = "🥇" // золото
		case 1:
			medal = "🥈" // серебро
		case 2:
			medal = "🥉" // бронза
		default:
			medal = ""
		}
		button := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", medal), fmt.Sprintf("username_%s", reader.Username)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("@%s", reader.Username), fmt.Sprintf("username_%s", reader.Username)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d мин.", reader.TotalMinutes), fmt.Sprintf("minutes_%d", reader.TotalMinutes)),
		}
		inlineButtons = append(inlineButtons, button)
	}

	inlineMarkup := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)

	err, text := utils.GetTranslation(ctx, queries, updates, "top")
	if err != nil {
		log.Println("Error getting translation:", err)
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = inlineMarkup
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
	}

	topReader, err := queries.GetTopReaders(ctx)
	top := topReader[0]
	err, text = utils.GetTranslation(ctx, queries, updates, "top_1")
	if err != nil {
		log.Println("Error getting translation:", err)
	}
	msg = tgbotapi.NewMessage(chatID, text+"\n"+top.Username+" - "+strconv.FormatInt(top.TotalMinutes, 10))
	msg.ParseMode = "HTML"
	time.Sleep(1 * time.Second)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
	}

	topStreaks, err := queries.GetTopStreaks(ctx)
	if err != nil {
		log.Println("Error getting top streaks:", err)
	} else {
		var inlineButtons [][]tgbotapi.InlineKeyboardButton
		for i, reader := range topStreaks {
			var medal string
			switch i {
			case 0:
				medal = "🥇" // золото
			case 1:
				medal = "🥈" // серебро
			case 2:
				medal = "🥉" // бронза
			default:
				medal = ""
			}
			button := []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", medal), fmt.Sprintf("username_%s", reader.Username)),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("@%s", reader.Username), fmt.Sprintf("username_%s", reader.Username)),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d 🔥", reader.StreakLength), fmt.Sprintf("minutes_%d", reader.StreakLength)),
			}
			inlineButtons = append(inlineButtons, button)
		}

		inlineMarkup := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
		you, err := queries.GetUserTopStreak(ctx, userid)
		if err != nil {
			log.Println("Error getting user top streaks:", err)
		} else {
			err, text := utils.GetTranslation(ctx, queries, updates, "top_2")
			if err != nil {
				log.Println("Error getting translation:", err)
			}
			msg := tgbotapi.NewMessage(chatID, text+you+"🔥")
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = inlineMarkup
			_, err = bot.Send(msg)
			if err != nil {
				log.Println("Error sending message:", err)
			}
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
		msg := tgbotapi.NewMessage(chatid, text)
		msg.ParseMode = "HTML"
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err.Error())
		}
		time.Sleep(2 * time.Second)
		_, text = utils.GetTranslation(ctx, queries, updates, "menu")
		msg = tgbotapi.NewMessage(chatid, text)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = utils.InlineMenu()
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err.Error())
		}
	}
	return err
}
