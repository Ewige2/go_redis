package main

import (
	"context"
	"go_redis/interface/tcp"
	"go_redis/lib/logger"
	"net"
)

type Config struct {
	address string
}

func ListenAndServeWithSognal(
	cfg *Config,
	handler tcp.Handler) error {
	closeChan := make(chan struct{}) //  实现优雅关闭
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
	ctx := context.Background()
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted  link")
		go func() {
			handler.Handle(ctx, conn)
		}()
	}

}
