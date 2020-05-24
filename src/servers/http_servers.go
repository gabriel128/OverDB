package servers

import (
	"net/rpc"
	"log"
	"strconv"
	"overdb/src/raft"
	"overdb/src/tm"
	"overdb/src/kvstore"
	"overdb/src/dialers"
	"net/http"
)


func StartHttpKvServer(ports []int, applyCh chan raft.ApplyMsg) (*raft.Raft, *kvstore.KvStore) {
	raftServer := raft.CreateRaftServer(ports[0], true)
	kv_server := kvstore.Create(raftServer, applyCh)

	go startKvRpcServer(ports[0], raftServer, &kv_server)

	client1, _ := dialers.DialHttp(ports[1])
	client2, _ := dialers.DialHttp(ports[2])

	peers := map[int]*rpc.Client{
		ports[1]: client1,
		ports[2]: client2,
	}

	raftServer.StartServer(peers, applyCh)

	return raftServer, &kv_server
}

func startKvRpcServer(port int, raft *raft.Raft, kv *kvstore.KvStore) {
	rpc.Register(raft)
	rpc.Register(kv)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":"+ strconv.Itoa(port), nil)

	if err != nil {
		log.Println(err.Error())
	}
}

func StartHttpTMServer(ports []int, applyCh chan raft.ApplyMsg) (*raft.Raft, *tm.TManager) {

	raftServer := raft.CreateRaftServer(ports[0], true)
	tm_server := tm.Create()

	go startTmRpcServer(ports[0], raftServer, &tm_server)

	client1, _ := dialers.DialHttp(ports[1])
	client2, _ := dialers.DialHttp(ports[2])

	peers := map[int]*rpc.Client{
		ports[1]: client1,
		ports[2]: client2,
	}

	raftServer.StartServer(peers, applyCh)
	tm_server.DialOthers()


	return raftServer, &tm_server
}

func startTmRpcServer(port int, raft *raft.Raft, tm *tm.TManager) {
	rpc.Register(raft)
	rpc.Register(tm)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":"+ strconv.Itoa(port), nil)

	if err != nil {
		log.Println(err.Error())
	}
}
