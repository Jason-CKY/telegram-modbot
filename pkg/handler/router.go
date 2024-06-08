package handler

import (
	"fmt"

	"github.com/Jason-CKY/telegram-modbot/pkg/schemas"
	"github.com/Jason-CKY/telegram-modbot/pkg/utils"
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	log.Info(update)
	numMembers, err := bot.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: update.FromChat().ID}})
	if err != nil {
		log.Error(err)
		return
	}
	chatSettings, _, err := schemas.InsertChatSettingsIfNotPresent(update.Message.Chat.ID, numMembers/2)
	if err != nil {
		log.Error(err)
		return
	}
	if update.Message != nil {
		if update.Message.IsCommand() {
			HandleCommand(update, bot, chatSettings)
		}
	}
	if update.Poll != nil {
		if update.Poll.IsClosed {
			for _, option := range update.Poll.Options {
				if option.Text == "Delete" && option.VoterCount >= chatSettings.Threshold {
					deleteMessage := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
					if _, err := bot.Request(deleteMessage); err != nil {
						log.Error(err)
						return
					}
				}
			}
		}
	}
}

func HandleCommand(update *tgbotapi.Update, bot *tgbotapi.BotAPI, chatSettings *schemas.ChatSettings) {
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
		if update.Message.ReplyToMessage == nil {
			return
		}
		var options = []string{"Delete", "Don't Delete"}
		question := fmt.Sprintf(`Poll to delete the message above. If >=%v of the group members vote to delete, the replied message shall be deleted.`, chatSettings.Threshold)
		pollConfig := tgbotapi.SendPollConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           update.Message.Chat.ID,
				ReplyToMessageID: update.Message.ReplyToMessage.MessageID,
			},
			Question:    question,
			Options:     options,
			IsAnonymous: true, // This is Telegram's default.
			OpenPeriod:  chatSettings.ExpiryTime,
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
