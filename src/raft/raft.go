package raft

import "math/rand"
import "time"
import "log"
import "io/ioutil"
import "net/rpc"
import "kgv/src/dialers"

var _ = ioutil.Discard

// For testing Mostly
func (rf *Raft) GetState() (int, bool) {
	rf.mu.Lock()
	isLeader := rf.state == "leader"
	currentTerm := rf.currentTerm
	rf.mu.Unlock()

	return currentTerm, isLeader
}

func (rf *Raft) Id() int {
	return rf.me
}
//

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	err := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	if err != nil {
		log.Println("Error on AppendEntries", err)

		time.Sleep(1 * time.Second)
		client, err1 := dialers.DialHttp(server)

		if err1 == nil {
			rf.peers[server] = client
			log.Println("Reconnecting", server)
		}
	}

	return err == nil
}

func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	err := rf.peers[server].Call("Raft.RequestVote", args, reply)
	if err != nil {
		log.Println("Error on sendRequestVote", err)
	}
	return err == nil
}

func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	becomeFollowerIfBiggerTerm(rf, args.Term)

	if args.Term < rf.currentTerm {
		reply.VoteGranted = false
	} else if rf.votedFor == -1 || rf.votedFor == args.CandidateId {
		rf.votedFor = args.CandidateId
		reply.VoteGranted = true
		rf.receivedHB = true
	} else {
		reply.VoteGranted = false
	}

	reply.Term = rf.currentTerm
	return nil
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	log.Printf("[%d] [%s] Received HB from [%d] for Term [%d]", rf.me, rf.state, args.LeaderId, args.Term)

	rf.receivedHB = true

	becomeFollowerIfBiggerTerm(rf, args.Term)

	if rf.state == "candidate"  {
		rf.state = "follower"
		rf.votedFor = -1
		log.Printf("[%d] [%s] became follower cause HB", rf.me, rf.state)
	}

	if args.Term < rf.currentTerm {
		reply.Success = false
	} else {
		reply.Success = true
	}

	reply.Term = rf.currentTerm
	return nil
}

func heartBeat(rf *Raft) {
	for {
		rf.mu.Lock()
		state := rf.state
		rf.mu.Unlock()

		if state == "leader" {
			for i, _ := range rf.peers {
				if i != rf.me {
					go func(i int, term int) {
						args := &AppendEntriesArgs{Term: term, LeaderId: rf.me}
						reply := AppendEntriesReply{}

						ok := rf.sendAppendEntries(i, args, &reply)

						if !ok {
							return
						}

						rf.mu.Lock()
						becomeFollowerIfBiggerTerm(rf, reply.Term)
						rf.mu.Unlock()
					}(i, rf.currentTerm)
				}
			}
		}
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}

func becomeFollowerIfBiggerTerm(rf *Raft, term int) {
	if term > rf.currentTerm {
		rf.currentTerm = term
		rf.state = "follower"
		rf.votedFor = -1
		log.Printf("[%d] [%s] became follower with currentTerm [%d]", rf.me, rf.state, rf.currentTerm)
	}
}
func sleepRandomElectionTime(rf *Raft) {
	for {
		rand.Seed(time.Now().UnixNano())
		min := 3000
		max := 6000
		rtime := rand.Intn(max - min + 1) + min
		log.Printf("[%d] [%s] New Election time [%d]", rf.me, rf.state, rtime)
		time.Sleep(time.Duration(rtime) * time.Millisecond)

		rf.mu.Lock()
		if rf.receivedHB {
			rf.receivedHB = false
		} else {
			rf.mu.Unlock()
			break
		}
		rf.mu.Unlock()
	}
}
func startElection(rf *Raft, term int) {
	rf.mu.Lock()
	rf.votedFor = rf.me
	rf.currentTerm++
	rf.state = "candidate"
	rf.receivedHB = true

	log.Printf("[%d] [%s] became candidate of term %d", rf.me, rf.state, rf.currentTerm)
	rf.mu.Unlock()

	voted := make(chan bool)
	for i, _ := range rf.peers {
		if i != rf.me {
			go func(i int, currentTerm int) {

				rf.mu.Lock()
				currentTerm2 := rf.currentTerm

				args := RequestVoteArgs{Term: currentTerm2, CandidateId: rf.me}
				reply := RequestVoteReply{}

				rf.mu.Unlock()

				ok := rf.sendRequestVote(i, &args, &reply)

				rf.mu.Lock()
				becomeFollowerIfBiggerTerm(rf, reply.Term)
				rf.mu.Unlock()

				voted <- ok && reply.VoteGranted
			}(i, rf.currentTerm)
		}
	}

	votes := 1
	for i, _ := range rf.peers {
		if i != rf.me && <-voted {
			votes++

			if votes > len(rf.peers)/2 && rf.state == "candidate" && rf.currentTerm == (term + 1) {
				rf.mu.Lock()
				rf.state = "leader"
				log.Printf("[%d] [%s] became Leader of term %d", rf.me, rf.state, rf.currentTerm)
				rf.mu.Unlock()
				return
			}
		}
	}
}

func electionTimer(rf *Raft) {
	for {
		sleepRandomElectionTime(rf)

		rf.mu.Lock()
		state := rf.state
		rf.mu.Unlock()

		if state == "leader" {
			continue
		}

		startElection(rf, rf.currentTerm)
	}
}

func (rf *Raft) Start(peers map[int]*rpc.Client) {
	rf.peers = peers

	go electionTimer(rf)
	go heartBeat(rf)
}

func CreateRaftServer(me int, logs bool) *Raft {
	if !logs {
		log.SetOutput(ioutil.Discard)
	}
	rf := &Raft{}
	rf.me = me

	rf.commitedIndex = 0
	rf.lastApplied = 0
	rf.votedFor = -1
	rf.state = "follower"
	rf.receivedHB = false

	return rf
}
