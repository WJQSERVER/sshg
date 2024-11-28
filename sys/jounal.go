package sys

import (
	"bufio"
	"fmt"
	"os/exec"
	"sshg/config"
	"sshg/core"
	"sshg/logger"
	"strings"
	"time"
)

// 日志模块
var (
	logw       = logger.Logw
	logInfo    = logger.LogInfo
	logWarning = logger.LogWarning
	logError   = logger.LogError
)

func CheckJounal(cfg *config.Config) {
	var lastReadTime time.Time

	for {
		// 如果是第一次运行，初始化 lastReadTime
		if lastReadTime.IsZero() {
			lastReadTime = time.Now()
		}

		// 使用 journalctl 读取上次读取时间到现在的 SSH 日志
		cmd := exec.Command("journalctl", "-u", "ssh.service", "--since", lastReadTime.Format("2006-01-02 15:04:05"))
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("Error creating StdoutPipe:", err)
			logError("Error creating StdoutPipe: %v", err)
			return
		}

		if err := cmd.Start(); err != nil {
			fmt.Println("Error starting command:", err)
			logError("Error starting command: %v", err)
			return
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			// 检查登录成功的日志
			if strings.Contains(line, "Accepted password for") {
				fmt.Println("SSH 登录成功:", line)
				logInfo("SSH 登录成功: %s", line)
				core.JournalCallback(0, line, cfg)
			}

			// 检查登录失败的日志
			if strings.Contains(line, "Failed password for") {
				fmt.Println("SSH 登录失败:", line)
				logWarning("SSH 登录失败: %s", line)
				core.JournalCallback(1, line, cfg)
			}

			// 检查连接关闭的日志
			if strings.Contains(line, "Connection closed by") {
				fmt.Println("SSH 连接关闭:", line)
				logInfo("SSH 连接关闭: %s", line)
				core.JournalCallback(100, line, cfg)
			}

			// 检查 kex_exchange_identification 错误
			if strings.Contains(line, "kex_exchange_identification") {
				fmt.Println("SSH 错误:", line)
				logError("SSH 错误: %s", line)
				core.JournalCallback(2, line, cfg)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from stdout:", err)
			logError("Error reading from stdout: %v", err)
		}

		if err := cmd.Wait(); err != nil {
			fmt.Println("Error waiting for command:", err)
			logError("Error waiting for command: %v", err)
		}

		// 更新 lastReadTime 为当前时间
		lastReadTime = time.Now()

		// 每隔 5 秒读取一次
		time.Sleep(5 * time.Second)
	}
}
