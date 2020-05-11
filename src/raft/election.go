package raft

import "math/rand"
import "time"
import log "github.com/sirupsen/logrus"

func sleepRandomElectionTime(rf *Raft) {for {
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
	log.Printf("[%d] [%s] just became candidate of term %d", rf.me, rf.state, rf.currentTerm)
	rf.mu.Unlock()

	voted := make(chan bool)
	for _, i := range otherPeers(rf) {
		go func(i int, currentTerm int) {
			rf.mu.Lock()
			currentTerm2 := rf.currentTerm

			lastEntry := lastLogEntry(rf)

			args := RequestVoteArgs{
				Term: currentTerm2,
				CandidateId: rf.me,
				LastLogTerm: lastEntry.Term,
				LastLogIndex: len(rf.log),
			}

			reply := RequestVoteReply{}

			rf.mu.Unlock()

			ok := rf.sendRequestVote(i, &args, &reply)

			rf.mu.Lock()
			becomeFollowerIfBiggerTerm(rf, reply.Term)
			rf.mu.Unlock()

			voted <- ok && reply.VoteGranted
		}(i, rf.currentTerm)
	}

	processVotes(rf, voted, term)
}

func processVotes(rf *Raft, voted chan bool, term int) {
	votes := 1
	for i, _ := range rf.peers {
		if i != rf.me && <-voted {
			votes++

			if votes > len(rf.peers)/2 && rf.state == "candidate" && rf.currentTerm == (term + 1) {
				transformToLeader(rf)
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

func transformToLeader(rf* Raft) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.state = "leader"
	log.Printf("[%d] [%s] is Leader of term %d", rf.me, rf.state, rf.currentTerm)
}
