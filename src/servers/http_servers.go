package servers

import (
	"net/rpc"
	"log"
	"strconv"
	"overdb/src/raft"
	"overdb/src/rm_server"
	"overdb/src/dialers"
	"net/http"
)


func startHttpServer(port int, raft *raft.Raft, rm_server *rm_server.RMServer) {
	rpc.Register(raft)
	rpc.Register(rm_server)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":"+ strconv.Itoa(port), nil)

	if err != nil {
		log.Println(err.Error())
	}
}

func StartHttpRPCServer(ports []int, rm_server rm_server.RMServer, applyCh chan raft.ApplyMsg) (*raft.Raft, rm_server.RMServer) {

	raftServer := raft.CreateRaftServer(ports[0], true)
	server := rm_server.Create()

	go startHttpServer(ports[0], raftServer, &server)

	client1, _ := dialers.DialHttp(ports[1])
	client2, _ := dialers.DialHttp(ports[2])

	peers := map[int]*rpc.Client{
		ports[1]: client1,
		ports[2]: client2,
	}

	raftServer.StartServer(peers, applyCh)

	return raftServer, server
}
