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

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     map[int]*rpc.Client // RPC end points of all peers
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Persistent State
	currentTerm int
	votedFor int
	log []string

	// Volatile state
	commitedIndex int
	lastApplied int

	// Volatile leader state
	nextIndex []int
	matchIndex []int

	state string
	receivedHB bool
}

type RequestVoteArgs struct {
	Term int
	CandidateId int
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

	Entries []string
	LeaderCommitIndex int
}

type AppendEntriesReply struct {
	Term int
	Success bool
}
