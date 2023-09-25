package monitor

import (
	"context"
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/stats/cpu"
	"github.com/grahovsky/system-stats-daemon/internal/stats/load"
	"github.com/grahovsky/system-stats-daemon/internal/storage"
)

func NewLoad(ctx context.Context, ms storage.Storage) {
	num := 0
	tiker := time.NewTicker(1 * time.Second)
	defer tiker.Stop()
	defer ms.Show()

	for {
		select {
		case <-tiker.C:
			d, err := load.GetStats()
			if err != nil {
				logger.Error(err.Error())
			}

			ms.Push(d, time.Now())
			num++
			if num >= 5 {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func NewCpu(ctx context.Context, ms storage.Storage) {
	num := 0
	tiker := time.NewTicker(1 * time.Second)
	defer tiker.Stop()
	defer ms.Show()

	for {
		select {
		case <-tiker.C:
			d, err := cpu.GetStats()
			if err != nil {
				logger.Error(err.Error())
			}
			ms.Push(d, time.Now())
			num++
			if num >= 5 {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
