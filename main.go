package main

import (
	sysinfo "fn-mcp-server/Tools/SysInfo"
	"log"
	"net/http"

	"fn-mcp-server/Tools"
	"fn-mcp-server/wsclient"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func makeServerForRequest(r *http.Request) *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "fnnas-mcp-server",
			Version: "1.0.0",
		}, nil,
	)
	if r.Header.Get("Fn-Os-Url") != "" {
		mcp.AddTool(
			server,
			&mcp.Tool{
				Name:        "appcgi.sysinfo.getHostName",
				Description: "获取飞牛设备的HostName和版本信息",
			},
			tools.GetHostName,
		)
		mcp.AddTool(
			server,
			&mcp.Tool{
				Name:        "appcgi.sysinfo.getHardwareInfo",
				Description: "获取硬件信息,包含CPU/内存/虚拟化/BIOS/磁盘等详细信息",
			},
			sysinfo.GetHardwareInfo,
		)
		if r.Header.Get("Token") != "" {
			mcp.AddTool(
				server,
				&mcp.Tool{
					Name:        "file.ls",
					Description: "获取文件列表,不传递参数默认获取主目录数据",
				},
				tools.GetFileLs,
			)
			mcp.AddTool(
				server,
				&mcp.Tool{
					Name:        "appcgi.sysinfo.getMachineId",
					Description: "获取设备ID",
				},
				sysinfo.GetMachineId,
			)
			mcp.AddTool(
				server,
				&mcp.Tool{
					Name:        "appcgi.network.net.list",
					Description: "获取网络硬件信息,系统上所有网络接口的列表,包含:物理网卡/虚拟网卡/回环接口/Docker网桥等,排序:通常按Index或名称排序",
				},
				sysinfo.GetNetworkNetList,
			)
		}

	}

	return server
}

// 全局连接管理器
var connectionManager *wsclient.ConnectionManager

func main() {
	// 初始化连接管理器
	connectionManager = wsclient.NewConnectionManager()
	// 确保在程序结束时停止连接管理器，清理资源
	defer connectionManager.Stop()

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return makeServerForRequest(r)
	}, nil)

	http.Handle("/mcp", handler)

	log.Println("MCP HTTP server running at :8080/mcp")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
