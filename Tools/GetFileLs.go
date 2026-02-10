package tools

import (
	"context"
	"fn-mcp-server/wsclient"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type FileLsInput struct {
	Path string `json:"path"` // 改为大写并添加JSON标签，MCP SDK才能正确解析
}

func GetFileLs(ctx context.Context, req *mcp.CallToolRequest, input FileLsInput) (*mcp.CallToolResult, any, error) {
	pass, err := check(req, true)
	if !pass {
		return nil, nil, err
	}
	identity := req.Session.ID()

	var wsClient *wsclient.FnOsWsBase
	wsClient, _ = connectionManager.GetConnection(identity)
	// 由于我们将path从指针类型改为了值类型，所以需要直接使用input.Path
	path := input.Path
	result := wsClient.GetFileLs(&path)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "获取成功",
			},
		},
	}, result.Files, nil
}
