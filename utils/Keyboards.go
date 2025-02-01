package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func InlineMenu() tgbotapi.InlineKeyboardMarkup {
	inline := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Read", "callback_read"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Stat", "callback_stat"),
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
			tgbotapi.NewInlineKeyboardButtonData("Регистрироваться", "callback_register"),
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
