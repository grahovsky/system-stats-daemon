package client

import (
	"fmt"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
)

func PrintResponse(data *pb.StatsResponse) {
	fmt.Println("\nCPU:")
	fmt.Println("\tuser mode:", data.GetCpuInfo().GetUser())
	fmt.Println("\tsystem mode:", data.GetCpuInfo().GetSystem())
	fmt.Println("\tidle time:", data.GetCpuInfo().GetIdle())

	fmt.Println("Disk:")
	fmt.Println("\tkbt:", data.GetDiskInfo().GetKbt())
	fmt.Println("\ttps:", data.GetDiskInfo().GetTps())

	fmt.Println("Avg load:")
	fmt.Println("\tavg load 1min:", data.GetLoadInfo().GetLoad_1Min())
	fmt.Println("\tavg load 5min:", data.GetLoadInfo().GetLoad_5Min())
	fmt.Println("\tavg load 15min:", data.GetLoadInfo().GetLoad_15Min())
}
