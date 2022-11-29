package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	chatName string
}

func NewBot(token, chatName string) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Bot{bot: bot, chatName: chatName}
}

func (b *Bot) SendMessage(message string) bool {
	msg := tgbotapi.NewMessageToChannel(b.chatName, message)
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return false
	}
	return true
}
