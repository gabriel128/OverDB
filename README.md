## OverDB

This is a distributed, fault-tolerant, sharded key-value store with distributed transactions. This is built on top of Raft so most of the fancy words are atributed to that algorithm.

### Is it production ready?

This just for me to write distributed system stuff but if you want to use it in production, like a friend of mine would say... Sure

### How to run it

Install all the dependencies, then:

(Linux)

just run ./run_servers.sh , it will lunch all the servers in different windows and run a client

(Mac)

install a VM with linux or modify the scripts so it runs the equivalent

### Why in go?

The obvious dad joke would be "I wanted to give it a go" but the reality is that I like playing with
concurrency and I like CSP stuff. I might re-write it in Haskell some day...

### Raft implementation

This project implements a version of Raft based on this [paper](https://raft.github.io/raft.pdf). I did some slight modifications to it though

- On conflict, the leader overrides the follower log in 1 step instead of do it incrementally. That removes the need of keeping track of PrevLogIndex and PrevLogTerm
- The first element in the log is always a spanshot/checkpoint


### Distributed Transactions

It implements something similar to determenistic transactions with only 1 phase commit, see the Docs mentioned below. Even if it's only 1 commit request, I call it *1.5 Phase commit*, again see the Docs for more info.

### Docs

The way the transactions work and the topology of the system is explained [HERE](https://www.notion.so/Distributed-Transactions-on-OverDB-8b8571e4e2ba4cc385fa34ee136f4dd2)

### Sharding

For now it has "consitently stupid hashing" by a hash function and modulo on the quantity of servers. So migrations are manual. I might change it later to just "consistent hashing"

### Potential (overengineered of course) TODO

- [x] Raft Leader Election
- [x] Raft replicated log
- [x] Raft persistent log
- [x] Log Compaction/Snapshots
- [X] Build a key value store on top of it
- [X] Sharding
- [ ] Distributed transactions
- [ ] Hybrid clocks, Vector clocks or Lamport timestamps (not needed but why not)
