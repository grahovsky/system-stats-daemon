//go:build linux

package disk

import (
	"strings"

	"github.com/grahovsky/system-stats-daemon/internal/executor"
	"github.com/grahovsky/system-stats-daemon/internal/models"
	"github.com/grahovsky/system-stats-daemon/internal/stats"
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

	kbt := stats.SafeParseFloat(fields[kbtPos])
	tps := stats.SafeParseFloat(fields[tpsPos])

	return &models.DiskInfo{
		Kbt: kbt,
		Tps: tps,
	}, nil
}
