package tcp

import (
	"context"
	"go_redis/interface/tcp"
	"go_redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	address string
}

func ListenAndServeWithSignal(
	cfg *Config,
	handler tcp.Handler) error {

	closeChan := make(chan struct{}) // 监测系统是否将程序中断
	// 监测系统关闭信号
	sigCh := make(chan os.Signal)
	// 监听指定信号，如果监听到信号将信号发送到chan中
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		// 挂起，退出，终止， 中断信号
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", cfg.address)
	if err != nil {
		return err
	}
	logger.Info("start   listen")
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(
	listener net.Listener,
	handler tcp.Handler,
	closeChan <-chan struct{}) {

	go func() {
		// 监听程序是否被关闭如果关闭停止所有的通信
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()

	// 关闭处理
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	// 服务器错误or退出的时候需要将还在处理的handler执行完毕
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted  link")
		// 等待队列中增加一个在处理的协程
		waitDone.Add(1)
		go func() {
			defer func() {
				// 如果handler出现错误也能执行
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
