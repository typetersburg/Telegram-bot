package main

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/typetersburg/telegram-bot/tg"
)

type config struct {
	Tg tg.Config
}

func newConfig() (*config, error) {
	// name of config file (without extension)
	viper.SetConfigName("config")
	viper.AddConfigPath("src")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "can't read config")
	}

	cfg := &config{}
	cfg.Tg.InitConfig()

	err = cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "config is invalid")
	}

	return cfg, nil
}

func (c config) validate() error {
	err := c.Tg.Validate()
	if err != nil {
		return errors.Wrap(err, "Invalid tg config")
	}

	return nil
}
