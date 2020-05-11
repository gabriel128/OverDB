package raft

import "time"

func heartBeat(rf *Raft) {
	for {
		rf.mu.Lock()
		state := rf.state
		rf.mu.Unlock()

		if state == "leader" {
			for _, i := range otherPeers(rf) {
				go func(i int, term int) {
					args := &AppendEntriesArgs{
						Term: term,
						LeaderId: rf.me,
						OverrideLogs: false,
						PrevLogIndex: rf.commitIndex,
						PrevLogTerm: rf.log[rf.commitIndex].Term,
						Entries: []LogEntry{},
						LeaderCommit: rf.commitIndex,
					}
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
		time.Sleep(time.Duration(config.heartBeatRateMs) * time.Millisecond)
	}
}
