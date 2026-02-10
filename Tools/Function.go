package tools

import (
	"encoding/base64"
	"fmt"
	"fn-mcp-server/wsclient"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func createWs(identity, FnOsUrl, Token string) error {
	wsClient := wsclient.NewFnOsWsBase(FnOsUrl)
	// 启动连接
	err := wsClient.Start("main")
	if err != nil {
		return err
	}
	if Token != "" {
		decodedBytes, err := base64.StdEncoding.DecodeString(Token)
		if err != nil {
			// 处理解码错误（如编码字符串无效、填充符错误等）
			log.Fatalf("标准 Base64 解码失败：%v", err)
			return err
		}
		// 将字节切片转为字符串（如果原内容是文本）
		decodedStr := string(decodedBytes)
		parts := strings.Split(decodedStr, ":")
		username := parts[0]
		password := parts[1]
		err = wsClient.Login(username, password)
		if err != nil {
			return err
		}
	}
	// 将连接添加到管理器
	connectionManager.AddConnection(identity, wsClient)

	return nil
}

func check(req *mcp.CallToolRequest, ifLogin bool) (bool, error) {
	identity := req.Session.ID()
	if identity == "" {
		return false, fmt.Errorf("身份验证失败: 缺少请求ID")
	}

	if connectionManager == nil {
		InitConnectionManager()
	}

	fnosUrl := req.Extra.Header.Get("Fn-Os-Url")
	if fnosUrl == "" {
		return false, fmt.Errorf("未获取到飞牛地址")
	}
	token := req.Extra.Header.Get("Token")
	if ifLogin {
		if token == "" {
			return false, fmt.Errorf("为获取到Token")
		}
	}
	var wsClient *wsclient.FnOsWsBase
	var exists bool

	wsClient, exists = connectionManager.GetConnection(identity)

	if !exists {
		err := createWs(identity, fnosUrl, token)
		if err != nil {
			return false, err
		}
		wsClient, exists = connectionManager.GetConnection(identity)
	}

	// 4. 检查连接状态
	wsClient.Mu.Lock()
	isRun := wsClient.IsRun
	wsClient.Mu.Unlock()

	if !isRun {
		// 连接已断开，从管理器中移除并创建新连接
		connectionManager.RemoveConnection(identity)

		err := createWs(identity, fnosUrl, token)
		if err != nil {
			return false, err
		}
		wsClient, exists = connectionManager.GetConnection(identity)
	}
	return true, nil
}
