// KeyCallBack/keyboard_callback.go
package key_callback

import (
	"context"
	"fmt"
	"log"
	"os"

	apicalls "art_chicago/api_calls"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleCallback(ctx context.Context, b *bot.Bot, update *models.Update, img_data [50]apicalls.ImageData) {
	for _, data := range img_data {
		if data.ID == 0 {
			break
		}
		path := fmt.Sprintf("%d/%s.jpg", update.CallbackQuery.From.ID, data.ImageID)
		file, err := os.Open(path)
		if err != nil {
			log.Printf("Ошибка открытия файла %s: %v\n", path, err)
			continue
		}
		defer file.Close()
		sendPhotoParams := &bot.SendPhotoParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID, //надо где то достать чат айди
			Photo:  &models.InputFileUpload{Filename: data.ImageID + ".jpg", Data: file},
		}

		if _, err := b.SendPhoto(ctx, sendPhotoParams); err != nil {
			log.Printf("Ошибка отправки фото: %v\n", err)
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
	}
}
