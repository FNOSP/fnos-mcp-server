package sysinfo

import (
	"context"
	Tools "fn-mcp-server/Tools"
	"fn-mcp-server/wsclient"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetHardwareInfo 获取飞牛设备的硬件信息
func GetHardwareInfo(ctx context.Context, req *mcp.CallToolRequest, input *Tools.EmptyInput) (*mcp.CallToolResult, *wsclient.HardwareInfoData, error) {
	pass, err := Tools.Check(req, true)
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
	wsClient, _ = Tools.ConnectionManager.GetConnection(identity)
	hardwareInfo := wsClient.GetNetworkNetList()

	return nil, &hardwareInfo.Data, nil
}
