package system

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// getCPUUsage returns the current cpu usage in a percentage
func GetCPUUsage() (float64, error) {
	result, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		return 0, fmt.Errorf("could not get CPU usage: %s", err)
	}

	return result[0], nil
}
