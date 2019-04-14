package tg

import (
	"log"
	"net/http"
	"sync"
	"time"

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

// workerContext provides context for workers.
type workerContext struct {
	id  int
	wg  *sync.WaitGroup
	bot *tg.BotAPI
}

// InitMessageWorkerPool creates a pool of workers who send messages to
// Telegram based on the API limit.
func (c Config) InitMessageWorkerPool(bot *tg.BotAPI) (
	*sync.WaitGroup, chan<- tg.Chattable) {
	// limiter to limit sending messages so that Telegram does not reject messages.
	limitTicker := time.NewTicker(time.Second / messagesPerSecondTgLimit)

	// Wait group for graceful shutdown.
	wg := &sync.WaitGroup{}
	wg.Add(c.MessageWorkers)

	workerCtx := workerContext{
		wg:  wg,
		bot: bot,
	}

	// Channel for sending messages.
	ch := make(chan tg.Chattable, messagesPerSecondTgLimit)

	for i := 0; i < c.MessageWorkers; i++ {
		workerCtx.id = i
		go messageWorker(workerCtx, ch, limitTicker)
	}

	return wg, ch
}

func messageWorker(
	ctx workerContext,
	msgCh <-chan tg.Chattable,
	limitTicker *time.Ticker,
) {
	defer ctx.wg.Done()
	for {
		c, ok := <-msgCh
		// Check if the channel is closed. The channel closes when the
		// application ends.
		if !ok {
			log.Printf("[TG][DEBUG] worker %d closed\n", ctx.id)
			return
		}

		// Rate limiter
		<-limitTicker.C

		_, err := ctx.bot.Send(c)
		if err != nil {
			log.Printf("[TG][ERROR] can't send message: %v\n", err)
		}
	}
}
