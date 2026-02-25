package sysinfo

import (
	"context"
	Tools "fn-mcp-server/Tools"
	"fn-mcp-server/wsclient"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetMachineId 获取飞牛设备的机器ID
func GetMachineId(ctx context.Context, req *mcp.CallToolRequest, input any) (*mcp.CallToolResult, *wsclient.MachineIdDto, error) {
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
	result := wsClient.GetMachineId()
	return nil, &result.Data, nil
}
