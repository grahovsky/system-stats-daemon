//go:build linux
// +build linux

package cpu

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/grahovsky/system-stats-daemon/internal/models"
)

func getCPUSample() (*models.CpuInfo, error) {
	contents, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	cpuI := models.CpuInfo{}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				cpuI.System += val // tally up all the numbers to get total ticks
				if i == 4 {        // idle is the 5th field in the cpu line
					cpuI.Idle = val
				}
			}
			return &cpuI, nil
		}
	}
	return &cpuI, nil
}