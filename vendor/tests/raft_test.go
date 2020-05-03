package test

import "testing"
import "comms"
import "log"
import "raft"
import "math/rand"
import "time"
// import "net"


func TestInitialElection(t *testing.T) {
	log.Println("[TEST] Testing Initial Election")
	peers := make(map[int]*raft.Raft)
	serversReady := make(chan bool)

	go func() {
		rf := comms.StartTcpRaftServer(8000, 8001, 8002, true)
		peers[8000] = rf

		serversReady<-true
	}()

	go func() {
		rf := comms.StartTcpRaftServer(8001, 8002, 8000, true)

		peers[8001] = rf

		serversReady<-true
	}()

	go func() {
		rf := comms.StartTcpRaftServer(8002, 8000, 8001, true)
		peers[8002] = rf

		serversReady<-true
	}()

	for i := 0; i < 3; i++ {
		<-serversReady
	}

	checkOneLeader(peers)

	time.Sleep(500 * time.Millisecond)
	term1 := checkTerms(peers)

	time.Sleep(2 * time.Second)

	term2 := checkTerms(peers)
	if term1 != term2 {
		log.Printf("warning: term changed even though there were no failures")
	}

	checkOneLeader(peers)
	log.Println("[TEST] Leader Election PASSED")
}

// func TestReElection(t *testing.T) {
//	log.Println("[TEST] Testing ReElection")
//	peers := make(map[int]*raft.Raft)
//	serversReady := make(chan bool)

//	go func() {
//		rf := comms.StartTcpRaftServer(8000, 8001, 8002, true)
//		peers[8000] = rf

//		serversReady<-true
//	}()

//	go func() {
//		rf := comms.StartTcpRaftServer(8001, 8002, 8000, true)

//		peers[8001] = rf

//		serversReady<-true
//	}()

//	go func() {
//		rf := comms.StartTcpRaftServer(8002, 8000, 8001, true)
//		peers[8002] = rf

//		serversReady<-true
//	}()

//	for i := 0; i < 3; i++ {
//		<-serversReady
//	}

//	leader1 := checkOneLeader(peers)
//	// if the leader disconnects, a new one should be elected.
//	log.Println("[TEST] Disconnecting leader", leader1.Id())

//	conn1 = true
//	// conn := *conn1

//	time.Sleep(10 * time.Second)



//	delete(peers, leader1.Id())

//	// leader2 := checkOneLeader(peers)

//	// if leader1.Id() == leader2.Id() {
//	//	log.Fatalln("No new leader")
//	// }

//	// // if the old leader rejoins, that shouldn't
//	// // disturb the new leader.
//	// log.Println("[TEST] Connecting old leader again", leader1.Id())
//	// leader1.Connect()
//	// peers[leader1.Id()] = leader1

//	// time.Sleep(10 * time.Second)
//	// leader3 := checkOneLeader(peers)

//	// if leader3.Id() != leader2.Id() {
//	//	log.Fatalln("Adding old leader back disturbed the new leader")
//	// }

//	// // if there's no quorum, no leader should
//	// // be elected.
//	// disconnect(leader2)
//	// disconnect((leader2 + 1) % 3)
//	// time.Sleep(2 * time.Second)
//	// checkNoLeader()

//	// // if a quorum arises, it should elect a leader.

//	// cfg.connect((leader2 + 1) % servers)
//	// cfg.checkOneLeader()

//	// // re-join of last node shouldn't prevent leader from existing.
//	// cfg.connect(leader2)
//	// cfg.checkOneLeader()

//	log.Println("[TEST] Leader ReElection PASSED")
// }

func checkTerms(peers map[int]*raft.Raft) int {
	term := -1
	for _, raft := range peers {
			xterm, _ := raft.GetState()
			if term == -1 {
				term = xterm
			} else if term != xterm {
				log.Fatalf("servers disagree on term")
			}
	}
	return term
}

func checkOneLeader(peers map[int]*raft.Raft) *raft.Raft {
	for iters := 0; iters < 10; iters++ {
		ms := 450 + (rand.Int63() % 100)
		time.Sleep(time.Duration(ms) * time.Millisecond)

		leaders := make(map[int][]*raft.Raft)
		for _, raft := range peers {
			if term, leader := raft.GetState(); leader {
				leaders[term] = append(leaders[term], raft)
			}
		}

		lastTermWithLeader := -1
		for term, leaders := range leaders {
			if len(leaders) > 1 {
				log.Fatalf("term %d has %d (>1) leaders", term, len(leaders))
			}
			if term > lastTermWithLeader {
				lastTermWithLeader = term
			}
		}

		if len(leaders) != 0 {
			return leaders[lastTermWithLeader][0]
		}
	}
	log.Fatalf("expected one leader, got none")
	return nil
}
