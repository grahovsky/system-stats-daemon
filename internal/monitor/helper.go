package monitor

import (
	"time"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

func (s *Server) LoadInfoAvg(at int64) *pb.LoadInfo {
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

func (s *Server) CPUInfoAvg(at int64) *pb.CPUInfo {
	timeAt := time.Now().Add(-time.Duration(at) * time.Second)

	var user, system, idle float64
	i := 0.0

	elems := s.cStorage.msc.GetElementsAt(timeAt)
	for el := range elems {
		elem := el.(*models.CPUInfo)
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

func (s *Server) DiskInfoAvg(at int64) *pb.DiskInfo {
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

func (s *Server) TalkersAvg(at int64) *pb.TalkersInfo {
	timeAt := time.Now().Add(-time.Duration(at) * time.Second)

	top1 := ""
	top2 := ""
	top3 := ""

	elems := s.cStorage.mst.GetElementsAt(timeAt)
	for el := range elems {
		elem := el.(*models.Talkers)
		top1 = elem.Top1
		top2 = elem.Top2
		top3 = elem.Top3
	}

	// to do: impl correctly avg for talkers

	return &pb.TalkersInfo{Top1: top1, Top2: top2, Top3: top3}
}

func (s *Server) CheckRT(tAt int64) bool {
	sa := <-s.cStorage.msdef.StoreAt()
	if tt, ok := sa.(time.Time); ok {
		return time.Now().After(tt.Add(time.Duration(tAt) * time.Second))
	}
	return false
}
