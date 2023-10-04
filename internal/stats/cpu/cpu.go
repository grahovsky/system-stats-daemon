package cpu

import (
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

func GetStats() (*models.CPUInfo, error) {
	cpuInfo, err := GetStatsOs()

	return cpuInfo, err
}
