//go:build linux

package load

import (
	"strings"

	"github.com/grahovsky/system-stats-daemon/internal/executor"
	"github.com/grahovsky/system-stats-daemon/internal/models"
	"github.com/grahovsky/system-stats-daemon/internal/stats"
)

const (
	load1MinPos  = 0
	load5MinPos  = 1
	load15MinPos = 2
)

func GetStatsOs() (*models.LoadAverageInfo, error) {
	res, err := executor.Exec("cat", []string{"/proc/loadavg"})
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(res)

	load1Min := stats.SafeParseFloat(fields[load1MinPos])
	load5Min := stats.SafeParseFloat(fields[load5MinPos])
	load15Min := stats.SafeParseFloat(fields[load15MinPos])

	return &models.LoadAverageInfo{
		Load1Min:  load1Min,
		Load5Min:  load5Min,
		Load15Min: load15Min,
	}, nil
}
