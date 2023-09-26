package load

import (
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

func GetStats() (*models.LoadAverageInfo, error) {
	cpuInfo, err := GetStatsOs()

	return cpuInfo, err
}
