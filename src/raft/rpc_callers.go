package raft

import "kgv/src/dialers"
import log "github.com/sirupsen/logrus"

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	err := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	if err != nil {
		log.Println("Error on AppendEntries", err)

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

func (rf *Raft) sendSetSnapshot(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	err := rf.peers[server].Call("Raft.SetSnapshot", args, reply)
	if err != nil {
		log.Println("Error on sendRequestVote", err)
	}
	return err == nil
}
