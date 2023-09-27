//go:build linux

package cpu

import (
	"strconv"
	"strings"

	"github.com/grahovsky/system-stats-daemon/internal/executor"
	"github.com/grahovsky/system-stats-daemon/internal/models"
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

	user, err := strconv.ParseFloat(fields[userPos], 64)
	if err != nil {
		return nil, err
	}

	system, err := strconv.ParseFloat(fields[systemPos], 64)
	if err != nil {
		return nil, err
	}

	idle, err := strconv.ParseFloat(fields[idlePos], 64)
	if err != nil {
		return nil, err
	}

	return &models.CPUInfo{
		User:   user,
		System: system,
		Idle:   idle,
	}, nil
}
