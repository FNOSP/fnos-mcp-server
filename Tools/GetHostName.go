package tools

import (
	"context"
	"fn-mcp-server/wsclient"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type EmptyInput struct{}

// GetHostName 获取飞牛设备的HostName和版本信息
func GetHostName(ctx context.Context, req *mcp.CallToolRequest, input EmptyInput) (*mcp.CallToolResult, any, error) {
	wsClient := wsclient.NewFnOsWsBase("http://fnos.xn--1jqw64a7tu.cn:25130/")
	err := wsClient.Start("main")
	if err != nil {
		return nil, nil, err
	}
	hostNameResult := wsClient.GetHostName()
	wsClient.Stop()

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "获取成功",
			},
		},
	}, hostNameResult.Data, nil
}
