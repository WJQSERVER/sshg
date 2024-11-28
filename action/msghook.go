package action

import (
	"net/http"
	"net/url"
	"sshg/config"
)

// 发送到主控推送
func MsgHook(msg string, cfg *config.Config) {
	// 构造请求cfg.MsgHook.Url的请求参数
	parsedURL, err := url.Parse(cfg.MsgHook.URL)
	// 使用MSG-Auth 作为Header 鉴权
	if err != nil {
		// 解析失败，记录日志
		return
	}
	// 构造请求参数
	req := &http.Request{
		Method: "POST",
		URL:    parsedURL,
		Header: http.Header{
			"Content-Type": {"application/json"},
			"MSG-Auth":     {cfg.MsgHook.Token},
		},
		Body: http.NoBody,
	}
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// 发送失败，记录日志
		return
	}
	defer resp.Body.Close()
}

// 从主控获取消息
