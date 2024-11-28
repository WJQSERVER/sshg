package main

import (
	"flag"
	"fmt"
	"log"

	"sshg/action"
	"sshg/config"
	"sshg/logger"
	"sshg/sys"
)

var (
	cfg        *config.Config
	configfile = "/data/go/config/config.toml"
)

// 日志模块
var (
	logw       = logger.Logw
	logInfo    = logger.LogInfo
	logWarning = logger.LogWarning
	logError   = logger.LogError
)

func ReadFlag() {
	cfgfile := flag.String("cfg", configfile, "config file path")
	flag.Parse()
	configfile = *cfgfile
}

func loadConfig() {
	var err error
	// 初始化配置
	cfg, err = config.LoadConfig(configfile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Loaded config: %v\n", cfg)
}

func setupLogger() {
	// 初始化日志模块
	var err error
	err = logger.Init(cfg.Log.LogFilePath, cfg.Log.MaxLogSize) // 传递日志文件路径
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logw("Logger initialized")
	logw("Init Completed")
}

func init() {
	ReadFlag()
	loadConfig()
	setupLogger()
}

func main() {
	//sys.CheckJounal(cfg)
	go func() {
		err := action.TGBotListen(cfg.TGBot.Token, cfg.TGBot.ChatID, cfg.Server.SSHPort)
		if err != nil {
			logError("Failed to start TGBot: %v", err)
		}
	}()
	sys.CheckJounal(cfg)
	defer logger.Close() // 确保在退出时关闭日志文件
}
