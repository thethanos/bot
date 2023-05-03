package config

import (
	"github.com/pelletier/go-toml"
)

type Config struct {
	TgToken string
}

func Load(path string) (*Config, error) {

	cfg, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	return &Config{
		TgToken: cfg.Get("bot.tg_token").(string),
	}, nil
}
