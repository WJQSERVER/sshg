package action

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TGBot(chatID int64, token string, msgtext string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	//发送消息(无阻塞)(协程)
	go sendMsg(bot, chatID, msgtext)
}

// 发送消息
func sendMsg(bot *tgbotapi.BotAPI, chatID int64, msgtext string) {
	//markdown格式
	msg := tgbotapi.NewMessage(chatID, msgtext)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func TGBotListen(token string, chatID int64, port string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	// 监听消息(无阻塞)(协程)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		go ListenMsg(bot, update, port)
	}
}

// 监听消息
func ListenMsg(bot *tgbotapi.BotAPI, update tgbotapi.Update, port string) {
	if update.Message == nil {
		return
	}
	// 处理命令/ban ip
	if update.Message.IsCommand() && update.Message.Command() == "ban" {
		// 获取消息文本并去掉命令
		args := update.Message.CommandArguments()
		// 分割参数，按照空格
		ipList := strings.Fields(args)
		if len(ipList) == 1 {
			// 处理命令/ban ip
			err := Ban(ipList[0], port)
			if err == nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "已封禁IP "+ipList[0])
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "封禁IP "+ipList[0]+" 失败")
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			}
		} else {
			// 如果没有提供参数，发送错误消息或使用提示帮助
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "请提供要封禁的IP地址")
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		}
	}
	// 处理命令/unban ip
	if update.Message.IsCommand() && update.Message.Command() == "unban" {
		// 获取消息文本并去掉命令
		args := update.Message.CommandArguments()
		// 分割参数，按照空格
		ipList := strings.Fields(args)
		if len(ipList) == 1 {
			// 处理命令/unban ip
			err := Unban(ipList[0], port)
			if err == nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "已解封IP "+ipList[0])
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "解封IP "+ipList[0]+" 失败")
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			}
		} else {
			// 如果没有提供参数，发送错误消息或使用提示帮助
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "请提供要解封的IP地址")
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		}
	}
}
