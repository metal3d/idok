package utils

import (
	"errors"
	"log"
	"net"
)

// return local ip that matches kodi network
// ignoring loopback and other net interfaces
func GetLocalInterfaceIP() (string, error) {
	ips, _ := net.LookupIP(GlobalConfig.Target)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Error while checking you interfaces: %v", err)
	}
	for _, ip := range ips {
		for _, iface := range ifaces {
			if iface.Flags&net.FlagLoopback != 0 {
				continue
			}

			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				_, subnet, _ := net.ParseCIDR(addr.String())
				switch v := addr.(type) {
				case *net.IPNet:

					if subnet.Contains(ip) {
						return v.IP.String(), nil
					}
				}

			}
		}
	}
	return "", errors.New("Unable to get local ip")
}
