//go:build linux
// +build linux

package disk

import (
	"strconv"
	"strings"

	"github.com/grahovsky/system-stats-daemon/internal/executor"
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

const (
	kbtPos = 29
	tpsPos = 30
)

func GetStatsOs() (*models.DiskInfo, error) {
	res, err := executor.Exec("iostat", []string{})
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(res)

	kbt, err := strconv.ParseFloat(fields[kbtPos], 64)
	if err != nil {
		return nil, err
	}

	tps, err := strconv.ParseFloat(fields[tpsPos], 64)
	if err != nil {
		return nil, err
	}

	return &models.DiskInfo{
		Kbt: kbt,
		Tps: tps,
	}, nil
}
