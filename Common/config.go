package Common

import (
	"context"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aaronangxz/SeaDinner/Log"
)

var (
	Config     tomlConfig
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
	UrlPrefix   string `toml:"url_prefix"`
	TokenPrefix string `toml:"token_prefix"`
}

type tomlConfig struct {
	Adhoc     bool            `toml:"adhoc"`
	Prefix    prefixConfig    `toml:"prefix"`
	Runtime   runtimeConfig   `toml:"runtime"`
	OrderTime orderTimeConfig `toml:"order_time"`
	GrayScale grayScaleConfig `toml:"grayscale"`
}

func LoadConfig() {
	ConfigPath = "../config.toml"
	if os.Getenv("HEROKU_DEPLOY") == "TRUE" || os.Getenv("TEST_DEPLOY") == "TRUE" {
		ConfigPath = "config.toml"
	}

	if _, err := toml.DecodeFile(ConfigPath, &Config); err != nil {
		Log.Error(ctx, "Reading config failed | %v | %v", err, ConfigPath)
		// log.Fatalln("Reading config failed | ", err, ConfigPath)
		return
	}
	Log.Info(ctx, "Reading config OK | %v", ConfigPath)
	// log.Println("Reading config OK", ConfigPath)
}
