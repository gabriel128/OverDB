## OverDB

This is a distributed, fault-tolerant, sharded key-value store with distributed transactions. This is built on top of Raft so most of the fancy words are atributed to that algorithm.

### Is it production ready?

This just for me to write distributed system stuff but if you want to use it in production, like a friend of mine would say... Sure

### How to run it

  WIP

### Why in go?

The obvious dad joke would be "I wanted to give it a go" but the reality is that I like playing with
concurrency and I like CSP stuff. I might re-write it in Haskell some day...

### Raft implementation

This project implements a version of Raft based on this [paper](https://raft.github.io/raft.pdf). I did some slight modifications to it though

- On conflict, the leader overrides the follower log in 1 step instead of do it incrementally. That removes the need of keeping track of PrevLogIndex and PrevLogTerm
- The first element in the log is always a spanshot/checkpoint

### Distributed Transactions

It implements a mix of determenistic transactions and 2 phase commit. I call it *1.5 Phase commit*

### Docs

The way the transactions work and the topology of the system is explained [HERE](https://www.notion.so/Distributed-Transactions-on-OverDB-8b8571e4e2ba4cc385fa34ee136f4dd2)

### Potential (overengineered of course) TODO

- [x] Raft Leader Election
- [x] Raft replicated log
- [x] Raft persistent log
- [x] Log Compaction/Snapshots
- [ ] Build a key value store on top of it
- [ ] Distributed transactions
- [ ] Sharding
- [ ] Hybrid clocks, Vector clocks or Lamport timestamps (not needed actually)
