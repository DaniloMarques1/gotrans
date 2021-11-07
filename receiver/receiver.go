package receiver

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	NoAddrFound   = "Something went wrong while trying to find the local address"
	InvalidHeader = "Invalid header"
)

func Execute() {
	localAddr, err := getLocalAddr()
	if err != nil {
		log.Fatal(err)
	}
	socket, err := net.Listen("tcp", localAddr+":5000")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := socket.Accept()
	if err != nil {
		log.Fatal(err)
	}

	// it is our intention (at least for now) to handle one connection at a time
	handleConn(conn)
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	path := getPathFromUser(conn.RemoteAddr().String())
	buffer := make([]byte, 2048)
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 2))
	n, err := conn.Read(buffer)
	for err == nil && n > 0 {
		n, err = conn.Read(buffer)
	}

	var header string
	for idx, char := range buffer {
		if char != '\n' {
			header += string(char)
		} else {
			buffer = buffer[idx+1:]
			break
		}
	}
	fileName, fileSize, _ := parseHeader(header)
	buffer = buffer[:fileSize]

	localFile, err := os.Create(fmt.Sprintf("%v/%v", path, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer localFile.Close()
	localFile.Write(buffer)

	conn.Write([]byte("OK\n")) // ignoring error
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

func getPathFromUser(senderAddr string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%v is trying to send you a file \n", senderAddr)
	fmt.Println("Type in the path in which you'd like to store the file")
	fmt.Printf("> ")
	var path string
	if scanner.Scan() {
		path = scanner.Text()
	}
	return path
}

func parseHeader(header string) (string, int, error) {
	header = strings.Replace(header, "\n", "", -1)
	splited := strings.Split(header, ";")
	if len(splited) != 2 {
		return "", 0, errors.New(InvalidHeader)
	}
	fileSize := splited[1]
	size, err := strconv.Atoi(fileSize)
	if err != nil {
		return "", 0, err
	}
	fileName := splited[0]
	return fileName, size, nil
}
