//go:build windows

package cpu

import (
	"strings"

	"github.com/grahovsky/system-stats-daemon/internal/executor"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/models"
	"github.com/grahovsky/system-stats-daemon/internal/stats"
)

func init() {
	executor.Exec("chcp", []string{"65001"})
	logger.Info("encoding set: UTF-8")
}

// криво работает, т.к. опрос идет раз в одну секунду, т.е. получается 1с + 1с задержка. Перевел на wmi.
func GetStatsOsBak() (*models.CPUInfo, error) {
	out, err := executor.Exec("typeperf", []string{
		`\Processor Information(_Total)\% Privileged Time`,
		`\Processor Information(_Total)\% User Time`,
		`\Processor Information(_Total)\% Idle Time`,
		"-sc",
		"1",
	})
	if err != nil {
		return nil, err
	}
	return parseData(out)
}

func parseData(output string) (*models.CPUInfo, error) {
	var err error
	lines := strings.Split(output, "\r\n")
	if len(lines) < 3 {
		return nil, stats.ErrInvalidData
	}

	fields := strings.Split(lines[2], ",")

	system := stats.SafeParseFloat(strings.Trim(fields[1], "\""))
	user := stats.SafeParseFloat(strings.Trim(fields[2], "\""))
	idle := stats.SafeParseFloat(strings.Trim(fields[3], "\""))
	if err != nil {
		return nil, err
	}
	cpuI := models.CPUInfo{
		System: system,
		User:   user,
		Idle:   idle,
	}
	return &cpuI, nil
}
