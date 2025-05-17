// KeyCallBack/keyboard_callback.go
package key_callback

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
		// sendPhotoParams := &bot.SendPhotoParams{
		// 	ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		// 	Photo:  &models.InputFileUpload{Filename: data.ImageID + ".jpg", Data: file},
		// }
		media_path := fmt.Sprintf("attach://%s", path)
		media := &models.InputMediaPhoto{
			Media:           media_path,
			Caption:         "Новое фото",
			ParseMode:       models.ParseModeMarkdown,
			MediaAttachment: file, // ВАЖНО: сюда кладёшь io.Reader с файлом
		}

		time.Sleep(500 * time.Millisecond)
		fmt.Println("updating photo")
		editPhotoParams := &bot.EditMessageMediaParams{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Media:     media,
		}

		if _, err := b.EditMessageMedia(ctx, editPhotoParams); err != nil {
			log.Printf("Ошибка отправки фото: %v\n", err)
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
	}
}
