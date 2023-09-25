package main

import (
	"context"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/config"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/monitor"
	memoryStorage "github.com/grahovsky/system-stats-daemon/internal/storage/memory"
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

	// ctx, done := context.WithCancel(ctx)
	go func() {
		// defer done()
		ms := memoryStorage.New()
		monitor.NewLoad(ctx, ms)
	}()

	go func() {
		// defer done()
		ms := memoryStorage.New()
		monitor.NewCpu(ctx, ms)
	}()

	// <-ctx.Done()

	time.Sleep(10 * time.Second)
}
