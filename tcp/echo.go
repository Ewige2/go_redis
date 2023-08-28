package tcp

import (
	"bufio"
	"context"
	"go_redis/lib/logger"
	"go_redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

/*
* @desc:  回复处理
* @Author: Ewige2
* @Date: 2023/8/28 19:05
 */

type EchoHandler struct {
	// 保存所有工作状态的client集合， 这里map的value为空结构体，将map当做set使用
	activeConn sync.Map
	// 关闭状态标识位
	closing atomic.Bool
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

// Close 关闭服务器
func (h *EchoHandler) Close() error {
	logger.Info("handler  shutting  down....")
	h.closing.Store(true)
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true

	})
	return nil

}

// 客户端连接的抽象
type EchoClient struct {
	Conn net.Conn
	// 当服务端开始发送数据时进入waiting ，阻止其他协程关闭连接
	Waiting wait.Wait
}

// Close 关闭连接
func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	c.Conn.Close()
	return nil
}

// 处理客户端连接
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Load() {
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	// 读取数据
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection  close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}
