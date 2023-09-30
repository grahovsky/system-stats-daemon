package monitor

import (
	"context"
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/config"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/stats/cpu"
	"github.com/grahovsky/system-stats-daemon/internal/stats/disk"
	"github.com/grahovsky/system-stats-daemon/internal/stats/load"
	"github.com/grahovsky/system-stats-daemon/internal/storage"
	mstorage "github.com/grahovsky/system-stats-daemon/internal/storage/memory"
)

type Server struct {
	cStorage *cStorage
	ctx      context.Context
}

type cStorage struct {
	msdef storage.Storage
	msl   storage.Storage
	msc   storage.Storage
	msd   storage.Storage
}

func NewMonitor(ctx context.Context) *Server {
	return &Server{ctx: ctx}
}

func (s *Server) StartMonitoring() {
	msdef := mstorage.New()
	msl := mstorage.New()
	msc := mstorage.New()
	msd := mstorage.New()

	s.cStorage = &cStorage{
		msdef: msdef,
		msl:   msl,
		msc:   msc,
		msd:   msd,
	}

	go func() {
		tiker := time.NewTicker(1 * time.Second)
		defer tiker.Stop()

		for {
			select {
			case <-tiker.C:
				Default(s.cStorage.msdef)
				if config.Settings.Stats.Collect.CPU {
					ScanCPU(s.cStorage.msc)
				}
				if config.Settings.Stats.Collect.LoadAverage {
					ScanLoad(s.cStorage.msl)
				}
				if config.Settings.Stats.Collect.DiskInfo {
					ScanDisk(s.cStorage.msd)
				}
			case <-s.ctx.Done():
				logger.Info("stopped stats scan..")
				return
			}
		}
	}()
	logger.Info("started stats scan")
}

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
