package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Jason-CKY/telegram-modbot/pkg/schemas"
	"github.com/Jason-CKY/telegram-modbot/pkg/utils"
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
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
		if update.Message.IsCommand() {
			HandleCommand(update, bot, chatSettings)
		}
	}
	if update.Poll != nil {
		if update.Poll.IsClosed {
			pollWithChatSettings, err := schemas.GetPollWithChatSettingsByPollId(update.Poll.ID)
			if err != nil {
				log.Error(err)
				return
			}
			if pollWithChatSettings == nil {
				log.Error("no poll found")
				return
			}
			for _, option := range update.Poll.Options {
				if option.Text == "Delete" {
					if option.VoterCount >= pollWithChatSettings.ChatSettings.Threshold {
						poll := schemas.Poll{
							PollId: pollWithChatSettings.PollId,
						}
						err = poll.Delete()
						if err != nil {
							log.Error(err)
							return
						}
						deleteMessage := tgbotapi.NewDeleteMessage(pollWithChatSettings.ChatSettings.ChatId, pollWithChatSettings.MessageId)
						if _, err := bot.Request(deleteMessage); err != nil {
							log.Error(err)
							if strings.Contains(err.Error(), "can't be deleted") {
								msg := tgbotapi.NewMessage(pollWithChatSettings.ChatSettings.ChatId, "Error deleting message. Please check if I have the permission to delete group messages.")
								if _, err := bot.Request(msg); err != nil {
									log.Error(err)
									return
								}
							}
							return
						} else {
							msg := tgbotapi.NewMessage(pollWithChatSettings.ChatSettings.ChatId, "Offending message has been deleted.")
							if _, err := bot.Request(msg); err != nil {
								log.Error(err)
								return
							}
						}
					} else {
						msg := tgbotapi.NewMessage(pollWithChatSettings.ChatSettings.ChatId, "Threshold votes not reached before poll expiry.")
						if _, err := bot.Request(msg); err != nil {
							log.Error(err)
							return
						}
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
	case "getconfig":
		msg.Text = fmt.Sprintf(`Current group config:
Threshold: %v
Poll Expiry: %v seconds`, chatSettings.Threshold, chatSettings.ExpiryTime)
	case "setthreshold":
		shouldProcess := update.FromChat().IsPrivate()
		if !update.FromChat().IsPrivate() {
			admins, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: update.Message.Chat.ID}})
			if err != nil {
				log.Error(err)
				return
			}
			for _, admin := range admins {
				if admin.User.ID == update.Message.From.ID {
					shouldProcess = true
				}
			}
		}
		if !shouldProcess {
			msg.Text = "Only chat administrators allowed to set configs"
		} else {
			texts := strings.Split(update.Message.Text, " ")
			if len(texts) == 2 {
				threshold, err := strconv.Atoi(texts[1])
				if err != nil {
					log.Error(err)
					return
				}
				numMembers, err := bot.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: update.FromChat().ID}})
				if err != nil {
					log.Error(err)
					return
				}
				if threshold > numMembers {
					msg.Text = fmt.Sprintf("Invalid threshold %v more than members in the group.", threshold)
				} else if threshold < 1 {
					msg.Text = "Invalid threshold cannot be less than 1."
				} else {
					chatSettings.Threshold = threshold
					err = chatSettings.Update()
					if err != nil {
						msg.Text = fmt.Sprintf("Error setting threshold\n%v", err.Error())
						log.Error(err)
					} else {
						msg.Text = fmt.Sprintf("threshold set to %v", threshold)
					}
				}
			} else {
				return
			}
		}
	case "setexpiry":
		shouldProcess := update.FromChat().IsPrivate()
		if !update.FromChat().IsPrivate() {
			admins, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: update.Message.Chat.ID}})
			if err != nil {
				log.Error(err)
				return
			}
			for _, admin := range admins {
				if admin.User.ID == update.Message.From.ID {
					shouldProcess = true
				}
			}
		}
		if !shouldProcess {
			msg.Text = "Only chat administrators allowed to set configs"
		} else {
			texts := strings.Split(update.Message.Text, " ")
			if len(texts) == 2 {
				expiry, err := strconv.Atoi(texts[1])
				if err != nil {
					log.Error(err)
					return
				}
				if expiry > utils.MAX_POLL_EXPIRY {
					msg.Text = fmt.Sprintf("Invalid expiry %v cannot be more than %v.", expiry, utils.MAX_POLL_EXPIRY)
				} else if expiry < utils.MIN_POLL_EXPIRY {
					msg.Text = fmt.Sprintf("Invalid threshold cannot be less than %v.", utils.MIN_POLL_EXPIRY)
				} else {
					chatSettings.ExpiryTime = expiry
					err = chatSettings.Update()
					if err != nil {
						msg.Text = fmt.Sprintf("Error setting expiry time\n%v", err.Error())
						log.Error(err)
					} else {
						msg.Text = fmt.Sprintf("expiry set as %v seconds", expiry)
					}
				}
			} else {
				return
			}
		}
	case "delete":
		if update.Message.ReplyToMessage == nil {
			msg.Text = "Please make sure to reply to the offending message when making request to delete."
			if _, err := bot.Request(msg); err != nil {
				log.Error(err)
				return
			}
			return
		}
		var options = []string{"Delete", "Don't Delete"}
		question := fmt.Sprintf(`Poll to delete the message above.  This poll will last for %v seconds, if >=%v of the group members vote to delete, the replied message shall be deleted.`,
			chatSettings.ExpiryTime,
			chatSettings.Threshold)
		pollConfig := tgbotapi.SendPollConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           update.Message.Chat.ID,
				ReplyToMessageID: update.Message.ReplyToMessage.MessageID,
			},
			Question:    question,
			Options:     options,
			IsAnonymous: true, // This is Telegram's default.
		}
		pollResponse, err := bot.Request(pollConfig)
		if err != nil {
			log.Error(err)
			return
		}
		buf, err := pollResponse.Result.MarshalJSON()
		if err != nil {
			log.Error(err)
			return
		}
		var messageResponse tgbotapi.Message
		err = json.Unmarshal(buf, &messageResponse)
		if err != nil {
			log.Error(err)
			return
		}

		go func(update *tgbotapi.Update, bot *tgbotapi.BotAPI, chatSettings *schemas.ChatSettings, messageResponse tgbotapi.Message) {
			time.Sleep(time.Duration(chatSettings.ExpiryTime) * time.Second)
			endPollConfig := tgbotapi.NewStopPoll(update.Message.Chat.ID, messageResponse.MessageID)
			_, err := bot.Request(endPollConfig)
			if err != nil {
				log.Error(err)
				return
			}
		}(update, bot, chatSettings, messageResponse)

		poll := schemas.Poll{
			PollId:    messageResponse.Poll.ID,
			MessageId: update.Message.ReplyToMessage.MessageID,
			ChatId:    update.Message.Chat.ID,
		}
		err = poll.Create()
		if err != nil {
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
