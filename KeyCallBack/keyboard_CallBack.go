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

func HandleCallback(ctx context.Context, b *bot.Bot, update *models.Update, img_data [50]apicalls.ImageData, user_count int) {

	if img_data[user_count].ID == 0 {
		return
	}
	path := fmt.Sprintf("%d/%s.jpg", update.CallbackQuery.From.ID, img_data[user_count].ImageID)
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Ошибка открытия файла %s: %v\n", path, err)
	}
	defer file.Close()

	media_path := fmt.Sprintf("attach://%s", path)

	var let_arr [35]string
	for num, let := range img_data[user_count].Dimensions {
		if string(let) == "(" {
			break
		}
		let_arr[num] = string(let)
	}
	fixDemension := strings.Join(let_arr[:], "")

	caption_data := fmt.Sprintf("Автор: _%s_\nОписание: _%s_\nРазмеры: _%s_\nКлассификация: _%s_\nДата создания: _%s_", img_data[user_count].ArtistTitle, img_data[user_count].CreditLine, fixDemension, img_data[user_count].Сlassification_title, img_data[user_count].Date_display)
	fmt.Print(caption_data)
	media := &models.InputMediaPhoto{
		Media:           media_path,
		Caption:         caption_data,
		ParseMode:       models.ParseModeMarkdownV1,
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
