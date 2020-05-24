package main

import (
	"log"
	"os"
	"time"
	"overdb/src/servers"
	"overdb/src/config"
	// "overdb/src/tm"
	"overdb/src/raft"
	"overdb/src/client"
	// "overdb/src/kvstore"
	"strconv"
)

func main() {
	log.Println("OverDB Starting ... ")

	serverType := os.Args[2]
	partitionNumber, _ := strconv.Atoi(os.Args[3])
	serverNumber, _ := strconv.Atoi(os.Args[4])

	commCh := make(chan raft.ApplyMsg)
	s_config := config.Servers

	if serverType == "tm" {

		currentServer := s_config.TransactionManager[serverNumber]
		peers := otherPeers(currentServer, s_config.TransactionManager)

		go func() {
			log.Println("TM server ", currentServer,  "Starting ... ")
			servers.StartHttpTMServer([]int{currentServer, peers[0], peers[1]}, commCh)
		}()

	} else if serverType == "kv" {

		currentServer := s_config.Kvstores[partitionNumber][serverNumber]
		peers := otherPeers(currentServer, s_config.Kvstores[partitionNumber])

		go func() {
			log.Println("KVstore partition server ", currentServer,  "Starting ... ")
			servers.StartHttpKvServer([]int{currentServer, peers[0], peers[1]}, commCh)
		}()

	} else if serverType == "client" {
		client := new(client.Client)
		client.Connect()

		client.Console()

	}

	for {
		time.Sleep(time.Minute)
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
