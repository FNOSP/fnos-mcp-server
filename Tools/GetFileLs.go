package tools

import (
	"context"
	"fmt"
	"fn-mcp-server/wsclient"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type FileLsInput struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
	V    string `json:"v"`
}

func GetFileLs(ctx context.Context, req *mcp.CallToolRequest, input *FileLsInput) (*mcp.CallToolResult, any, error) {
	pass, err := check(req, true)
	if !pass {
		return nil, nil, err
	}
	identity := req.Session.ID()

	var wsClient *wsclient.FnOsWsBase
	wsClient, _ = connectionManager.GetConnection(identity)
	var path *string = nil
	if input.Name != "" && input.Uid != "" && input.V != "" {
		_path := fmt.Sprintf("vol%s/%s/%s", input.V, input.Uid, input.Name)
		path = &_path
	}
	result := wsClient.GetFileLs(path)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "获取成功",
			},
		},
	}, result.Files, nil
}
