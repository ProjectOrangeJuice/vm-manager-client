package system

import (
	"container-manager/shared"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func GetFreeStorageSpace() ([]shared.StorageResult, error) {
	// Get the output of the `df` command.
	cmd := exec.Command("df")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not execute the command, %s", err)
	}

	// Split the output into lines.
	lines := strings.Split(string(output), "\n")
	storageList := make([]shared.StorageResult, 0)

	// Iterate over the lines and print the storage free information.
	for index, line := range lines {
		// Split the line into fields.
		fields := strings.Fields(line)
		if len(fields) < 5 || index == 0 {
			continue // This line can't be read
		}

		size := fields[1][:len(fields[1])-1]
		size_float, err := strconv.ParseFloat(size, 64)
		if err != nil {
			log.Printf("Could not convert size to float, %s", err)
			continue
		}

		used := fields[2][:len(fields[2])-1]
		used_float, err := strconv.ParseFloat(used, 64)
		if err != nil {
			log.Printf("Could not convert used size to float, %s", err)
			continue
		}

		storageList = append(storageList, shared.StorageResult{
			Name:       fields[0],
			TotalSpace: size_float,
			UsedSpace:  used_float,
			Mount:      fields[5],
		})
	}

	return storageList, nil
}
