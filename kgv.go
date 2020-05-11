package main

import (
	"strings"
	"fmt"
	"os"
	"kgv/src/raft"
	"strconv"
	"log"
	// "time"
	"kgv/src/servers"
)

func main() {
	fmt.Println("KGV Starting")

	port1, _ := strconv.Atoi(os.Args[2])
	port2, _ := strconv.Atoi(os.Args[3])
	port3, _ := strconv.Atoi(os.Args[4])

	cmd_response := make(chan raft.ApplyMsg)
	var rf *raft.Raft

	go func() {
		rf = servers.StartHttpRPCServer(port1, port2, port3, cmd_response)
	}()

	// go servers.StartTcpRaftServer(port1, port2, port3, true)
	// go servers.StartTcpRaftServer(port2, port3, port1)
	// go servers.StartTcpRaftServer(port3, port1, port2)

	var text string

	for {
		if rf == nil {
			continue
		}

		fmt.Println("Enter a command, possible commands: ")
		fmt.Println("1. getlogs [Show logs in current server]")
		fmt.Println("2. anything [Anything that is not getlogs will get inserted in the logs]")
		fmt.Print("> ")

		fmt.Scanln(&text)
		trimmedInput := strings.TrimSpace(text)

		if trimmedInput == "" {
			continue
		}

		if "getlogs" == strings.ToLower(trimmedInput) {
			log.Printf("[Logs] %+v\n", rf.GetLogs())
		} else {
			index, term, isLeader := rf.SendCommand(trimmedInput)

			log.Printf("[Sending command] on index %d, Term %d \n\n", index, term)

			if isLeader {
				msg := <-cmd_response

				if msg.CommandIndex != index {
					log.Printf("[Error] applying command %+v", msg)
				} else {
					log.Printf("[Success] Message applied %+v", msg)
					log.Printf("[Success] Log are now: %+v", rf.GetLogs())
				}
			} else {
				log.Printf("Can't apply commant to not leader\n")
			}
		}
	}
}
