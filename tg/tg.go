package tg

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// messagesPerSecondTgLimit is the message limit per second of the Telegram
// Bot API.
const messagesPerSecondTgLimit = 30

// New creates a new BotAPI instance.
func (c Config) New() (*tg.BotAPI, error) {
	bot, err := tg.NewBotAPI(c.Token)
	if err != nil {
		return nil, errors.Wrap(err, "can't init connect to tg bot api")
	}

	bot.Debug = c.Debug

	return bot, nil
}
