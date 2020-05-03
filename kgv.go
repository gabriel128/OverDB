package main

import (
	"sync"
	"fmt"
	"net/rpc"
	"net"
	"log"
	"time"
	"os"
	"strconv"
	"raft"
	// "runtime"
	"net/http"
)

func dialTCP(port int) (*rpc.Client, error) {
	var client *rpc.Client
	var err error

	for {
		client, err = rpc.Dial("tcp", "localhost:" + strconv.Itoa(port))

		if err != nil {
			log.Println("Dial for port", strconv.Itoa(port), "failed", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	return client, err
}

func dialHttp(port int) (*rpc.Client, error) {
	var client *rpc.Client
	var err error

	for {
		client, err = rpc.DialHTTP("tcp", "localhost:" + strconv.Itoa(port))

		if err != nil {
			log.Println("Dial for port", strconv.Itoa(port), "failed", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	return client, err
}


func startHttpServer(port int, raft *raft.Raft) {
	rpc.Register(raft)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":"+ strconv.Itoa(port), nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func startTCPServer(port int, wg *sync.WaitGroup) {
	raft := new(raft.Raft)

	handler := rpc.NewServer()
	handler.Register(raft)

	ln, err := net.Listen("tcp", "localhost:" + strconv.Itoa(port))

	if err != nil {
		fmt.Printf("Error on listen: %s\n", err)
		return
	}

	for {
		cxn, err := ln.Accept()

		if err != nil {
			log.Printf("Error on accepting connection\n", err)
			return
		}

		go handler.ServeConn(cxn)
	}
}

func main() {
	fmt.Println("KGV Starting")

	port1, _ := strconv.Atoi(os.Args[2])
	port2, _ := strconv.Atoi(os.Args[3])
	port3, _ := strconv.Atoi(os.Args[4])

	raftServer := raft.CreateRaftServer(port1)

	go startHttpServer(port1, raftServer)

	client1, _ := dialHttp(port2)
	client2, _ := dialHttp(port3)

	peers := map[int]*rpc.Client{
		port2: client1,
		port3: client2,
	}

	raftServer.Start(peers)

	for {

	}
}
