package kvstore

import "overdb/src/raft"
import "log"
import "errors"
import "strconv"
import "strings"
import "encoding/json"
import "sync"
// import "time"

const max_log_entries int = 2

type Element struct {
	Val string
	Txn int
}

type KvStore struct {
	mu sync.Mutex
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
	kv.store = make(map[string][]Element)

	recoverStateFromLogs(&kv, raft.GetLogs())
	return kv
}

func (kv *KvStore) Get(args *GetArgs, reply *GetReply) error {
	log.Printf("[KvStore state] %+v", kv.store)

	vals, ok := kv.store[args.Key]

	if ok {
		val := vals[len(vals)-1]
		reply.Txn = val.Txn
		reply.Val = val.Val
	} else {
		reply.Txn = -1
		reply.Val = "nil"
	}
	return nil
}

func (kv *KvStore) Put(args *PutArgs, reply *PutReply) error {
	log.Println("KvStore Put")

	index, _, isLeader := kv.raft.SendCommand("Put,"+args.Key+","+args.Val+","+strconv.Itoa(args.Txn))
	// command := LogCommand{"Put", *args}
	// json_command, err := json.Marshal(command)

	// if err != nil {
	//	reply.Err = "Parse Error"
	//	return errors.New("parse_error")
	// }

	// index, _, isLeader := kv.raft.SendCommand(json_command)


	if !isLeader {
		reply.IsLeader = false
		reply.Err = "Not a leader"

		go func(kv *KvStore) {
			<-kv.raftCommCh
			log.Printf("[KvStore got here]")
			storePut(kv, args)
		}(kv)

		return errors.New("not_a_leader")
	}

	msg := <-kv.raftCommCh


	if msg.CommandIndex != index {
		log.Printf("[Error] applying command %+v", msg)

		reply.IsLeader = true
		reply.Err = "Index out of sync"

		return errors.New("index_out_of_sync")
	} else {
		storePut(kv, args)

		if len(kv.raft.GetLogs()) > (max_log_entries + 1) {
			go handleSnapshotCreation(kv)
		}
		log.Printf("[Success] Message applied %+v", msg)
	}

	return nil
}

func recoverStateFromLogs(kv *KvStore, logs []raft.LogEntry) {
	for _,logEntry := range(logs) {
		if logEntry.IsSnapshot && len(logEntry.Data) > 0 {
			if err := json.Unmarshal(logEntry.Data, &kv.store); err != nil {
				log.Printf("Error Unmarshalling stuff")
			}
		} else if !logEntry.IsSnapshot {
			split_input := strings.Split(logEntry.Command.(string), ",")

			if split_input[0] == "Put" {
				txn, _ := strconv.Atoi(split_input[3])
				args := &PutArgs{Key: split_input[1], Val: split_input[2], Txn: txn}
				storePut(kv, args)
			}
		}
	}
}

func storePut(kv *KvStore, args *PutArgs) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if vals, ok := kv.store[args.Key]; ok {
		kv.store[args.Key] = append(vals, Element{Val: args.Val, Txn: args.Txn})
	} else {
		kv.store[args.Key] = []Element{Element{Val: args.Val, Txn: args.Txn}}
	}
}

func handleSnapshotCreation(kv *KvStore) {
	// time.Sleep(1 * time.Second)

	kv.mu.Lock()
	defer kv.mu.Unlock()

	json_state, err := json.Marshal(kv.store)
	kv.raft.TakeSnapshot(json_state)

	if err != nil {
		log.Println("Error on taking snapshot")
	}
}
