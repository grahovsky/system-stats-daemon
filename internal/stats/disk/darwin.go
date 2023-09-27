//go:build darwin

package disk

import "github.com/grahovsky/system-stats-daemon/internal/models"

func GetStatsOs() (*models.DiskInfo, error) {
	// no implemented
	return &models.DiskInfo{}, nil
}
