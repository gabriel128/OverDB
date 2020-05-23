package kvstore

import "overdb/src/rm_server"

type KVStore struct {

}

func (kv KVStore) Create() rm_server.RMServer {
	return kv
}
