package main

import (
	apicalls "art_chicago/api_calls"
	"art_chicago/db"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Хранилище состояний пользователей
var userStates = make(map[int64]string)

func main() {
	// Получаем токен бота из переменной окружения
	botToken := Token
	if botToken == "" {
		log.Fatal("Token not set")
	}

	// Создаем экземпляр бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// Настраиваем канал для получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Инициалиируем базу данных
	db.Base_init_db()

	updates := bot.GetUpdatesChan(u)

	// Обрабатываем входящие обновления
	for update := range updates {
		if update.Message == nil { // Игнорируем всё, кроме сообщений
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		// Check exist user and create
		userID := strconv.FormatInt(update.Message.From.ID, 10)
		username := update.Message.Chat.UserName
		go db.CreateNewUser(&userID, &username)

		user_int_id := update.Message.From.ID
		if state, exists := userStates[user_int_id]; exists {
			switch state {
			case "awaiting_search":
				resp := apicalls.Full_text_search(update.Message.Text, user_int_id)
				if len(resp) > 4000 {
					resp = resp[:4000] + "... [ответ сокращен]"
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
				msg.ReplyToMessageID = update.Message.MessageID

				if _, err := bot.Send(msg); err != nil {
					log.Printf("Ошибка отправки: %v", err)
				}

				delete(userStates, user_int_id)
				continue
			}
		}
		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "Этот бот позволяет искать любые произведения искусства, которые хранятся в Чикагском университете искусств. По любымы багам/вопросам пишите @rayhartt"
		case "base_search":
			userStates[user_int_id] = "awaiting_search"
			msg.Text = "Пожалуйста введите поисковой запрос на английском языке"
		default:
			msg.Text = "Такой команды нет"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

	}
}
