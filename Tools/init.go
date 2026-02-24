package tools

import "fn-mcp-server/wsclient"

type EmptyInput struct{}

// ConnectionManager 全局连接管理器实例，导出供其他包使用
var ConnectionManager *wsclient.ConnectionManager

// InitConnectionManager 初始化连接管理器
func InitConnectionManager() {
	ConnectionManager = wsclient.NewConnectionManager()
}
