/* 
 :: ServiceDiscovery
 :: Contributed by @slimdestro
*/

package main

import (
	"fmt"
	"net"
	"strings"
)

func main() { 
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}
 
	for _, i := range interfaces { 
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			return
		} 
		for _, addr := range addrs { 
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
 
			if ip == nil || strings.Contains(ip.String(), ":") {
				continue
			}
			fmt.Println(ip.String())
		}
	}
}