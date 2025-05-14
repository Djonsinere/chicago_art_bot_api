package main

import (
	"art_chicago/db"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

		// Check exist user and create

		id := strconv.FormatInt(update.Message.From.ID, 10)
		username := update.Message.Chat.UserName
		db.CreateNewUser(&id, &username)

		// Создаем копию полученного сообщения
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		reply.ReplyToMessageID = update.Message.MessageID

		// Отправляем ответ
		if _, err := bot.Send(reply); err != nil {
			log.Println("Ошибка при отправке сообщения:", err)
		}
	}
}
