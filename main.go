package main

import (
	"log"
	"net/http"

	"fn-mcp-server/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CityArgs struct {
	City string `json:"city"`
}

var tzMap = map[string]string{
	"nyc":    "America/New_York",
	"sf":     "America/Los_Angeles",
	"boston": "America/New_York",
}

func main() {

	// 创建 MCP Server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "fnnas-mcp-server",
			Version: "1.0.0",
		},
		nil,
	)

	// 注册 tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "appcgi.sysinfo.getHostName",
		Description: "获取飞牛设备的HostName和版本信息",
	}, tools.GetHostName)

	// 必须提供 options（v1.2.0 要求）
	handler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return server
		},
		&mcp.StreamableHTTPOptions{},
	)

	http.Handle("/mcp", handler)

	log.Println("MCP HTTP server running at :8080/mcp")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
