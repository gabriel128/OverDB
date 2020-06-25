package kvstore

import "overdb/src/raft"
import "log"
import "errors"
import "strconv"
import "strings"
import "encoding/json"
import "sync"
// import "time"

const max_log_entries int = 5

type Element struct {
	Val string
	Txn int
}

type KvStore struct {
	mu sync.Mutex
	raft *raft.Raft
	raftCommCh chan(raft.ApplyMsg)
	leaderResponse chan(raft.ApplyMsg)
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

func Create(rf *raft.Raft, commCh chan raft.ApplyMsg) KvStore {
	kv := KvStore{}
	kv.raft = rf
	kv.raftCommCh = commCh
	kv.store = make(map[string][]Element)
	kv.leaderResponse = make(chan raft.ApplyMsg)

	go loadLogsInMemory(&kv)
	recoverStateFromLogs(&kv, rf.GetLogs())
	return kv
}

func loadLogsInMemory(kv *KvStore) {
	for {
		if len(kv.raft.GetLogs()) > (max_log_entries + 1) {
			log.Println("Taking snapshot")
			handleSnapshotCreation(kv)
		}

		msg := <-kv.raftCommCh
		log.Println("Message is ", msg)


		_, isLeader := kv.raft.GetState()

		if msg.Command == "Snapshot" {
			recoverStateFromLogs(kv, kv.raft.GetLogs())
		} else {
			storePut(kv, msg.Command.(string))
			log.Printf("Putting in memory %+v \n", msg)
		}

		if isLeader {
			kv.leaderResponse<-msg
		}
	}
}

func (kv *KvStore) Get(args *GetArgs, reply *GetReply) error {
	log.Printf("[KvStore state] %+v", kv.store)

	vals, ok := kv.store[args.Key]

	log.Printf("[KvStore state] Vals %+v", vals)

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

	log.Println("Got Args ", args)

	command := "Put,"+args.Key+","+args.Val+","+strconv.Itoa(args.Txn)
	index, _, isLeader := kv.raft.SendCommand(command)

	if !isLeader {
		reply.IsLeader = false
		reply.Err = "Not a leader"

		return errors.New("not_a_leader")
	}

	msg := <-kv.leaderResponse

	log.Printf("[KvStore Leader got here]")

	if msg.CommandIndex != index {
		log.Printf("[Error] applying command %+v", msg)

		reply.IsLeader = true
		reply.Err = "Index out of sync"

		return errors.New("index_out_of_sync")
	} else {
		// storePut(kv, command)

		log.Printf("[Success] command applied %+v", command)
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
			storePut(kv, logEntry.Command.(string))
		}
	}
}

func storePut(kv *KvStore, command string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	log.Println("Command about to put is", command)
	split_input := strings.Split(command, ",")

	if split_input[0] == "Put" {
		txn, _ := strconv.Atoi(split_input[3])
		args := &PutArgs{Key: split_input[1], Val: split_input[2], Txn: txn}

		if vals, ok := kv.store[args.Key]; ok {
			kv.store[args.Key] = append(vals, Element{Val: args.Val, Txn: args.Txn})
		} else {
			kv.store[args.Key] = []Element{Element{Val: args.Val, Txn: args.Txn}}
		}
	}
}

func handleSnapshotCreation(kv *KvStore) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	json_state, err := json.Marshal(kv.store)
	kv.raft.TakeSnapshot(json_state)

	if err != nil {
		log.Println("Error on taking snapshot")
	}
}
