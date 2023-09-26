//go:build windows
// +build windows

package cpu

import (
	"github.com/StackExchange/wmi"
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

type Win32_PerfFormattedData_PerfOS_Processor struct {
	PercentProcessorTime  uint64
	PercentUserTime       uint64
	PercentPrivilegedTime uint64
	PercentIdleTime       uint64
}

func getCPUSample() (*models.CpuInfo, error) {
	var dst []Win32_PerfFormattedData_PerfOS_Processor
	query := "SELECT * FROM Win32_PerfFormattedData_PerfOS_Processor WHERE Name = '_Total'"
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, err
	}

	cpuLoad := dst[0]
	return &models.CpuInfo{
		System: cpuLoad.PercentProcessorTime,
		User:   cpuLoad.PercentUserTime,
		Idle:   cpuLoad.PercentIdleTime,
	}, nil
}
