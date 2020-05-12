## KGV (it's pronounced KGB)

### The overenGineered Key Value store (though the acronym is sorted differently)
This is a distributed, fault-tolerant and sharded key-value store. This is built on top of Raft so most of the fancy words are atributed to that algorithm.



This is a Key Value Store that aims to have everything even if it's not actually necessary (e.g add vector clocks to it)

### Is it production ready?

This just for me to write distributed system stuff but if you want to use it in production, like a friend of mine would say... Sure

### Assumptions

For now this assumes only 3 raft servers, it works with more (and less) but I need to modify main() to allow n quantity of raft servers.

### How to start it
For now just do the following in three different terminals

```
go run kgv.go -- 8000 8001 8002
go run kgv.go -- 8001 8000 8002
go run kgv.go -- 8002 8000 8001
```
or

```
./run_servers.sh
```
if you have the gnome-terminal in your system.

The first parameter is the "local" server port and the other 2 are the other servers

### Why in go?

The obvious dad joke would be "I wanted to give it a go" but the reality is that I like playing with
concurrency and I like CSP stuff. I might re-write it in Haskell some day...

### Potential (overengineered of course) TODO

- [x] Raft Leader Election
- [x] Raft replicated log
- [x] Raft persistent log
- [x] Log Compaction/Snapshots
- [ ] Build a key value store on top of it
- [ ] Distributed transactions
- [ ] Sharded servers
- [ ] Vector clocks (or Lamport timestamps)
- [ ] Operational and Denotational semantics
- [ ] Proof of correctness of something in it (with agda maybe?)
- [ ] Some other unecessary but cool thing
