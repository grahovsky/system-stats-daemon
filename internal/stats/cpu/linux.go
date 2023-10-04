//go:build linux

package cpu

import (
	"strings"

	"github.com/grahovsky/system-stats-daemon/internal/executor"
	"github.com/grahovsky/system-stats-daemon/internal/models"
	"github.com/grahovsky/system-stats-daemon/internal/stats"
)

const (
	userPos   = 14
	systemPos = 16
	idlePos   = 19
)

func GetStatsOs() (*models.CPUInfo, error) {
	res, err := executor.Exec("iostat", []string{"-c"})
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(res)

	user := stats.SafeParseFloat(fields[userPos])
	system := stats.SafeParseFloat(fields[systemPos])
	idle := stats.SafeParseFloat(fields[idlePos])

	return &models.CPUInfo{
		User:   user,
		System: system,
		Idle:   idle,
	}, nil
}
