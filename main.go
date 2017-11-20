package main

import (
	"log"
	sc "github.com/maddevsio/simple-config"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	config := sc.NewSimpleConfig("./config", "yml")
	token := config.GetString("token")
	chatid := int64(config.Get("chatid").(int))

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	msg := tgbotapi.NewMessage(chatid, "Hello")
	bot.Send(msg)
}
