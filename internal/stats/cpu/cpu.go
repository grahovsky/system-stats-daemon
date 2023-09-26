package cpu

import (
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

func GetStats() (*models.CpuInfo, error) {
	cpuInfo, err := GetStatsOs()

	return cpuInfo, err
}
