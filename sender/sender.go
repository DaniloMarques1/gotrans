package sender

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Sender struct {
	// path for file that you want to transfer
	path         string
	receiverAddr string
}

func NewSender(path, addr string) *Sender {
	return &Sender{
		path:         path,
		receiverAddr: addr,
	}
}

func Execute() {
	path, address := getPathAndAddrFromUser()
	sender := NewSender(path, address)
	sender.Send()
}

func getPathAndAddrFromUser() (string, string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("You need to provide the receiver ip address (ip:port) make sure the receiver is waiting:")
	fmt.Print("> ")
	var address string
	if scanner.Scan() {
		address = scanner.Text()
	}

	fmt.Println("Provide the file path you want to transfer:")
	fmt.Print("> ")
	var path string
	if scanner.Scan() {
		path = scanner.Text()
	}

	return path, address
}

func (sender *Sender) Send() {
	file, err := os.Open(sender.path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	conn, err := net.Dial("tcp", sender.receiverAddr)
	if err != nil {
		log.Fatal(err) // TODO
	}
	defer conn.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileName := fileInfo.Name()
	fileSize := fileInfo.Size()
	fileMode := fileInfo.Mode()

	fmt.Printf("Sending file %v, with %v bytes and mode %d to computer %v\n",
		fileName, fileSize, fileMode, sender.receiverAddr)

	header := fmt.Sprintf("%s;%d;%d\n", fileName, fileSize, fileMode)
	if _, err := conn.Write([]byte(header)); err != nil {
		log.Fatal(err)
	}

	// read the header response
	if _, err = bufio.NewReader(conn).ReadString('\n'); err != nil {
		log.Fatal(err) // TODO
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err) // TODO
	}
	if _, err := conn.Write(bytes); err != nil {
		log.Fatal(err) // TODO
	}

	// read the body(file) response
	_, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err) // TODO
	}
	fmt.Println("File transfered successfully")
}
