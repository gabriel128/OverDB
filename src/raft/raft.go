package raft

import "sync/atomic"
import "io/ioutil"
import "net/rpc"
import log "github.com/sirupsen/logrus"
import "time"
import "os"

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func (rf *Raft) GetState() (int, bool) {
	rf.mu.Lock()
	isLeader := rf.state == "leader"
	currentTerm := rf.currentTerm
	log.Printf("[%d] [%s]", rf.me, rf.state)
	rf.mu.Unlock()

	return currentTerm, isLeader
}

func (rf *Raft) GetLogs() []LogEntry {
	return rf.log[1:]
}

func (rf *Raft) Id() int {
	return rf.me
}

func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}


//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}


//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

func (rf *Raft) Start(command interface{}) (int, int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.state == "leader" {
		go appendToLog(rf, command)
		return len(rf.log), rf.currentTerm, true
	} else {
		return -1, -1, false
	}
}

func sendAppendEntriesRPC(rf *Raft, i int, args AppendEntriesArgs) (bool , AppendEntriesReply) {
	reply := AppendEntriesReply{}

	ok := rf.sendAppendEntries(i, &args, &reply)

	rf.mu.Lock()
	becomeFollowerIfBiggerTerm(rf, reply.Term)
	rf.mu.Unlock()

	return ok, reply
}

func otherPeers(rf *Raft) []int {
	others := []int{}
	for i, _ := range rf.peers {
		if i != rf.me {
			others = append(others, i)
		}
	}
	return others
}

func appendToLog(rf *Raft, command interface{}) {
	rf.mu.Lock()

	newEntryTerm := rf.currentTerm
	rf.log = append(rf.log, LogEntry{Term: rf.currentTerm, Command: command})
	rf.mu.Unlock()

	success := make(chan bool)
	for _, i := range otherPeers(rf) {
		go func(i int, term int) {
			overrideLogs := false
			for {
				rf.mu.Lock()

				if rf.state != "leader" || rf.currentTerm != newEntryTerm {
					rf.mu.Unlock()
					break
				}

				prevLogIndex := len(rf.log) - 2
				var entries []LogEntry;

				if overrideLogs {
					entries = rf.log
				} else {
					entries = []LogEntry{lastLogEntry(rf)}
				}

				// log.Printf("PrevLogIndex for [%d] is %d and Entries sent is %+v, lastApplied %d", i, prevLogIndex, entries, lastApplied)

				args := AppendEntriesArgs{
					Term: term,
					LeaderId: rf.me,
					OverrideLogs: overrideLogs,
					PrevLogIndex: prevLogIndex,
					PrevLogTerm: rf.log[prevLogIndex].Term,
					Entries: entries,
					LeaderCommit: rf.commitIndex,
				}

				rf.mu.Unlock()

				ok, reply := sendAppendEntriesRPC(rf, i, args)

				if ok && reply.Success {
					success<-true
					break
				} else if ok {
					overrideLogs = true
				}
			}
		}(i, rf.currentTerm)
	}

	appendSuccesses := 1
	for range otherPeers(rf) {
		if <-success {
			appendSuccesses++

			rf.mu.Lock()
			if appendSuccesses > len(rf.peers)/2 && newEntryTerm == rf.currentTerm {
				rf.commitIndex++
				rf.mu.Unlock()
				return
			}
			rf.mu.Unlock()
		}
	}
}

func applyLastCommit(rf *Raft, applyCh chan ApplyMsg) {
	for {
		rf.mu.Lock()
		if rf.commitIndex > rf.lastApplied {
			// log.Printf("[%d] [%s] [Applied log] %+v, commitIndex: %d", rf.me, rf.state, rf.log, rf.commitIndex)
			rf.lastApplied++
			logEntry := rf.log[rf.lastApplied]

			go func(lastApplied int) {
				applyCh <- ApplyMsg{
					CommandValid: true,
					Command: logEntry.Command,
					CommandIndex: lastApplied}
			}(rf.lastApplied)

		}
		rf.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
}

func becomeFollowerIfBiggerTerm(rf *Raft, term int) {
	if term > rf.currentTerm {
		rf.currentTerm = term
		rf.state = "follower"
		rf.votedFor = -1
		log.Printf("[%d] [%s] just became follower with currentTerm [%d]", rf.me, rf.state, rf.currentTerm)
	}
}

func lastLogEntry(rf *Raft) LogEntry {
	if len(rf.log) == 0 {
		return LogEntry{Term: 0, Command: 0}
	} else {
		return rf.log[len(rf.log)-1]
	}
}

// func Make(peers []*labrpc.ClientEnd, me int, persister *Persister, applyCh chan ApplyMsg) *Raft {
//	// log.SetOutput(ioutil.Discard)

//	rf := &Raft{}
//	rf.peers = peers
//	rf.persister = persister
//	rf.me = me

//	rf.commitIndex = 0
//	rf.lastApplied = 0
//	rf.votedFor = -1
//	rf.state = "follower"
//	rf.receivedHB = false
//	rf.log = append(rf.log, LogEntry{Term: 0, Command: 0})

//	go electionTimer(rf)
//	go heartBeat(rf)
//	go applyLastCommit(rf, applyCh)

//	// initialize from state persisted before a crash
//	rf.readPersist(persister.ReadRaftState())

//	return rf
// }

func (rf *Raft) StartServer(peers map[int]*rpc.Client, applyCh chan ApplyMsg) {
	rf.peers = peers

	go electionTimer(rf)
	go heartBeat(rf)
	go applyLastCommit(rf, applyCh)
}

func CreateRaftServer(me int, logs bool) *Raft {
	if !logs {
		log.SetOutput(ioutil.Discard)
	}
	rf := &Raft{}
	rf.me = me

	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.votedFor = -1
	rf.state = "follower"
	rf.receivedHB = false
	rf.log = append(rf.log, LogEntry{Term: 0, Command: 0})


	return rf
}
