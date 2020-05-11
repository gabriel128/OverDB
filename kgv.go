package main

import (
	"strings"
	"fmt"
	"os"
	"kgv/src/raft"
	"strconv"
	"log"
	// "time"
	// "bufio"
	// "runtime"
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
	//	for {
	//		if rf == nil {
	//			continue
	//		}

	//		log.Printf("[%d][Logs] %+v", rf.Id(), rf.GetLogs())

	//		time.Sleep(3*time.Second)
	//	}
	// }()
	// input := bufio.NewReader(os.Stdin)
	var text string
	for {
		if rf == nil {
			continue
		}

		fmt.Print("Enter command [GetLog] [Anything]: ")

		// text, _ := reader.ReadString('\n')

		fmt.Scanln(&text)
		if "getlog" == strings.ToLower(text) {
			log.Printf("[Logs] %+v\n", rf.GetLogs())
		} else {
			rf.Start(text)
		}


		// _,_,isLeader := rf.Start("holaPepe")


		// response := <-cmd_response

		// log.Printf("[%d] Reponse is %+v, isLeader is %t", rf.Id(), response, isLeader)

		// time.Sleep(5*time.Second)
	}
}
