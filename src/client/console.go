package client

import (
	"strings"
	"fmt"
	"os"
	// "strconv"
	"log"
	// "time"
	"bufio"
)

func (client *Client) Console() {
	reader := bufio.NewReader(os.Stdin)

	for {
		var text string

		fmt.Printf("\n\n")
		fmt.Println("Enter a command, possible commands: ")
		fmt.Println("1. get,key")
		fmt.Println("2. put,key,val")
		fmt.Print("> ")

		text, _ = reader.ReadString('\n')

		trimmedInput := strings.TrimSpace(strings.ToLower(text))

		if trimmedInput == "" {
			continue
		}

		split_input := strings.Split(trimmedInput, ",")

		if len(split_input) < 2 || len(split_input) > 3 {
			log.Printf("Invalid Input")
			continue
		}

		command := split_input[0]
		key := split_input[1]

		if "get" == command {
			val, txn := client.Get(key)

			log.Printf("\nVal: %s, Txn: %d", val, txn)
		} else if "put" == command {
			val := split_input[2]
			log.Println("About to put ", val)

			client.Put(key, val)


		} else {
			log.Printf("Not Valid command sorry")
		}
	}
}
