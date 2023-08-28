package system

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

// returns total, free
func GetRAMUsage() (uint64, uint64, error) {
	result, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, fmt.Errorf("could not get RAM usage: %s", err)
	}
	return result.Total, result.Available, nil
}
