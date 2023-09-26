package service

import (
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/monitor"
	"github.com/grahovsky/system-stats-daemon/internal/storage"
	mstorage "github.com/grahovsky/system-stats-daemon/internal/storage/memory"
)

type cStorage struct {
	msl storage.Storage
	msc storage.Storage
	msd storage.Storage
}

func (s *StatsMonitoringSever) StartMonitoring() {
	msl := mstorage.New()
	msc := mstorage.New()
	msd := mstorage.New()

	s.cStorage = &cStorage{
		msl: msl,
		msc: msc,
		msd: msd,
	}

	go func() {
		tiker := time.NewTicker(1 * time.Second)
		defer tiker.Stop()

		for {
			select {
			case <-tiker.C:
				monitor.ScanLoad(s.cStorage.msl)
				monitor.ScanCpu(s.cStorage.msc)
				monitor.ScanDisk(s.cStorage.msd)
			case <-s.ctx.Done():
				logger.Info("scan stoped..")
				return
			}
		}
	}()

	// // for debug
	// go func() {
	// 	tiker := time.NewTicker(5 * time.Second)
	// 	defer tiker.Stop()

	// 	for {
	// 		select {
	// 		case <-tiker.C:
	// 			elems := s.cStorage.msl.GetElements(5)
	// 			for e := range elems {
	// 				fmt.Printf("%+v\n", e)
	// 			}
	// 		case <-s.ctx.Done():
	// 			return
	// 		}
	// 	}
	// }()
}
