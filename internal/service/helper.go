package service

import (
	"time"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

func (s *StatsMonitoringSever) LoadInfoAvg(at int64) *pb.LoadInfo {
	timeAt := time.Now().Add(-time.Duration(at) * time.Second)

	var e1, e5, e15 float64
	i := 0.0

	elems := s.cStorage.msl.GetElementsAt(timeAt)
	for el := range elems {
		elem := el.(*models.LoadAverageInfo)
		e1 += elem.Load1Min
		e5 += elem.Load5Min
		e15 += elem.Load15Min
		i++
	}

	if i == 0 {
		return &pb.LoadInfo{}
	}

	return &pb.LoadInfo{Load_1Min: e1 / i, Load_5Min: e5 / i, Load_15Min: e15 / i}
}

func (s *StatsMonitoringSever) CpuInfoAvg(at int64) *pb.CPUInfo {
	timeAt := time.Now().Add(-time.Duration(at) * time.Second)

	var user, system, idle float64
	i := 0.0

	elems := s.cStorage.msc.GetElementsAt(timeAt)
	for el := range elems {
		elem := el.(*models.CpuInfo)
		user += elem.User
		system += elem.System
		idle += elem.Idle
		i++
	}

	if i == 0 {
		return &pb.CPUInfo{}
	}

	return &pb.CPUInfo{User: user / i, System: system / i, Idle: idle / i}
}

func (s *StatsMonitoringSever) DiskInfoAvg(at int64) *pb.DiskInfo {
	timeAt := time.Now().Add(-time.Duration(at) * time.Second)

	var kbt, tps float64
	i := 0.0

	elems := s.cStorage.msd.GetElementsAt(timeAt)
	for el := range elems {
		elem := el.(*models.DiskInfo)
		kbt += elem.Kbt
		tps += elem.Tps
		i++
	}

	if i == 0 {
		return &pb.DiskInfo{}
	}

	return &pb.DiskInfo{Kbt: kbt / i, Tps: tps / i}
}

func (s *StatsMonitoringSever) checkRT(t_at int64) bool {
	// return time.Now().After(s.m_at.Add(time.Duration(t_at) * time.Second))
	// if limit storage, date stats != date server start, add technical storage
	sa := <-s.cStorage.msdef.StoreAt()
	logger.Info(sa.(time.Time).String())
	return time.Now().After(sa.(time.Time).Add(time.Duration(t_at) * time.Second))
}
