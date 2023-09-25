package cpu

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/grahovsky/system-stats-daemon/internal/models"
)

func GetStats() (*models.CpuInfo, error) {
	cpuIdle, cpuTotal := getCPUSample()

	return &models.CpuInfo{
		Idle:  cpuIdle,
		Total: cpuTotal,
	}, nil
}

func getCPUSample() (idle, total uint64) {
	contents, err := os.ReadFile("/proc/stat")
	if err != nil {
		return
	}
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
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}
