package servers

import (
	"net/rpc"
	"net"
	"log"
	"strconv"
	"overdb/src/raft"
	"overdb/src/dialers"
)


func startTcpServer(port int, raft *raft.Raft) {
	handler := rpc.NewServer()
	handler.Register(raft)

	ln, err := net.Listen("tcp", "localhost:" + strconv.Itoa(port))

	if err != nil {
		log.Printf("Error on listen: %s\n", err)
		return
	}

	for {
		cxn, err := ln.Accept()

		if err != nil {
			log.Printf("Error on accepting connection\n", err)
			return
		}

		handler.ServeConn(cxn)
	}
}

func StartTcpRaftServer(port1 int, port2 int, port3 int, logs bool, applyCh chan raft.ApplyMsg) *raft.Raft {
	raftServer := raft.CreateRaftServer(port1, logs)

	go startTcpServer(port1, raftServer)

	client1, _ := dialers.DialTCP(port2)
	client2, _ := dialers.DialTCP(port3)

	peers := map[int]*rpc.Client{
		port2: client1,
		port3: client2,
	}

	raftServer.StartServer(peers, applyCh)

	return raftServer
}
