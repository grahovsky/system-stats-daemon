package disk

import (
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

func GetStats() (*models.DiskInfo, error) {
	cpuInfo, err := GetStatsOs()

	return cpuInfo, err
}
