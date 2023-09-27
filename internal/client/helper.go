package client

import (
	"fmt"
	"math"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
)

const dec = 100

func PrintResponse(data *pb.StatsResponse) {
	if !EmptyLoad(data.GetLoadInfo()) {
		fmt.Println("\nLoad Avg:")
		fmt.Println("\tload avg 1min:", round(data.GetLoadInfo().GetLoad_1Min()))
		fmt.Println("\tload avg 5min:", round(data.GetLoadInfo().GetLoad_5Min()))
		fmt.Println("\tload avg 15min:", round(data.GetLoadInfo().GetLoad_15Min()))
	}

	if !EmptyCPU(data.GetCpuInfo()) {
		fmt.Println("CPU:")
		fmt.Println("\tuser:", round(data.GetCpuInfo().GetUser()))
		fmt.Println("\tsystem:", round(data.GetCpuInfo().GetSystem()))
		fmt.Println("\tidle:", round(data.GetCpuInfo().GetIdle()))
	}

	if !EmptyDisk(data.GetDiskInfo()) {
		fmt.Println("Disk:")
		fmt.Println("\tkbt:", round(data.GetDiskInfo().GetKbt()))
		fmt.Println("\ttps:", round(data.GetDiskInfo().GetTps()))
	}
}

func round(v float64) float64 {
	return math.Round(v*dec) / dec
}

func EmptyLoad(load *pb.LoadInfo) bool {
	return load.Load_1Min == 0.0 && load.Load_5Min == 0.0 && load.Load_15Min == 0.0
}

func EmptyCPU(cpu *pb.CPUInfo) bool {
	return cpu.User == 0.0 && cpu.System == 0.0 && cpu.Idle == 0.0
}

func EmptyDisk(disk *pb.DiskInfo) bool {
	return disk.Kbt == 0.0 && disk.Tps == 0.0
}
