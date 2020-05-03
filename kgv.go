package main

import (
	"fmt"
	"os"
	"strconv"
	// "runtime"
	"comms"
)

func main() {
	fmt.Println("KGV Starting")

	port1, _ := strconv.Atoi(os.Args[2])
	port2, _ := strconv.Atoi(os.Args[3])
	port3, _ := strconv.Atoi(os.Args[4])

	go comms.StartHttpRPC(port1, port2, port3)
	// go comms.StartTcpRaftServer(port1, port2, port3, true)
	// go comms.StartTcpRaftServer(port2, port3, port1)
	// go comms.StartTcpRaftServer(port3, port1, port2)


	for {

	}
}
