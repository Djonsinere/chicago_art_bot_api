// KeyCallBack/keyboard_callback.go
package key_callback

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	apicalls "art_chicago/api_calls"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var (
	numericKeyboard = models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "<", CallbackData: "<"},
				{Text: ">", CallbackData: ">"},
			},
		},
	}
	//userNumStates = make(map[int64]int)
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

		media_path := fmt.Sprintf("attach://%s", path)

		var let_arr [35]string
		for num, let := range data.Dimensions {
			if string(let) == "(" {
				break
			}
			let_arr[num] = string(let)
		}
		fixDemension := strings.Join(let_arr[:], "")

		caption_data := fmt.Sprintf("Автор: %s\nОписание: %s\nРазмеры: %s\nКлассификация: %s", data.ArtistTitle, data.CreditLine, fixDemension, data.Сlassification_title) //добавть дату создания
		fmt.Print(caption_data)
		media := &models.InputMediaPhoto{
			Media:           media_path,
			Caption:         caption_data,
			ParseMode:       models.ParseModeHTML,
			MediaAttachment: file,
		}

		time.Sleep(1000 * time.Millisecond)
		fmt.Println("updating photo")
		editPhotoParams := &bot.EditMessageMediaParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Media:       media,
			ReplyMarkup: &numericKeyboard,
		}

		if _, err := b.EditMessageMedia(ctx, editPhotoParams); err != nil {
			log.Printf("Ошибка отправки фото: %v\n", err)
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
	}
}
