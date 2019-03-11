package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type config struct {
	Debug bool
}

const configFileName = "src/config.yaml"

func newConfig() (*config, error) {
	// name of config file (without extension)
	viper.SetConfigName("config")
	viper.AddConfigPath("src")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "can't read config")
	}

	viper.AutomaticEnv()

	cfg := &config{
		Debug: viper.GetBool("debug"),
	}

	return cfg, nil
}
