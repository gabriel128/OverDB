package client

import "overdb/src/config"
import "overdb/src/dialers"
import "overdb/src/kvstore"
import "net/rpc"
import "log"

type Client struct {
	txn int
	// leader int
	// leaderCached bool
	kvstores map[int][]*rpc.Client
}

func (client Client) Connect() {
	for k, servers := range(config.Servers.Kvstores) {
		client.kvstores[k] = make([]*rpc.Client,3)

		for i, port := range(servers)  {
			rpc_client, _ := dialers.DialHttp(port)
			client.kvstores[k][i] = rpc_client
		}
	}
}

func (client *Client) Put(key string, val string) bool {
	args := kvstore.PutArgs{key, val, client.txn}
	reply := kvstore.PutReply{}

	store_for_key := key_from_hash(key)

	error := false

	for i, _ := range(client.kvstores[store_for_key]) {
		err := client.kvstores[store_for_key][i].Call("KvStore.Put", args, reply)

		if err != nil {
			error = true
			log.Println("Error on Put", err)
		}

		if reply.IsLeader {
			break
		}
	}

	return error
}

func (client *Client) Get(key string) (string, int) {
	args := kvstore.GetArgs{key}
	reply := kvstore.GetReply{}

	store_for_key := key_from_hash(key)

	for i, _ := range(client.kvstores[store_for_key]) {
		err := client.kvstores[store_for_key][i].Call("KvStore.Get", args, reply)

		if err != nil {
			log.Println("Error on Get", err)
		}

		if reply.IsLeader {
			break
		}
	}

	return reply.Val, reply.Txn
}

func key_from_hash(key string) int {
	return 0
}
