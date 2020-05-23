package kvstore

// import "overdb/src/rm_server"
import "log"

type KvStore struct {

}

func Create() KvStore {
	kv := KvStore{}
	return kv
}

func (kv *KvStore) Put(args *int, reply *int) error {
	log.Println("Put")

	return nil
}
