//go:build darwin

package cpu

import "github.com/grahovsky/system-stats-daemon/internal/models"

func GetStatsOs() (*models.CPUInfo, error) {
	// no implemented
	return &models.CPUInfo{}, nil
}
