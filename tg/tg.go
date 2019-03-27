package tg

import (
	"net/http"

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

// GetChanWithUpdates receives updates from the Telegram API, or updates come
// through the webhook.
func (c Config) GetChanWithUpdates(bot *tg.BotAPI) (<-chan tg.Update, error) {
	// Remove webhook if already crated
	if _, err := bot.RemoveWebhook(); err != nil {
		return nil, errors.Wrap(err, "can't remove webhook")
	}

	// If the bot will be launched using update retrieval, not via webhook.
	if !c.isUseTgWebhooks() {
		upd := tg.NewUpdate(0)
		upd.Timeout = c.UpdateTimeout
		upd.Limit = c.UpdateLimit
		return bot.GetUpdatesChan(upd)
	}

	webhookConfig := tg.NewWebhook(c.WebhookURL + c.Token)
	webhookConfig.MaxConnections = c.WebhookMaxConnections

	_, err := bot.SetWebhook(webhookConfig)
	if err != nil {
		return nil, errors.Wrap(err, "can't set webhook "+c.WebhookURL)
	}

	go func() {
		err = http.ListenAndServe(c.WebhookHostname, nil)
	}()
	if err != nil {
		return nil, errors.Wrap(err, "can't init server "+c.WebhookHostname)
	}

	return bot.ListenForWebhook("/" + bot.Token), nil
}
