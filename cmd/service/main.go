package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	grpcServer := service.NewStatsMonitoringSever(ctx)
	defer grpcServer.Stop()

	go func() {
		err := grpcServer.Start()
		if err != nil {
			logger.Error(err.Error())
			cancel()
			return
		}
	}()

	<-ctx.Done()
}
