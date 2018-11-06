package main

import (
	"flag"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func newBool(val bool) *bool {
	b := val
	return &b
}

func main() {
	botToken := flag.String("token", "", "Bot API token")
	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(*botToken)
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.Sticker != nil {
			delMsgConfig := tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			}

			if _, err := bot.DeleteMessage(delMsgConfig); err != nil {
				log.Print(err)
				continue
			}

			restrictConfig := tgbotapi.RestrictChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: update.Message.Chat.ID,
					UserID: update.Message.From.ID,
				},
				UntilDate:             int64(update.Message.Date) + 60,
				CanSendMessages:       newBool(true),
				CanSendMediaMessages:  newBool(true),
				CanSendOtherMessages:  newBool(false),
				CanAddWebPagePreviews: newBool(true),
			}

			if _, err := bot.RestrictChatMember(restrictConfig); err != nil {
				log.Print(err)
				continue
			}
		}
	}
}
