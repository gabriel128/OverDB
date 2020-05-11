package main

import (
	"strings"
	"fmt"
	"os"
	"kgv/src/raft"
	"strconv"
	"log"
	// "time"
	"kgv/src/comms"
)

func main() {
	fmt.Println("KGV Starting")

	port1, _ := strconv.Atoi(os.Args[2])
	port2, _ := strconv.Atoi(os.Args[3])
	port3, _ := strconv.Atoi(os.Args[4])

	cmd_response := make(chan raft.ApplyMsg)
	var rf *raft.Raft

	go func() {
		rf = comms.StartHttpRPCServer(port1, port2, port3, cmd_response)
	}()

	// go comms.StartTcpRaftServer(port1, port2, port3, true)
	// go comms.StartTcpRaftServer(port2, port3, port1)
	// go comms.StartTcpRaftServer(port3, port1, port2)
	// go func() {

	var text string

	for {
		if rf == nil {
			continue
		}

		fmt.Println("Enter a command, possible commands: ")
		fmt.Println("1. getlog [Show logs in current server]")
		fmt.Println("2. anything [Anything that is not getlogs will get inserted in the logs]")
		fmt.Print("> ")

		fmt.Scanln(&text)
		trimmedInput := strings.TrimSpace(text)
		if "getlog" == strings.ToLower(trimmedInput) {
			log.Printf("[Logs] %+v\n", rf.GetLogs())
		} else {
			_, _, isLeader := rf.Start(trimmedInput)

			msg := <-cmd_response
			log.Println(msg)
			if !isLeader {
				log.Printf("[Logs] Sorry Not a leader\n")
			}
		}
	}
}
