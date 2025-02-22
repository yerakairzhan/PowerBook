package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"time"
)

func InlineMenu() tgbotapi.InlineKeyboardMarkup {
	LoadConfig()

	inline := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Read", "callback_read"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📊 Table", TableURL),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏆 Top", "callback_top"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌍 Language", "callback_lang"),
		),
	)
	return inline
}

func InlineLang() tgbotapi.InlineKeyboardMarkup {
	inline := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "callback_ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "callback_en"),
			tgbotapi.NewInlineKeyboardButtonData("🇰🇿 Қазақша	", "callback_kz"),
		),
	)
	return inline
}

func InlineRegister() tgbotapi.InlineKeyboardMarkup {
	inline := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Регистрация", "callback_register"),
		),
	)
	return inline
}

func InlineTimer() tgbotapi.InlineKeyboardMarkup {
	inline := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🕒 15:00", "timer_15_00"),
			tgbotapi.NewInlineKeyboardButtonData("🕕 18:00", "timer_18_00"),
			tgbotapi.NewInlineKeyboardButtonData("🕘 21:00", "timer_21_00"),
		),
	)
	return inline
}

func GenerateCalendarKeyboard(year int, month int, readMinutes map[int]int) tgbotapi.InlineKeyboardMarkup {
	daysOfWeek := []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	months := []string{"December", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November"}
	var keyboard [][]tgbotapi.InlineKeyboardButton
	var weekRow []tgbotapi.InlineKeyboardButton

	for _, day := range daysOfWeek {
		weekRow = append(weekRow, tgbotapi.NewInlineKeyboardButtonData(day, "ignore"))
	}
	keyboard = append(keyboard, weekRow)

	// Get first weekday and total days in month
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	startWeekday := int(firstDay.Weekday()) // Sunday = 0
	if startWeekday == 0 {
		startWeekday = 7 // Adjust for Monday start
	}
	daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()

	var row []tgbotapi.InlineKeyboardButton
	// Fill empty slots before the first day
	for i := 1; i < startWeekday; i++ {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
	}

	// Fill in the actual days
	for day := 1; day <= daysInMonth; day++ {
		minutes := readMinutes[day]

		row = append(row, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(minutes), fmt.Sprintf("day_%d.%d.%d", day, month, year)))

		// Break at the end of each week
		if len(row) == 7 {
			keyboard = append(keyboard, row)
			row = nil
		}
	}

	// Add remaining row if not complete
	if len(row) > 0 {
		for len(row) < 7 {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "ignore"))
		}
		keyboard = append(keyboard, row)
	}

	navRow := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", fmt.Sprintf("calendar_%d_%d", year, month-1)),
		tgbotapi.NewInlineKeyboardButtonData("📆"+months[month], fmt.Sprintf("calendar_%d_%d", time.Now().Year(), int(time.Now().Month()))),
		tgbotapi.NewInlineKeyboardButtonData("➡️ Next", fmt.Sprintf("calendar_%d_%d", year, month+1)),
	}
	keyboard = append(keyboard, navRow)

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}
