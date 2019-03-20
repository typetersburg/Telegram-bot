package tg

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/viper"
)

// Config represents the variables that are needed to work with Telegram API.
type Config struct {
	Token                 string
	WebhookURL            string
	WebhookHostname       string
	WebhookMaxConnections int
	UpdateTimeout         int // in seconds
	UpdateLimit           int
	MessageWorkers        int
	Debug                 bool
}

// InitConfig initializes the configuration for telegram.
func (c *Config) InitConfig() {
	// Default values.
	viper.SetDefault("tg.message_workers", messagesPerSecondTgLimit)

	// Init keys for variables.
	c.Token = viper.GetString("tg.token")
	c.WebhookURL = viper.GetString("tg.webhook_url")
	c.WebhookHostname = viper.GetString("tg.webhook_hostname")
	c.WebhookMaxConnections = viper.GetInt("tg.webhook_max_connections")
	c.UpdateTimeout = viper.GetInt("tg.update_timeout")
	c.UpdateLimit = viper.GetInt("tg.update_limit")
	c.MessageWorkers = viper.GetInt("tg.message_workers")
	c.Debug = viper.GetBool("tg.debug")
}

// Validate validates the telegram configuration.
func (c Config) Validate() error {
	err := validation.ValidateStruct(&c,
		validation.Field(&c.Token, validation.Required),
	)
	if c.isUseTgWebhooks() {
		err = validation.ValidateStruct(&c,
			validation.Field(&c.WebhookURL, is.URL, validation.Required),
			validation.Field(&c.WebhookHostname, is.Host, validation.Required),
		)
	}
	return err
}

// isUseTgWebhooks checks whether variables are set to use the telegram API
// in webhooks mode.
func (c Config) isUseTgWebhooks() bool {
	return c.WebhookURL != "" || c.WebhookHostname != ""
}
