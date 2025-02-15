package tcp

import (
	"bufio"
	"context"
	"godis/lib/logger"
	"godis/lib/sync/atomic"
	"godis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

type EchoHandler struct {
	activeConn sync.Map //保存client的集合
	closing    atomic.Boolean
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close 关闭客户端连接
func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	_ = c.Conn.Close()
	return nil
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// 关闭中handler 拒绝处理新连接
		_ = conn.Close()
		return
	}

	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		// 发送数据前设置为waiting状态，阻止连接被关闭
		client.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

// Close 关闭服务器
func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
