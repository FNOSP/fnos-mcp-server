package tools

import (
	"context"
	"fmt"
	"fn-mcp-server/wsclient"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type FileLsInput struct {
	Path *string `json:"path,omitempty" jsonschema:"飞牛系统的文件路径，以/开头"`
	Uid  *string `json:"uid,omitempty" jsonschema:"文件所属用户的ID"`
	V    *string `json:"v,omitempty" jsonschema:"存储空间标识"`
}

func GetFileLs(ctx context.Context, req *mcp.CallToolRequest, input *FileLsInput) (*mcp.CallToolResult, *wsclient.GetFileLsRetDto, error) {
	pass, err := Check(req, true)
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
	wsClient, _ = ConnectionManager.GetConnection(identity)
	var path *string = nil
	if input != nil {
		// 检查是否有任意字段被提供
		hasAny := input.Path != nil || input.Uid != nil || input.V != nil

		if hasAny {
			// 如果提供了任意一个，就必须全部提供，否则报错
			if input.Path == nil || input.Uid == nil || input.V == nil {
				_err := fmt.Errorf("参数验证失败：Path, Uid, V 必须同时提供，不允许部分为空（当前 Path=%v, Uid=%v, V=%v）",
					input.Path != nil, input.Uid != nil, input.V != nil)
				return &mcp.CallToolResult{
						Content: []mcp.Content{
							&mcp.TextContent{
								Text: "获取失败，" + _err.Error(),
							},
						},
					}, nil, fmt.Errorf("参数验证失败：Path, Uid, V 必须同时提供，不允许部分为空（当前 Path=%v, Uid=%v, V=%v）",
						input.Path != nil, input.Uid != nil, input.V != nil)
			}

			// 解引用指针获取实际值
			_path := fmt.Sprintf("vol%s/%s/%s", *input.V, *input.Uid, *input.Path)
			path = &_path
		}
	}

	result := wsClient.GetFileLs(path)
	return nil, &result, nil
}
