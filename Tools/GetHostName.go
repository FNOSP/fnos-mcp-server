package tools

import (
	"context"
	"fn-mcp-server/wsclient"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetHostName 获取飞牛设备的HostName和版本信息
func GetHostName(ctx context.Context, req *mcp.CallToolRequest, input EmptyInput) (*mcp.CallToolResult, *wsclient.GetHostNameDataRetDto, error) {
	pass, err := Check(req, false)
	if !pass {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "获取失败" + err.Error(),
				},
			},
		}, nil, err
	}
	identity := req.Session.ID()

	var wsClient *wsclient.FnOsWsBase

	wsClient, _ = connectionManager.GetConnection(identity)

	hostNameResult := wsClient.GetHostName()

	return nil, &hostNameResult.Data, nil
}
