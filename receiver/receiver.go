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
)

const (
	NoAddrFound   = "Something went wrong while trying to find the local address"
	InvalidHeader = "Invalid header"
)

const PORT = 5000

type InfoFile struct {
	name string
	size int
	path string
	mode os.FileMode
}

func (infoFile *InfoFile) GetFullPath() string {
	return fmt.Sprintf("%v/%v", infoFile.path, infoFile.name)
}

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

	header, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	conn.Write([]byte("OK\n")) //

	infoFile, err := parseHeader(header)
	if err != nil {
		log.Fatal(err)
	}
	infoFile.path = path

	buffer := make([]byte, infoFile.size)
	if _, err = conn.Read(buffer); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Storing the file %v in path %v\n", infoFile.name, infoFile.path)

	localFile, err := os.OpenFile(infoFile.GetFullPath(),
		os.O_CREATE|os.O_WRONLY, infoFile.mode)
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
	fmt.Println("Type in the path in which you'd like to store the file:")
	fmt.Print("> ")
	var path string
	if scanner.Scan() {
		path = scanner.Text()
	}
	return path
}

func parseHeader(header string) (*InfoFile, error) {
	header = strings.Replace(header, "\n", "", -1)
	splited := strings.Split(header, ";")
	if len(splited) != 3 {
		return nil, errors.New(InvalidHeader)
	}
	fileSize := splited[1]
	size, err := strconv.Atoi(fileSize)
	if err != nil {
		return nil, err
	}
	fileMode := splited[2]
	mode, err := strconv.Atoi(fileMode)
	if err != nil {
		return nil, err
	}
	fileName := splited[0]

	infoFile := InfoFile{
		name: fileName,
		size: size,
		mode: os.FileMode(mode),
	}

	return &infoFile, nil
}
