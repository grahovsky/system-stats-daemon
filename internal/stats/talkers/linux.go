//go:build linux

package talkers

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/grahovsky/system-stats-daemon/internal/models"
	"github.com/grahovsky/system-stats-daemon/internal/storage"
)

const (
	senderPos       = 1
	receiverPos     = 0
	senderRatePos   = 5
	receiverRatePos = 4
)

func GetStatsOs(ctx context.Context, st storage.Storage) {
	cmd := exec.Command("iftop", "-t", "-o 40")

	stdout, err := cmd.StdoutPipe()
	// stdout, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("iftop start error")
		return
	}

	if err := cmd.Start(); err != nil {
		logger.Error(fmt.Errorf("iftop start error. %w", err).Error())
		return
	}

	scanner := bufio.NewScanner(stdout)
	var prevLine string
	var currentLine string
	var talkers models.Talkers

	for scanner.Scan() {
		currentLine = scanner.Text()

		if strings.Contains(prevLine, " 1 ") && talkers.Top1 == "" {
			talkers.Top1 = getValue(prevLine, currentLine)
		}
		if strings.Contains(prevLine, " 2 ") && talkers.Top2 == "" {
			talkers.Top2 = getValue(prevLine, currentLine)
		}
		if strings.Contains(prevLine, " 3 ") && talkers.Top3 == "" {
			talkers.Top3 = getValue(prevLine, currentLine)
		}

		if talkers.Top3 != "" {
			pushed := talkers
			talkers = models.Talkers{}
			st.Push(&pushed, time.Now())
		}

		prevLine = currentLine

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func getValue(sender, receiver string) string {
	fieldsSender := strings.Fields(sender)
	fieldsReceiver := strings.Fields(receiver)

	if len(fieldsReceiver) < 7 && len(fieldsSender) < 6 {
		return ""
	}

	return fmt.Sprintf("%s=>%s %s + %s",
		fieldsSender[senderPos],
		fieldsReceiver[receiverPos],
		fieldsSender[senderRatePos],
		fieldsReceiver[receiverRatePos])
}
