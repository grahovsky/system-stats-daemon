package main

import (
	"context"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/grahovsky/system-stats-daemon/internal/config"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/service"
)

func main() {
	logger.Info(config.Settings.Server.Host)
	logger.Info(strconv.FormatBool(config.Settings.Stats.Collect.Cpu))

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
