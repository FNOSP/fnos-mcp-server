package tools

import "fn-mcp-server/wsclient"

type EmptyInput struct{}

// 全局连接管理器实例
var connectionManager *wsclient.ConnectionManager

// InitConnectionManager 初始化连接管理器
func InitConnectionManager() {
	connectionManager = wsclient.NewConnectionManager()
}
