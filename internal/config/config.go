package config

import (
	"github.com/pelletier/go-toml"
)

type Config struct {
	TgToken    string
	ModelsURL  string
	GalleryURL string
	PsqlHost   string
	PsqlPort   int64
	PsqlUser   string
	PsqlPass   string
	PsqlDb     string
}

func Load(path string) (*Config, error) {

	cfg, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	return &Config{
		TgToken:    cfg.Get("bot.tg_token").(string),
		ModelsURL:  cfg.Get("bot.models_url").(string),
		GalleryURL: cfg.Get("bot.gallery_url").(string),
		PsqlHost:   cfg.Get("postgres.host").(string),
		PsqlPort:   cfg.Get("postgres.port").(int64),
		PsqlUser:   cfg.Get("postgres.user").(string),
		PsqlPass:   cfg.Get("postgres.password").(string),
		PsqlDb:     cfg.Get("postgres.dbname").(string),
	}, nil
}
