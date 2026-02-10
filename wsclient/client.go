package wsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// FnosURL 全局FNOS服务器URL，实际项目中应该从环境变量或配置文件中读取
var FnosURL = "http://fnos.xn--1jqw64a7tu.cn:25130/"

// FnOsWsBase WebSocket客户端核心结构体
type FnOsWsBase struct {
	LoginRetDto
	FnosUrl        string
	conn           *websocket.Conn
	IsRun          bool
	Mu             sync.Mutex             // 保护pendingFutures、isRun和conn
	ConnMu         sync.Mutex             // 保护WebSocket连接的并发写入
	PendingFutures map[string]chan []byte // reqid -> 结果通道
	Done           chan struct{}          // 用于停止goroutine

	Iv      []byte
	Key     string
	IsLogin bool
}

// NewFnOsWsBase 创建WebSocket客户端实例
func NewFnOsWsBase(fnosUrl string, token *string) *FnOsWsBase {

	l := LoginRetDto{}
	if token != nil {
		l = LoginRetDto{
			Token: *token,
		}
	}

	return &FnOsWsBase{
		FnosUrl:        fnosUrl,
		IsRun:          false,
		PendingFutures: make(map[string]chan []byte),
		Done:           make(chan struct{}),
		Iv:             generateIv(),
		Key:            generateKey(),
		LoginRetDto:    l,
	}
}

// connect 建立WebSocket连接并启动处理循环
func (f *FnOsWsBase) connect(wsType string) error {
	// 解析URL并确定ws/wss协议
	parsedUrl, err := url.Parse(f.FnosUrl)
	if err != nil {
		return fmt.Errorf("解析URL失败: %w", err)
	}

	protocol := "ws"
	if parsedUrl.Scheme == "https" {
		protocol = "wss"
	}

	// 构建WebSocket连接地址
	wsUrl := fmt.Sprintf("%s://%s/websocket?type=%s", protocol, parsedUrl.Host, wsType)

	headers := http.Header{}
	// 建立连接
	if f.Token != "" {
		headers.Add("fnos-token", f.Token)
	}
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, headers)
	if err != nil {
		return fmt.Errorf("连接WebSocket失败: %w", err)
	}

	f.Mu.Lock()
	f.conn = conn
	f.IsRun = true
	f.Mu.Unlock()
	log.Info().Msg("WS连接已建立")

	// 启动读消息goroutine
	go f.readMessageLoop()

	// 启动心跳和活跃检测
	go f.sendHeartbeat()
	go f.sendActive()

	return nil
}

// readMessageLoop 循环读取WebSocket消息
func (f *FnOsWsBase) readMessageLoop() {
	defer func() {
		f.Mu.Lock()
		f.IsRun = false
		f.Mu.Unlock()
		_ = f.conn.Close()
		log.Error().Msg("消息读取循环退出")
	}()

	for {
		select {
		case <-f.Done:
			return
		default:
			// 读取消息
			_, msgBytes, err := f.conn.ReadMessage()
			if err != nil {
				log.Error().Err(err).Msg("读取消息失败")
				return
			}

			log.Info().Str("message", string(msgBytes)).Msg("接收到消息")

			// 解析JSON消息
			var msg map[string]interface{}
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				log.Error().Err(err).Msg("解析消息JSON失败")
				continue
			}

			// 获取reqid并匹配pending请求
			reqID, ok := msg["reqid"].(string)
			if !ok || reqID == "" {
				continue
			}

			f.Mu.Lock()
			ch, exists := f.PendingFutures[reqID]
			if exists {
				// 发送结果到通道并清理
				ch <- msgBytes
				close(ch)
				delete(f.PendingFutures, reqID)
			}
			f.Mu.Unlock()
		}
	}
}

// Send 发送消息并等待返回（带超时）
func (f *FnOsWsBase) Send(msg interface{}, reqID string, timeout time.Duration) (map[string]interface{}, error) {
	// 等待连接建立（最多5秒）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("等待连接建立超时")
		default:
			f.Mu.Lock()
			runStatus := f.IsRun
			f.Mu.Unlock()
			if runStatus {
				goto connected
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

connected:
	var msgBytes []byte
	var err error

	if reflect.TypeOf(msg).Kind() == reflect.String {
		msgBytes = []byte(msg.(string))
	} else {
		// 将消息转为JSON
		msgBytes, err = json.Marshal(msg)
		if err != nil {
			return nil, fmt.Errorf("序列化消息失败: %w", err)
		}
	}

	// 创建结果通道并加入pending
	ch := make(chan []byte, 1)
	f.Mu.Lock()
	f.PendingFutures[reqID] = ch
	f.Mu.Unlock()

	// 发送消息
	f.ConnMu.Lock()
	err = f.conn.WriteMessage(websocket.TextMessage, msgBytes)
	f.ConnMu.Unlock()
	if err != nil {
		f.Mu.Lock()
		delete(f.PendingFutures, reqID)
		f.Mu.Unlock()
		close(ch)
		return nil, fmt.Errorf("发送消息失败: %w", err)
	}

	// 等待返回（带超时）
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		f.Mu.Lock()
		delete(f.PendingFutures, reqID)
		f.Mu.Unlock()
		close(ch)
		log.Error().Msg("获取返回值超时")
		return nil, ctx.Err()
	case resBytes := <-ch:
		var res map[string]interface{}
		if err := json.Unmarshal(resBytes, &res); err != nil {
			return nil, fmt.Errorf("解析返回消息失败: %w", err)
		}
		return res, nil
	}
}

// sendHeartbeat 发送心跳包（每10秒）
func (f *FnOsWsBase) sendHeartbeat() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-f.Done:
			return
		case <-ticker.C:
			f.Mu.Lock()
			conn := f.conn
			runStatus := f.IsRun
			f.Mu.Unlock()
			if conn == nil || !runStatus {
				return
			}

			heartbeatMsg := map[string]interface{}{
				"req": "ping",
			}
			f.ConnMu.Lock()
			err := conn.WriteJSON(heartbeatMsg)
			f.ConnMu.Unlock()
			if err != nil {
				log.Warn().Err(err).Msg("发送心跳包失败")
			} else {
				log.Debug().Msg("发送心跳包")
			}
		}
	}
}

// sendActive 发送活跃检测（每60秒）
func (f *FnOsWsBase) sendActive() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-f.Done:
			return
		case <-ticker.C:
			f.Mu.Lock()
			conn := f.conn
			runStatus := f.IsRun
			f.Mu.Unlock()
			if conn == nil || !runStatus {
				return
			}

			activeMsg := map[string]interface{}{
				"req":   "user.active",
				"reqid": f.GetReqId(),
			}
			f.ConnMu.Lock()
			err := conn.WriteJSON(activeMsg)
			f.ConnMu.Unlock()
			if err != nil {
				log.Warn().Err(err).Msg("发送user.active失败")
			} else {
				log.Debug().Msg("发送user.active")
			}
		}
	}
}

// Start 启动WebSocket客户端
func (f *FnOsWsBase) Start(wsType string) error {
	if err := f.connect(wsType); err != nil {
		return err
	}
	return nil
}

// Stop 停止WebSocket客户端
func (f *FnOsWsBase) Stop() {
	close(f.Done)
	f.Mu.Lock()
	f.IsRun = false
	f.Mu.Unlock()
	if f.conn != nil {
		_ = f.conn.Close()
	}
	log.Info().Msg("WS客户端已停止")
}
