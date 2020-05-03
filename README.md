## KGV (ex KGB)

### The Gabe Key Value store (though the acronym is sorted differently)

This is a distributed, fault-tolerant and sharded key-value store. This is built on top of Raft so most of the fancy words are atributed to that algorithm.


### The future of this

It's probaby going to be more things with fancy names in the future.


### Is it production ready?

This just for me to write distributed system stuff but if you want to use it in production, like a friend of mine would say... Sure

### Assumptions

For now this assumes only 3 servers, it should work with more though but I need to modify main() to allow n quantity of raft servers.

### How to start it
For now just do the following in three different terminals

```
go run kgv.go -- 8000 8001 8002
go run kgv.go -- 8001 8000 8002
go run kgv.go -- 8002 8000 8001
```

the first parameter is the "local" server port and the other 2 are the other servers

### TODO

- [x] Raft Leader Election
- [ ] Raft Append new log Entries
- [ ] Raft persistent state
- [ ] Build a key value store on top of it
- [ ] Log Compaction (Maybe)
