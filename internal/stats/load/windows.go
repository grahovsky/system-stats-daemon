//go:build windows
// +build windows

package load

import "github.com/grahovsky/system-stats-daemon/internal/models"

func GetStatsOs() (*models.LoadAverageInfo, error) {
	// no implemented
	return &models.LoadAverageInfo{}, nil
}
