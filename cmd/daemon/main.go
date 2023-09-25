package main

import (
	"context"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/grahovsky/system-stats-daemon/internal/config"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
)

func main() {
	logger.Info(config.Settings.Server.Host)
	logger.Info(strconv.FormatBool(config.Settings.Metrics.Collect.Cpu))

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	<-ctx.Done()
}
