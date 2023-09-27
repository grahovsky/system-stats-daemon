package service

import (
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/monitor"
	"github.com/grahovsky/system-stats-daemon/internal/storage"
	mstorage "github.com/grahovsky/system-stats-daemon/internal/storage/memory"
)

type cStorage struct {
	msdef storage.Storage
	msl   storage.Storage
	msc   storage.Storage
	msd   storage.Storage
}

func (s *StatsMonitoringSever) StartMonitoring() {
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
	s.m_at = time.Now()

	go func() {
		tiker := time.NewTicker(1 * time.Second)
		defer tiker.Stop()

		for {
			select {
			case <-tiker.C:
				monitor.Default(s.cStorage.msdef)
				monitor.ScanLoad(s.cStorage.msl)
				monitor.ScanCpu(s.cStorage.msc)
				monitor.ScanDisk(s.cStorage.msd)
			case <-s.ctx.Done():
				logger.Info("scan stoped..")
				return
			}
		}
	}()
}
