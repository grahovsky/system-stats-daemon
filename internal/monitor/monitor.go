package monitor

import (
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/stats/cpu"
	"github.com/grahovsky/system-stats-daemon/internal/stats/disk"
	"github.com/grahovsky/system-stats-daemon/internal/stats/load"
	"github.com/grahovsky/system-stats-daemon/internal/storage"
)

func Default(ms storage.Storage) {
	ms.Push(&struct{}{}, time.Now())
}

func ScanLoad(ms storage.Storage) {
	d, err := load.GetStats()
	if err != nil {
		logger.Error(err.Error())
	}
	ms.Push(d, time.Now())
}

func ScanCPU(ms storage.Storage) {
	d, err := cpu.GetStats()
	if err != nil {
		logger.Error(err.Error())
	}
	ms.Push(d, time.Now())
}

func ScanDisk(ms storage.Storage) {
	d, err := disk.GetStats()
	if err != nil {
		logger.Error(err.Error())
	}
	ms.Push(d, time.Now())
}
