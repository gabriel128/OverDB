package raft

func (rf *Raft) TakeSnapshot(data []byte) bool {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.state == "leader" {
		go snapshotLog(rf, data)

		return true
	} else {
		return false
	}
}

func snapshotLog(rf *Raft, data []byte) {
	rf.mu.Lock()

	logEntry := LogEntry{Term: rf.currentTerm, IsSnapshot: true, Data: data}
	state := rf.state
	rf.log = []LogEntry{logEntry}

	rf.mu.Unlock()

	success := make(chan bool)
	if state == "leader" {
		for _, i := range otherPeers(rf) {
			go func(i int, term int) {
				args := &AppendEntriesArgs{
					Term: term,
					LeaderId: rf.me,
					OverrideLogs: false,
					Entries: rf.log,
					LeaderCommit: rf.commitIndex,
				}
				reply := AppendEntriesReply{}

				ok := rf.sendSetSnapshot(i, args, &reply)

				if !ok {
					return
				}

				rf.mu.Lock()
				becomeFollowerIfBiggerTerm(rf, reply.Term)
				rf.mu.Unlock()

				success<-true
			}(i, rf.currentTerm)
		}
	}


	snapshotSuccesses := 1

	for range otherPeers(rf) {
		if <-success {
			snapshotSuccesses++

			rf.mu.Lock()
			if snapshotSuccesses > len(rf.peers)/2 {
				rf.commitIndex = 0
				rf.lastApplied = 0
				rf.mu.Unlock()
				rf.persist()
				return
			}
			rf.mu.Unlock()
		}
	}
}
