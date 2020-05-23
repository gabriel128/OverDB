package raft

import log "github.com/sirupsen/logrus"

func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// log.Printf("[%d] [%s] Received RequestVote from [%d] for Term [%d]", rf.me, rf.state, args.CandidateId, args.Term)

	becomeFollowerIfBiggerTerm(rf, args.Term)

	if args.Term < rf.currentTerm {
		reply.VoteGranted = false
	} else if (rf.votedFor == -1 || rf.votedFor == args.CandidateId) && logUpToDate(rf, *args)  {
		rf.votedFor = args.CandidateId
		rf.receivedHB = true
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}

	reply.Term = rf.currentTerm

	return nil
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if len(args.Entries) > 0{
		// log.Printf("[%d] [%s] Received entries are: %+v", rf.me, rf.state, args.Entries)
	} else {
		log.Printf("[%d] [%s] Received AppendEntries from [%d] for Term [%d]", rf.me, rf.state, args.LeaderId, args.Term)
	}

	rf.receivedHB = true

	becomeFollowerIfBiggerTerm(rf, args.Term)

	if rf.state == "candidate" {
		rf.state = "follower"
		rf.votedFor = -1
		log.Printf("[%d] [%s] just became follower cause HB or log from leader", rf.me, rf.state)
	}

	reply.Term = rf.currentTerm
	reply.Success = true

	if args.Term < rf.currentTerm || args.LeaderCommit < rf.commitIndex {
		reply.Success = false
		return nil
	}

	lastLogEntry := lastLogEntry(rf)
	lastIndex := len(rf.log) - 1

	if args.OverrideLogs {
		rf.log = args.Entries
	} else if lastIndex > args.PrevLogIndex && lastLogEntry.Term == args.PrevLogTerm {
		rf.log = rf.log[:args.PrevLogIndex+1]
		appendEntriesToLog(rf, args)
	} else if lastIndex == args.PrevLogIndex && lastLogEntry.Term == args.PrevLogTerm {
		appendEntriesToLog(rf, args)
	} else {
		reply.Success = false
	}

	if args.LeaderCommit > rf.commitIndex {
		rf.commitIndex = min(args.LeaderCommit, len(rf.log) - 1)
	}

	return nil
}

func (rf *Raft) SetSnapshot(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	rf.mu.Lock()

	rf.receivedHB = true

	becomeFollowerIfBiggerTerm(rf, args.Term)

	if rf.state == "candidate" {
		rf.state = "follower"
		rf.votedFor = -1
		log.Printf("[%d] [%s] just became follower cause HB or log from leader", rf.me, rf.state)
	}

	reply.Term = rf.currentTerm
	reply.Success = true

	if args.Term < rf.currentTerm || args.LeaderCommit < rf.commitIndex {
		reply.Success = false
		return nil
	} else {
		rf.log = args.Entries
		rf.commitIndex = 0
		rf.lastApplied = 0
	}
	rf.mu.Unlock()
	rf.persist()

	return nil
}

func appendEntriesToLog(rf *Raft, args *AppendEntriesArgs) {
	for _, entry := range args.Entries {
		rf.log = append(rf.log, entry)
	}
}

func min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func logUpToDate(rf *Raft, args RequestVoteArgs) bool {
	if args.LastLogTerm == lastLogEntry(rf).Term {
		return args.LastLogIndex >= len(rf.log)
	} else {
		return args.LastLogTerm > lastLogEntry(rf).Term
	}
}
