package main

import (
	"strings"
	"fmt"
	"os"
	"overdb/src/raft"
	// "strconv"
	"log"
	"time"
	// "overdb/src/servers"
	"bufio"
)

func raft_console() {
	fmt.Println("KGV Starting ... ")

	// port1, _ := strconv.Atoi(os.Args[2])
	// port2, _ := strconv.Atoi(os.Args[3])
	// port3, _ := strconv.Atoi(os.Args[4])

	cmd_response := make(chan raft.ApplyMsg)
	var rf *raft.Raft
	// var rf1 *raft.Raft
	// var rf2 *raft.Raft
	// var rf3 *raft.Raft


	// go func() {
	//	rf = servers.StartHttpRPCServer(port1, port2, port3, cmd_response)
	// }()

	// go func() {
	//	rf1 = servers.StartTcpRaftServer(port1, port2, port3, true, cmd_response)
	// }()

	// go func() {
	//	rf2 = servers.StartTcpRaftServer(port2, port3, port1, true, cmd_response)
	// }()

	// go func() {
	//	rf3 = servers.StartTcpRaftServer(port3, port1, port2, true, cmd_response)
	// }()

	time.Sleep(5 * time.Second)

	reader := bufio.NewReader(os.Stdin)

	for {
		var text string

		if rf == nil {
			continue
		}

		fmt.Printf("\n\n")
		fmt.Println("Enter a command, possible commands: ")
		fmt.Println("1. get [Show logs in current server]")
		fmt.Println("2. put X [puts X value in the log]")
		fmt.Println("3. snapshot [create snapshot]")
		fmt.Println("4. raft [gets the raft state]")
		fmt.Print("> ")

		// fmt.Scanf("%s", &str)
		// fmt.Scan(&text)
		text, _ = reader.ReadString('\n')

		trimmedInput := strings.TrimSpace(strings.ToLower(text))

		if trimmedInput == "" {
			continue
		}

		fmt.Println("")

		if "get" == trimmedInput {
			log.Printf("\n[Logs] %+v", rf.GetLogs())
		} else if "raft" == trimmedInput {
			log.Printf("\n[Raft State] %+v", rf)
		} else if "snapshot" == trimmedInput {
			isLeader := rf.TakeSnapshot("somedata")

			if !isLeader {
				log.Printf("\n Can't snapshot a non leader")
			}

		} else if strings.HasPrefix(trimmedInput, "put ") {
			command := strings.TrimSpace(strings.TrimLeft(text, "put"))

			log.Println("About to put", strings.TrimSpace(command))
			index, term, isLeader := rf.SendCommand(command)

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
				log.Printf("Can't apply command to not leader\n")
			}
		} else  {
			log.Printf("Not Valid command sorry")
		}
	}
}
