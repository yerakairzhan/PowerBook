package main

import (
	db "PowerBook/db/sqlc"
	"PowerBook/handlers"
	"PowerBook/utils"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	utils.LoadConfig()
	err := utils.LoadTranslations()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(utils.DBDriver, utils.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	bot, err := tgbotapi.NewBotAPI(utils.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	db := db.New(conn)

	handlers.SetupHandlers(bot, db)
}
