package Processors

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	Config     tomlConfig
	ConfigPath string
)

type OrderTimeConfig struct {
	Hour    int `toml:"hour"`
	Minutes int `toml:"minutes"`
	Seconds int `toml:"seconds"`
}
type runtimeConfig struct {
	RetryTimes int `toml:"retry_times"`
}

type prefixConfig struct {
	UrlPrefix   string `toml:"url_prefix"`
	TokenPrefix string `toml:"token_prefix"`
}

type tomlConfig struct {
	Adhoc     bool            `toml:"adhoc"`
	Prefix    prefixConfig    `toml:"prefix"`
	Runtime   runtimeConfig   `toml:"runtime"`
	OrderTime OrderTimeConfig `toml:"order_time"`
}

func LoadConfig() {
	ConfigPath = "config.toml"
	if os.Getenv("HEROKU_DEPLOY") == "FALSE" {
		ConfigPath = "../config.toml"
	}

	if _, err := toml.DecodeFile(ConfigPath, &Config); err != nil {
		log.Fatalln("Reading config failed", err, ConfigPath)
	}
}
