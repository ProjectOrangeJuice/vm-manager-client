package system

import (
	"log"
	"net"

	"github.com/ProjectOrangeJuice/vm-manager-server/shared"
)

func GetNetworkDetails() ([]shared.Network, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error getting interfaces: %s", err)
		return nil, err
	}

	var network []shared.Network
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Printf("Error reading network details: %s", err)
			return nil, err
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if !v.IP.IsLoopback() && v.IP.To4() != nil {
					network = append(network, shared.Network{
						IP:   v.IP.String(),
						MAC:  iface.HardwareAddr.String(),
						Name: iface.Name,
					})
				}
			}
		}
	}

	return network, nil
}
