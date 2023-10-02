//go:build linux

package talkers

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/storage"
)

func GetStatsOs(ctx context.Context, st storage.Storage) {
	cmd := exec.Command("iftop")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("iftop start error")
		return
	}

	if err := cmd.Start(); err != nil {
		logger.Error("iftop start error")
		return
	}

	tiker := time.NewTicker(1 * time.Second)
	defer tiker.Stop()
	scanner := bufio.NewScanner(stdout)

	for {
		select {
		case <-tiker.C:
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				logger.Error(err.Error())
			}
		case <-ctx.Done():
			logger.Info("stopped talkers scan..")
			return
		}
	}
}
