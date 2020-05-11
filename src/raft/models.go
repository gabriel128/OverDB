package raft

import "sync"
import "net/rpc"
//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

type LogEntry struct {
	Term int
	Command interface{}
}

type RaftDao struct {
	CurrentTerm int
	VotedFor int
	Log []LogEntry
}

type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     map[int]*rpc.Client // RPC end points of all peers
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Persistent State
	currentTerm int
	votedFor int
	// log []interface{}
	log []LogEntry
	lastApplied int
	// Volatile state
	commitIndex int

	state string
	receivedHB bool
}

type RequestVoteArgs struct {
	Term int
	CandidateId int
	LastLogIndex int
	LastLogTerm int
}

type RequestVoteReply struct {
	Term int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term int
	LeaderId int
	PrevLogIndex int
	PrevLogTerm int
	OverrideLogs bool

	Entries []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term int
	Success bool
}
