package core

import (
	"errors"
	"fmt"
	"regexp"
	"sshg/action"
	"sshg/config"
	"sshg/logger"
	"time"
)

var (
	logw       = logger.Logw
	logInfo    = logger.LogInfo
	logWarning = logger.LogWarning
	logError   = logger.LogError
)

func Core() (err error) {
	return nil
}

func Callback(code int, msg string, cfg *config.Config) error {
	/*	switch code {
		case 0:
			logInfo(msg)
			fmsg, _, err := formatMsg(msg, 0, cfg)
			if err != nil {
				return err
			}
			action.TGBot(cfg.TGBot.ChatID, cfg.TGBot.Token, fmsg)
		case 1:
			fmsg, _, err := formatMsg(msg, 1, cfg)
			if err != nil {
				return err
			}
			action.TGBot(cfg.TGBot.ChatID, cfg.TGBot.Token, fmsg)
		case 2:
			fmsg, _, err := formatMsg(msg, 2, cfg)
			if err != nil {
				return err
			}
			action.TGBot(cfg.TGBot.ChatID, cfg.TGBot.Token, fmsg)
		case 100:
			fmsg, _, err := formatMsg(msg, 100, cfg)
			if err != nil {
				return err
			}
			action.TGBot(cfg.TGBot.ChatID, cfg.TGBot.Token, fmsg)

		default:
			return errors.New("invalid code")
		}
		return nil*/
	vaildCode := []int{0, 1, 2, 100}
	if !contains(vaildCode, code) {
		return errors.New("invalid code")
	}

	if code == 0 {
		logInfo(msg)
	}

	fmsg, _, err := formatMsg(msg, code, cfg)
	if err != nil {
		return err
	}
	logw(fmsg)
	err = action.TGBot(cfg.TGBot.ChatID, cfg.TGBot.Token, fmsg)
	if err != nil {
		return err
	}
	return nil
}

// 检查切片中是否包含指定的元素
func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func ToBot(code int, msg string, cfg *config.Config) error {
	vaildCode := []int{0, 1, 2, 100}
	if !contains(vaildCode, code) {
		return errors.New("invalid code")
	}

	if code == 0 {
		logInfo(msg)
	}

	fmsg, _, err := formatMsg(msg, code, cfg)
	if err != nil {
		return err
	}
	logw(fmsg)
	err = action.TGBot(cfg.TGBot.ChatID, cfg.TGBot.Token, fmsg)
	if err != nil {
		return err
	}
	return nil
}

func formatMsg(msg string, code int, cfg *config.Config) (fmsg string, ip string, err error) {
	var pattern string
	switch code {
	case 0:
		// 匹配 Accepted 日志
		pattern = `(\w+\s+\d+\s+\d+:\d+:\d+).*?Accepted password for (\w+) from ([\w.:]+) port (\d+)`
	case 1:
		// 匹配 Failed 日志
		pattern = `(\w+\s+\d+:\d+:\d+).*?Failed password for (\w+) from ([\w.:]+) port (\d+)`
	case 2:
		// 匹配 kex_exchange_identification 日志
		pattern = `(\w+\s+\d+\s+\d+:\d+:\d+).*?kex_exchange_identification: (\w+)`

	case 100:
		pattern = `(\w+\s+\d+\s+\d+:\d+:\d+).*?Connection closed by ([\w.:]+) port (\d+)`
	default:
		return "", "", errors.New("invalid code")
	}

	// 编译正则表达式
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(msg)

	// 检查匹配结果
	if len(matches) < 2 {
		return "", "", errors.New("log format not recognized")
	}

	// 解析时间并格式化
	timestamp := matches[1]
	layout := "Jan 2 15:04:05"
	parsedTime, err := time.Parse(layout, timestamp)
	if err != nil {
		return "", "", err
	}

	// 格式化时间为 MM-DD HH:MM:SS
	formattedTime := parsedTime.Format("01-02 15:04:05")

	// 构建结果字符串
	var (
		result    string
		ipAddress string
	)
	if code == 2 {
		result = fmt.Sprintf("主机 %s \n**Time:** %s \n**kex_exchange_identification:** %s", cfg.Server.Hostname, formattedTime, matches[2])
		logWarning(result)
	} else if code == 100 {
		ipAddress := matches[2]
		port := matches[3]
		result = fmt.Sprintf("主机 %s \n**Time:** %s \n**Connection closed by** %s **port** %s", cfg.Server.Hostname, formattedTime, ipAddress, port)
		logWarning(result)
	} else if code == 0 {
		username := matches[2]
		ipAddress := matches[3]
		port := matches[4]
		//result = fmt.Sprintf("主机 %s \n**Time:** %s \n**Login as** %s from [%s]:%s \n**Login Succeeded**", cfg.Server.Hostname, formattedTime, username, ipAddress, port)
		result = fmt.Sprintf("**主机** %s \n**Time:** %s \n**Login as** %s from [%s]:%s \n**Login Succeeded**\n![👍](tg://emoji?id=5368324170671202286)", cfg.Server.Hostname, formattedTime, username, ipAddress, port)
		logInfo(result)
	} else if code == 1 {
		username := matches[2]
		ipAddress := matches[3]
		port := matches[4]
		result = fmt.Sprintf("主机 %s \n**Time:** %s \n**Login as** %s from [%s]:%s \n**Login Failed**", cfg.Server.Hostname, formattedTime, username, ipAddress, port)
		logWarning(result)
	}
	return result, ipAddress, nil
}
