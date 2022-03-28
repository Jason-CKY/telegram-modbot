import os
from telegram.ext import ExtBot

BOT_TOKEN = os.getenv('BOT_TOKEN')
Bot = ExtBot(token=BOT_TOKEN)

# PUBLIC_URL = os.getenv('PUBLIC_URL')
CONFIG_COMMAND_MESSAGE = f"Get configs by typing '/getconfig'. \n" +\
                        f"Set your own threshold by typing '/setthreshold@{Bot.get_me().username} <number>'\n " +\
                        f"Set your own threshold by typing '/setexpiry@{Bot.get_me().username} <number>'"

START_MESSAGE = f"I am a Bot that moderates chat groups. Just add me into a group chat and " +\
                f"give me permissions to send polls and delete messages. Summon me in the " +\
                f"group chat using '/delete@{Bot.get_me().username}' and reply to the message in question. " +\
                f"I will then send a poll to collect other members' opinions. If the number of votes " +\
                f"in favour of deleting the message >= certain threshold, I will close the poll and delete the message in question." +\
                f"Polls are only active for the expiry time the group admin sets, and requests will need to be resent." +\
                CONFIG_COMMAND_MESSAGE

SUPPORT_MESSAGE =   f"My source code is hosted on https://github.com/Jason-CKY/Telegram-Bots/tree/main. Consider \n" +\
                    f"Post any issues with this bot on the github link, and feel free to contribute to the source code with a " +\
                    f"pull request."
POLL_EXPIRY = 120
MAX_EXPIRY = 600
MIN_EXPIRY = 10
DEV_CHAT_ID = os.getenv('DEV_CHAT_ID')

# import this line to avoid importing commands before
# defining the rest of the config as commands also import
# the configs from this file

from app import commands
'''
start - Help on how to use this Bot
help - Help on how to use this Bot
delete - Reply to a message with this command to initiate poll to delete
getconfig - Get current threshold and expiry time for this group chat
setthreshold - Set a threshold for this group chat
setexpiry - Set a expiry time for the poll
support - Support me on github!
'''
COMMANDS = {
    '/start': commands.start,
    '/help': commands.start,
    '/delete': commands.delete,
    '/getconfig': commands.get_config,
    '/setthreshold': commands.set_threshold,
    '/setexpiry': commands.set_expiry,
    '/support': commands.support
}
