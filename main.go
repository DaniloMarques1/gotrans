package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/danilomarques1/gotrans/receiver"
	"github.com/danilomarques1/gotrans/sender"
)

const (
	SENDER   = 1
	RECEIVER = 2
)

// errors
const (
	InvalidInput  = "Invalid input"
	InvalidChoice = "You can only select 1 for sender or 2 for receiver"
)

func main() {
	choice, err := menu()
	if err != nil {
		log.Fatal(err)
	}

	if choice == SENDER {
		sender.Execute()
	} else if choice == RECEIVER {
		receiver.Execute()
	} else {
		// TODO
		log.Fatal(errors.New(InvalidChoice))
	}
}

func menu() (int, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var choice int

	fmt.Println("Pick one")
	fmt.Println("1. Sender")
	fmt.Println("2. Receiver")
	fmt.Print("> ")

	if scanner.Scan() {
		input := scanner.Text()
		var err error
		choice, err = strconv.Atoi(input)
		if err != nil {
			return -1, errors.New(InvalidInput)
		}
	}

	return choice, nil
}
