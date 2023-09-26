package client

import (
	"fmt"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
)

func PrintResponse(data *pb.StatsResponse) {
	fmt.Println("\nLoad Avg:")
	fmt.Println("\tload avg 1min:", data.GetLoadInfo().GetLoad_1Min())
	fmt.Println("\tload avg 5min:", data.GetLoadInfo().GetLoad_5Min())
	fmt.Println("\tload avg 15min:", data.GetLoadInfo().GetLoad_15Min())

	fmt.Println("CPU:")
	fmt.Println("\tuser:", data.GetCpuInfo().GetUser())
	fmt.Println("\tsystem:", data.GetCpuInfo().GetSystem())
	fmt.Println("\tidle:", data.GetCpuInfo().GetIdle())

	fmt.Println("Disk:")
	fmt.Println("\tkbt:", data.GetDiskInfo().GetKbt())
	fmt.Println("\ttps:", data.GetDiskInfo().GetTps())
}
