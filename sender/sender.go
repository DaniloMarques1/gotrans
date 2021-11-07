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
	fmt.Println("You need to provide the receiver ip address (ip:port)")
	fmt.Print("> ")
	var address string
	if scanner.Scan() {
		address = scanner.Text()
	}

	fmt.Println("Provide the file path you want to transfer")
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

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	content := fmt.Sprintf("%s;%d\n%s",
		fileInfo.Name(), fileInfo.Size(), string(bytes))

	conn, err := net.Dial("tcp", sender.receiverAddr)
	if err != nil {
		log.Fatal(err) // TODO
	}
	defer conn.Close()

	// writing the header with the file name and size
	if _, err = conn.Write([]byte(content)); err != nil {
		log.Fatal(err) // TODO
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err) // TODO
	}
	log.Printf("Response = %v\n", response)
}
