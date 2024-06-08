package handler

import (
	"github.com/Jason-CKY/telegram-modbot/pkg/utils"
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
		if update.Message.IsCommand() {
			HandleCommand(update, bot)
		}
	}
}

func HandleCommand(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	botUser, _ := bot.GetMe()
	// Extract the command from the Message.
	switch update.Message.Command() {
	case "help":
		msg.Text = utils.GetHelpMessage(botUser.String())
	case "start":
		msg.Text = utils.GetHelpMessage(botUser.String())
	case "support":
		msg.Text = utils.SUPPORT_MESSAGE
	case "delete":
		var options = []string{"Delete", "Don't Delete"}
		pollConfig := tgbotapi.SendPollConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: update.Message.Chat.ID,
			},
			Question:    "Poll to delete message.",
			Options:     options,
			IsAnonymous: true, // This is Telegram's default.
			OpenPeriod:  10,
		}
		if _, err := bot.Request(pollConfig); err != nil {
			log.Error(err)
			return
		}
		return
	default:
		return
	}

	if _, err := bot.Request(msg); err != nil {
		log.Error(err)
		return
	}
}
