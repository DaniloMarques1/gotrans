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

const PORT = 5000

func Execute() {
	localAddr, err := getLocalAddr()
	if err != nil {
		log.Fatal(err)
	}
	socket, err := net.Listen("tcp", fmt.Sprintf("%v:%v", localAddr, PORT))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Waiting to receive file...")
	conn, err := socket.Accept()
	if err != nil {
		log.Fatal(err)
	}

	// it is our intention (at least for now) to handle one connection at a time
	// basically, a receiver will receive from only one sender at a time
	handleConn(conn)
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Receiving file from %v\n", conn.RemoteAddr())

	path := getPathFromUser(conn.RemoteAddr().String())
	buffer := make([]byte, 2048)
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 2))
	n, err := conn.Read(buffer)
	for err == nil && n > 0 {
		localBuffer := make([]byte, 2048)
		n, err = conn.Read(localBuffer)
		if n > 0 {
			buffer = append(buffer, localBuffer...)
		}
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
	fileName, fileSize, err := parseHeader(header)
	if err != nil {
		log.Fatal(err)
	}
	buffer = buffer[:fileSize]
	fmt.Printf("Storing the file %v in path %v\n", fileName, path)

	localFile, err := os.Create(fmt.Sprintf("%v/%v", path, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer localFile.Close()
	_, err = localFile.Write(buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File stored successfully")

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
	fmt.Print("Type in the path in which you'd like to store the file:")
	fmt.Print("> ")
	var path string
	if scanner.Scan() {
		path = scanner.Text()
	}
	return path
}

func parseHeader(header string) (string, int, error) {
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
