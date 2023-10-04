//go:build windows

package cpu

import (
	"github.com/StackExchange/wmi"
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

type Win32PerfFormattedDataPerfOSProcessor struct {
	PercentProcessorTime  uint64
	PercentUserTime       uint64
	PercentPrivilegedTime uint64
	PercentIdleTime       uint64
}

func GetStatsOs() (*models.CPUInfo, error) {
	var dst []Win32PerfFormattedDataPerfOSProcessor
	query := "SELECT * FROM Win32_PerfFormattedData_PerfOS_Processor WHERE Name = '_Total'"
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, err
	}

	cpuLoad := dst[0]
	return &models.CPUInfo{
		System: float64(cpuLoad.PercentProcessorTime),
		User:   float64(cpuLoad.PercentUserTime),
		Idle:   float64(cpuLoad.PercentIdleTime),
	}, nil
}
