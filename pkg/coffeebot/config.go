package coffeebot

import (
	"coffeeBot/internal/slackbot"
	"flag"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type DiscordConfig struct {
	AccessToken string `viper:"AccessToken"`
}

type NetConfig struct {
	Port int `viper:"Port"`
}

type Config struct {
	Net     NetConfig       `viper:"Net"`
	Slack   slackbot.Config `viper:"Slack"`
	Discord DiscordConfig   `viper:"Discord"`
}

func createViper() *viper.Viper {
	format := flag.String("format", "yaml", "Config file format. Valid values: 'yaml' (default) and 'json'")
	flag.Parse()
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.config")
	v.AddConfigPath("/etc/coffeebot")
	v.AddConfigPath("configs")
	v.SetConfigType(*format)
	v.SetConfigName("coffeebot")
	v.SetDefault("CoffeeBot.Net.Port", 80)
	return v
}

func readConfig() (_ Config, err error) {
	v := createViper()
	if err = v.ReadInConfig(); err != nil {
		return
	}

	cfg := struct {
		CoffeeBot Config `viper:"CoffeeBot"`
	}{}

	if err = v.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) {
		c.TagName = "viper"
	}); err != nil {
		return
	}

	return cfg.CoffeeBot, nil
}
