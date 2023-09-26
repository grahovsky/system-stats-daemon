//go:build windows
// +build windows

package disk

import "github.com/grahovsky/system-stats-daemon/internal/models"

func GetStatsOs() (*models.DiskInfo, error) {
	// no implemented
	return &models.DiskInfo{}, nil
}
