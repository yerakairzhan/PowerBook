package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var BotToken string
var DBDriver string
var DBSource string

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла : ", err)
	}

	BotToken = os.Getenv("BOT_TOKEN")
	DBDriver = os.Getenv("DB_DRIVER")
	DBSource = os.Getenv("DB_SOURCE")
}
