package common

import (
	"context"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aaronangxz/SeaDinner/log"
)

var (
	//Config toml config object
	Config tomlConfig
	//ConfigPath path of config.toml
	ConfigPath string
	ctx        = context.TODO()
)

type grayScaleConfig struct {
	Percentage int64 `toml:"percentage"`
}

type orderTimeConfig struct {
	Hour    int `toml:"hour"`
	Minutes int `toml:"minutes"`
	Seconds int `toml:"seconds"`
}

type runtimeConfig struct {
	RetryTimes                 int `toml:"retry_times"`
	MenuRefreshIntervalSeconds int `toml:"menu_refresh_interval_seconds"`
}

type prefixConfig struct {
	URLPrefix   string `toml:"url_prefix"`
	TokenPrefix string `toml:"token_prefix"`
}

type tomlConfig struct {
	Adhoc     bool            `toml:"adhoc"`
	Prefix    prefixConfig    `toml:"prefix"`
	Runtime   runtimeConfig   `toml:"runtime"`
	OrderTime orderTimeConfig `toml:"order_time"`
	GrayScale grayScaleConfig `toml:"grayscale"`
}

//LoadConfig Loads config.toml
func LoadConfig() {
	ConfigPath = "../config.toml"
	if os.Getenv("HEROKU_DEPLOY") == "TRUE" || os.Getenv("TEST_DEPLOY") == "TRUE" {
		ConfigPath = "config.toml"
	}

	if _, err := toml.DecodeFile(ConfigPath, &Config); err != nil {
		log.Error(ctx, "Reading config failed | %v | %v", err, ConfigPath)
		return
	}
	log.Info(ctx, "Reading config OK | %v", ConfigPath)
}
