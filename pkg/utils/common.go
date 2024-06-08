package utils

import "fmt"

var (
	LogLevel      = "info"
	DirectusHost  = "http://localhost:8055"
	DirectusToken = "directus-access-token"
	BotToken      = "my-bot-token"
)

func GetHelpMessage(username string) string {
	return fmt.Sprintf(`I am a Bot that moderates chat groups. Just add me into a group chat and give me permissions to send polls and delete messages. 

Summon me in the group chat using /delete@%v and reply to the message in question. I will then send a poll to collect other members' opinions. 

If the number of votes in favour of deleting the message >= certain threshold, I will close the poll and delete the message in question. Polls are only active for the expiry time the group admin sets, and requests will need to be resent.`,
		username,
	)
}

const SUPPORT_MESSAGE string = `My source code is hosted on https://github.com/Jason-CKY/telegram-modbot. Post any issues with this bot on the github link, and feel free to contribute to the source code with a pull request.`
const POLL_EXPIRY int = 120
const MIN_POLL_EXPIRY int = 10
const MAX_POLL_EXPIRY int = 600
