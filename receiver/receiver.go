package receiver

import (
	"errors"
	"log"
	"net"
)

const (
	NoAddrFound = "Something went wrong while trying to find the local address"
)

func Execute() {
	localAddr, err := getLocalAddr()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Local addr = %v\n", localAddr)
}

func getLocalAddr() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New(NoAddrFound)
}
