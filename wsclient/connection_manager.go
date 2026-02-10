package wsclient

import (
	"sync"
	"time"
)

// ConnectionManager 连接管理器，统一管理不同身份的WebSocket连接
type ConnectionManager struct {
	mu             sync.RWMutex
	connections    map[string]*FnOsWsBase // key: 身份标识（如用户名、token等）
	lastAccessTime map[string]time.Time   // key: 身份标识，value: 最后访问时间
	cleanupTicker  *time.Ticker           // 定期清理定时器
	done           chan struct{}          // 用于停止清理goroutine
}

// NewConnectionManager 创建连接管理器实例
func NewConnectionManager() *ConnectionManager {
	cm := &ConnectionManager{
		connections:    make(map[string]*FnOsWsBase),
		lastAccessTime: make(map[string]time.Time),
		cleanupTicker:  time.NewTicker(1 * time.Minute), // 每分钟检查一次
		done:           make(chan struct{}),
	}

	// 启动清理goroutine
	go cm.cleanupExpiredConnections()

	return cm
}

// GetConnection 获取指定身份的连接，并更新最后访问时间
func (cm *ConnectionManager) GetConnection(identity string) (*FnOsWsBase, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, exists := cm.connections[identity]
	if exists {
		// 更新最后访问时间
		cm.lastAccessTime[identity] = time.Now()
	}

	return conn, exists
}

// AddConnection 添加新连接，并设置初始最后访问时间
func (cm *ConnectionManager) AddConnection(identity string, conn *FnOsWsBase) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 添加连接
	cm.connections[identity] = conn
	// 设置初始最后访问时间
	cm.lastAccessTime[identity] = time.Now()
}

// RemoveConnection 移除连接
func (cm *ConnectionManager) RemoveConnection(identity string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if conn, exists := cm.connections[identity]; exists {
		conn.Stop()
		delete(cm.connections, identity)
		delete(cm.lastAccessTime, identity)
	}
}

// GetAllConnections 获取所有连接
func (cm *ConnectionManager) GetAllConnections() map[string]*FnOsWsBase {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// 返回副本，避免外部修改
	copy := make(map[string]*FnOsWsBase, len(cm.connections))
	for k, v := range cm.connections {
		copy[k] = v
	}
	return copy
}

// CountConnections 获取连接数量
func (cm *ConnectionManager) CountConnections() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.connections)
}

// Stop 停止连接管理器，清理所有资源
func (cm *ConnectionManager) Stop() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 停止清理定时器
	close(cm.done)
	cm.cleanupTicker.Stop()

	// 清理所有连接
	for identity, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, identity)
		delete(cm.lastAccessTime, identity)
	}
}

// cleanupExpiredConnections 定期清理过期连接（超过5分钟未访问的连接）
func (cm *ConnectionManager) cleanupExpiredConnections() {
	for {
		select {
		case <-cm.done:
			return
		case <-cm.cleanupTicker.C:
			cm.cleanupConnections()
		}
	}
}

// cleanupConnections 清理过期连接的实际实现
func (cm *ConnectionManager) cleanupConnections() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	now := time.Now()
	expiration := 5 * time.Minute // 5分钟过期时间

	for identity, lastAccess := range cm.lastAccessTime {
		// 检查是否过期
		if now.Sub(lastAccess) > expiration {
			// 关闭并移除过期连接
			if conn, exists := cm.connections[identity]; exists {
				conn.Stop()
				delete(cm.connections, identity)
			}
			delete(cm.lastAccessTime, identity)
		}
	}
}
