package main

import (
	"context"
	"fmt"
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
	msl := memoryStorage.New()
	go func() {
		monitor.NewLoad(ctx, msl)
	}()

	msc := memoryStorage.New()
	go func() {
		monitor.NewCpu(ctx, msc)
	}()

	tiker := time.NewTicker(5 * time.Second)
	defer tiker.Stop()

	for {
		select {
		case <-tiker.C:
			elems := msc.GetElements(5)
			for e := range elems {
				fmt.Printf("%+v\n", e)
			}
		case <-ctx.Done():
			return
		}
	}
}
