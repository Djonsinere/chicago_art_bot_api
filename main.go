// main.go
package main

import (
	key_callback "art_chicago/KeyCallBack"
	apicalls "art_chicago/api_calls"
	"art_chicago/db"
	"context"
	"log"
	"strconv"

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
	userStates  = make(map[int64]string)
	userReqData = make(map[int64][50]apicalls.ImageData)
	userCounter = make(map[int64]int)
)

func main() {
	botToken := Token
	if botToken == "" {
		log.Fatal("Token not set")
	}

	ctx := context.Background()
	db.Base_init_db()

	// Создаем обработчики

	b, err := bot.New(botToken,
		bot.WithDefaultHandler(defaultHandler),
	)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/search", bot.MatchTypeExact, baseSearchHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, defaultHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, callbackHandler)

	if err != nil {
		log.Panic(err)
	}

	b.Start(ctx)
}

// Обработчики команд
func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	sendText(b, update.Message.Chat.ID, "Этот бот позволяет искать произведения искусства, которые представленны в открытом доступе в Чикагском университете искусства(https://www.artic.edu). Максимальное количество изображений за 1 запрос = 50. Бот использует систему поиска full search определенную для Чикагского университета искусства, поэтому результаты такие какие есть. Если у вас есть какой либо вопрос или вы нашли баг - пишите @rayhartt")
}

func baseSearchHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userCounter[update.Message.From.ID] = 1
	userStates[update.Message.From.ID] = "awaiting_search"
	sendText(b, update.Message.Chat.ID, "Пожалуйста введите поисковой запрос на английском языке")
}

// Основной обработчик
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userID := strconv.FormatInt(update.Message.From.ID, 10)
	username := update.Message.From.Username
	go db.CreateNewUser(&userID, &username)

	if state, exists := userStates[update.Message.From.ID]; exists && state == "awaiting_search" {
		resp := apicalls.Full_text_search(update.Message.Text, update.Message.From.ID)
		userReqData[update.Message.From.ID] = resp

		sendMsgParams := &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Тут будет ответ",
			ReplyMarkup: &numericKeyboard,
		}
		b.SendMessage(ctx, sendMsgParams)
		delete(userStates, update.Message.From.ID)
	}
}

// Обработчик callback-запросов
func callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	press_button := update.CallbackQuery.Data
	if userCounter[update.CallbackQuery.From.ID] == 0 {
		userCounter[update.CallbackQuery.From.ID] = 1
	}
	switch press_button {
	case "<":
		if userCounter[update.CallbackQuery.From.ID] > 1 {
			userCounter[update.CallbackQuery.From.ID] -= 1
		}
	case ">":
		if userCounter[update.CallbackQuery.From.ID] < 50 {
			userCounter[update.CallbackQuery.From.ID] += 1
		}
	}
	key_callback.HandleCallback(ctx, b, update, userReqData[update.CallbackQuery.From.ID], userCounter[update.CallbackQuery.From.ID])
	//delete(userReqData, update.CallbackQuery.From.ID)
}

// Вспомогательная функция отправки текста
func sendText(b *bot.Bot, chatID int64, text string) {
	b.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
}
