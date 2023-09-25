package load

import (
	"strconv"
	"strings"

	"github.com/grahovsky/system-stats-daemon/internal/executor"
	"github.com/grahovsky/system-stats-daemon/internal/models"
)

const (
	load1MinPos  = 0
	load5MinPos  = 1
	load15MinPos = 2
)

// GetStats ...
func GetStats() (*models.LoadAverageInfo, error) {
	res, err := executor.Exec("cat", []string{"/proc/loadavg"})
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(res))

	load1Min, err := strconv.ParseFloat(fields[load1MinPos], 64)
	if err != nil {
		return nil, err
	}

	load5Min, err := strconv.ParseFloat(fields[load5MinPos], 64)
	if err != nil {
		return nil, err
	}

	load15Min, err := strconv.ParseFloat(fields[load15MinPos], 64)
	if err != nil {
		return nil, err
	}

	return &models.LoadAverageInfo{
		Load1Min:  load1Min,
		Load5Min:  load5Min,
		Load15Min: load15Min,
	}, nil
}
