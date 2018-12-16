package main

import (
	"fmt"
	"net"
)

var (
	localIPArray []string
)

func init() {
	addrs, e := net.InterfaceAddrs()
	if e != nil {
		panic(fmt.Sprintf("get local ip failed, %v", e))
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIPArray = append(localIPArray, ipnet.IP.String())
			}
		}
	}

	fmt.Println(localIPArray)
}
