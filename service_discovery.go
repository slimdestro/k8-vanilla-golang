package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// Get all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Iterate over all network interfaces
	for _, i := range interfaces {
		// Get all addresses of the interface
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Iterate over all addresses
		for _, addr := range addrs {
			// Get the IP address
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Check if the IP address is a valid IPv4 address
			if ip == nil || strings.Contains(ip.String(), ":") {
				continue
			}

			// Print the IP address
			fmt.Println(ip.String())
		}
	}
}
      