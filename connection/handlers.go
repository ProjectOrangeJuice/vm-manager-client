package connection

import (
	"fmt"
	"log"

	"github.com/ProjectOrangeJuice/vm-manager-client/system"
	"github.com/ProjectOrangeJuice/vm-manager-server/shared"
)

func (c *Connection) sendBackStorage() {
	log.Print("Sending storage info")
	storages, err := system.GetFreeStorageSpace()
	if err != nil {
		log.Printf("Error getting storage info, %s", err)
		return
	}

	out, err := createEvent("STORAGE", storages)
	if err != nil {
		log.Printf("Error creating event, %s", err)
		return
	}
	fmt.Fprintf(c.Conn, "%s\n", out)
	log.Printf("Sent storage info %s", out)
}

func (c *Connection) sendBackSystem() {
	log.Print("Sending system info")
	cpu, err := system.GetCPUUsage()
	if err != nil {
		log.Printf("Error getting cpu info, %s", err)
		return
	}

	totalRam, freeRam, err := system.GetRAMUsage()
	if err != nil {
		log.Printf("Error getting ram info, %s", err)
		return
	}
	outStruct := shared.SystemResult{
		CPUUseage:   cpu,
		TotalMemory: totalRam,
		FreeMemory:  freeRam,
	}

	out, err := createEvent("SYSTEM", outStruct)
	if err != nil {
		log.Printf("Error creating event, %s", err)
		return
	}
	fmt.Fprintf(c.Conn, "%s\n", out)
	log.Printf("Sent system info %s", out)

}
