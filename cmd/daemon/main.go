package main

import (
	"strconv"

	"github.com/grahovsky/system-stats-daemon/internal/config"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
)

func main() {
	logger.Info(config.Settings.Server.Host)
	logger.Info(strconv.FormatBool(config.Settings.Metrics.Collect.Cpu))
}
