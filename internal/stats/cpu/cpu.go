package cpu

import (
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

func GetStats() (*models.CpuInfo, error) {
	cpuInfo, err := getCPUSample()

	return cpuInfo, err
}
