//go:build windows
// +build windows

package cpu

import (
	"strconv"
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

// криво работает, т.к. опрос идет раз в одну секунду, т.е. получается 1с + 1с задержка. Перевел на wmi
func GetStatsOs_() (*models.CpuInfo, error) { //nolint
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
	return parseData(string(out))
}

func parseData(output string) (*models.CpuInfo, error) { //nolint
	var err error
	lines := strings.Split(output, "\r\n")
	if len(lines) < 3 {
		return nil, stats.ErrInvalidData
	}

	fields := strings.Split(lines[2], ",")

	system, err := strconv.ParseFloat(strings.Trim(fields[1], "\""), 64)
	if err != nil {
		return nil, err
	}

	user, err := strconv.ParseFloat(strings.Trim(fields[2], "\""), 64)
	if err != nil {
		return nil, err
	}

	idle, err := strconv.ParseFloat(strings.Trim(fields[3], "\""), 64)
	if err != nil {
		return nil, err
	}
	cpuI := models.CpuInfo{
		System: system,
		User:   user,
		Idle:   idle,
	}
	return &cpuI, nil
}
