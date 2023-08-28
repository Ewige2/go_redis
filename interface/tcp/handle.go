package tcp

import (
	"context"
	"net"
)

/*
* @desc: 处理TCP的连接
* @Author: Ewige2
* @Date: 2023/8/28 10:16
 */

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
