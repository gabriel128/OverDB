package tm

import "overdb/src/config"
import "overdb/src/dialers"
import "net/rpc"
import "log"

type TManager struct {
	kvstores map[int][]*rpc.Client

}

func Create() TManager {
	tm := TManager{}
	tm.kvstores = make(map[int][]*rpc.Client)
	return tm
}

func (tm TManager) DialOthers() {
	for k, servers := range(config.Servers.Kvstores) {
		tm.kvstores[k] = make([]*rpc.Client,3)

		for i, port := range(servers)  {
			client, _ := dialers.DialHttp(port)
			tm.kvstores[k][i] = client
		}
	}
}

func (tm *TManager) Commit(server int, args *int, reply *int) bool {
	err := tm.kvstores[0][0].Call("KvStore.Put", args, reply)

	if err != nil {
		log.Println("Error on Commit", err)
	}
	return err == nil
}
