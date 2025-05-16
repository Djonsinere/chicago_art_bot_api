package keycallback

import (
	apicalls "art_chicago/api_calls"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, chat_id int64, user_id int64, img_data [50]apicalls.ImageData) {
	for _, data := range img_data {
		if data.ID == 0 {
			break
		}
		path := fmt.Sprintf("%d/%s.jpg", user_id, data.ImageID)
		photo := tgbotapi.NewPhoto(chat_id, tgbotapi.FilePath(path))
		if _, err := bot.Send(photo); err != nil {
			log.Printf("\nОшибка отправки фото: %v\n", err)
		}
		// Подтверждаем получение callback
		bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	}
}
