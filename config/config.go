package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server  ServerConfig
	TGBot   TGBotConfig
	Log     LogConfig
	MsgHook MsgHookConfig
}

type ServerConfig struct {
	Hostname string `toml:"hostname"`
	SSHPort  string `toml:"sshport"`
}

type TGBotConfig struct {
	Token  string `toml:"token"`
	ChatID int64  `toml:"chatid"`
}

type LogConfig struct {
	LogFilePath string `toml:"logfilepath"`
	MaxLogSize  int    `toml:"maxlogsize"`
}

type MsgHookConfig struct {
	Enabled bool   `toml:"enabled"`
	URL     string `toml:"url"`
	Token   string `toml:"token"`
}

// LoadConfig 从 TOML 配置文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
