package client

import (
	"fmt"
	"math"

	pb "github.com/grahovsky/system-stats-daemon/internal/api/stats_service"
)

const dec = 100

func PrintResponse(data *pb.StatsResponse) {
	fmt.Println("\nLoad Avg:")
	fmt.Println("\tload avg 1min:", round(data.GetLoadInfo().GetLoad_1Min()))
	fmt.Println("\tload avg 5min:", round(data.GetLoadInfo().GetLoad_5Min()))
	fmt.Println("\tload avg 15min:", round(data.GetLoadInfo().GetLoad_15Min()))

	fmt.Println("CPU:")
	fmt.Println("\tuser:", round(data.GetCPUInfo().GetUser()))
	fmt.Println("\tsystem:", round(data.GetCPUInfo().GetSystem()))
	fmt.Println("\tidle:", round(data.GetCPUInfo().GetIdle()))

	fmt.Println("Disk:")
	fmt.Println("\tkbt:", round(data.GetDiskInfo().GetKbt()))
	fmt.Println("\ttps:", round(data.GetDiskInfo().GetTps()))
}

func round(v float64) float64 {
	return math.Round(v*dec) / dec
}
