package main

import (
	"log"
	"os"
	"overdb/src/servers"
	"overdb/src/transaction_manager"
	"overdb/src/raft"
	"strconv"
)

func main() {
	log.Println("OverDB Starting ... ")

	serverType := os.Args[2]
	serverNumber, _ := strconv.Atoi(os.Args[3])
	commCh := make(chan raft.ApplyMsg)
	config := servers.ServersConfig

	if serverType == "tm" {
		currentServer := config.TransactionManager[serverNumber]
		peers := otherPeers(currentServer, config.TransactionManager)

		go func() {
			log.Println("TM server ", currentServer,  "Starting ... ")
			tm := transaction_manager.TransactionManager{}
			servers.StartHttpRPCServer([]int{currentServer, peers[0], peers[1]}, tm, commCh)
		}()
	}

	for {

	}
}

func otherPeers(currentServer int, allServers []int) []int{
	peers := []int{}
	for _, peer := range(allServers) {
		if peer != currentServer {
			peers = append(peers, peer)
		}
	}
	return peers
}
