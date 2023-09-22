package config

import (
	"github.com/pelletier/go-toml"
)

type Mode string

const (
	DEBUG   Mode = "debug"
	RELEASE Mode = "release"
)

type Config struct {
	TgToken     string
	Mode        Mode
	PsqlHost    string
	PsqlPort    int64
	PsqlUser    string
	PsqlPass    string
	PsqlDb      string
}

func Load(path string) (*Config, error) {

	cfg, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	return &Config{
		TgToken:     cfg.Get("bot.tg_token").(string),
		Mode:        Mode(cfg.Get("bot.mode").(string)),
		PsqlHost:    cfg.Get("postgres.host").(string),
		PsqlPort:    cfg.Get("postgres.port").(int64),
		PsqlUser:    cfg.Get("postgres.user").(string),
		PsqlPass:    cfg.Get("postgres.password").(string),
		PsqlDb:      cfg.Get("postgres.dbname").(string),
	}, nil
}
