package Processors

import (
	"log"

	"github.com/BurntSushi/toml"
)

var (
	Config     tomlConfig
	ConfigPath string
)

type runtimeConfig struct {
	RetryTimes int `toml:"retry_times"`
}

type prefixConfig struct {
	UrlPrefix   string `toml:"url_prefix"`
	TokenPrefix string `toml:"token_prefix"`
}

type tomlConfig struct {
	Prefix  prefixConfig  `toml:"prefix"`
	Runtime runtimeConfig `toml:"runtime"`
}

func loadConfig() {
	ConfigPath = "./config.toml"
	if _, err := toml.DecodeFile(ConfigPath, &Config); err != nil {
		log.Fatalln("Reading config failed", err)
	}
}
