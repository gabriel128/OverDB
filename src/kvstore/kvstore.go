package kvstore

import "overdb/src/raft"
import "log"
import "errors"
// import "strconv"
import "encoding/json"

type Element struct {
	val string
	txn int
}

type KvStore struct {
	raft *raft.Raft
	raftCommCh chan(raft.ApplyMsg)
	store map[string][]Element
}

type PutArgs struct {
	Key string
	Val string
	Txn int
}

type PutReply struct {
	IsLeader bool
	Err string
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	IsLeader bool
	Val string
	Txn int
	Err string
}

type LogCommand  struct {
	Command string
	Args PutArgs
}

func Create(raft *raft.Raft, commCh chan raft.ApplyMsg) KvStore {
	kv := KvStore{}
	kv.raft = raft
	kv.raftCommCh = commCh

	// TODO: Recreate from logs

	return kv
}

func (kv *KvStore) Get(args *GetArgs, reply *GetReply) error {
	vals, ok := kv.store[args.Key]

	if ok {
		val := vals[len(vals)-1]
		reply.Txn = val.txn
		reply.Val = val.val
	}
	return nil
}

func (kv *KvStore) Put(args *PutArgs, reply *PutReply) error {
	log.Println("Put")

	// index, _, isLeader := kv.raft.SendCommand("Put,a" + args.Val + "," + strconv.Itoa(args.txn))
	command := LogCommand{"Put", *args}
	json_command, err := json.Marshal(command)

	if err != nil {
		reply.Err = "Parse Error"
		return errors.New("parse_error")
	}

	index, _, isLeader := kv.raft.SendCommand(json_command)

	if !isLeader {
		reply.IsLeader = false
		reply.Err = "Not a leader"

		return errors.New("not_a_leader")
	}

	msg := <-kv.raftCommCh

	if msg.CommandIndex != index {
		log.Printf("[Error] applying command %+v", msg)

		reply.IsLeader = true
		reply.Err = "Index out of sync"

		return errors.New("index_out_of_sync")

	} else {
		store_put(kv, *args)
		log.Printf("[Success] Message applied %+v", msg)
	}

	return nil
}

func store_put(kv *KvStore, args PutArgs) {
	if vals, ok := kv.store[args.Key]; ok {
		kv.store[args.Key] = append(vals, Element{val: args.Val, txn: args.Txn})
	} else {
		kv.store[args.Key] = []Element{Element{val: args.Val, txn: args.Txn}}
	}
}
