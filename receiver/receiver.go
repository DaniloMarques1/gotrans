package receiver

import (
	"errors"
	"fmt"
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
	socket, err := net.Listen("tcp", localAddr+":5000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// it is our intention (at least for now) to hanlde onde connection at a time
		handleConn(conn)
	}
}

// TODO get the file name
// open a file based on file name
// copy the rest of the bytes in the conn to file
func handleConn(conn net.Conn) {
	fmt.Println("Handle conn...")
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
