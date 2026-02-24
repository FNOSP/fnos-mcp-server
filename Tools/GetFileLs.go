package tools

import (
	"context"
	"fmt"
	"fn-mcp-server/wsclient"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type FileLsInput struct {
	Path *string `json:"path,omitempty"` // 指针类型表示可选参数，omitempty表示为空时不序列化
	Uid  *string `json:"uid,omitempty"`
	V    *string `json:"v,omitempty"`
}

func GetFileLs(ctx context.Context, req *mcp.CallToolRequest, input *FileLsInput) (*mcp.CallToolResult, *wsclient.GetFileLsRetDto, error) {
	pass, err := check(req, true)
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
	var path *string = nil

	// 检查input是否为nil，以及所有字段是否都提供了值
	if input != nil && input.Path != nil && input.Uid != nil && input.V != nil {
		// 解引用指针获取实际值
		_path := fmt.Sprintf("vol%s/%s/%s", *input.V, *input.Uid, *input.Path)
		path = &_path
	}

	result := wsClient.GetFileLs(path)
	return nil, &result, nil
}
