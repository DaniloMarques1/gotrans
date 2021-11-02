package sender

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
	fmt.Println(file.Name())
}

// returns only the file name of the given path
func (sender *Sender) PathNameOnly() string {
	idx := strings.LastIndex(sender.path, "/")
	if idx == -1 {
		return sender.path
	}
	return sender.path[idx+1:]
}
