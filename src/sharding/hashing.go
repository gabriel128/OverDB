package sharding

import "hash/fnv"
import "overdb/src/config"
import "log"

// Consistently stupid hashing
func GetShard(key string) int {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(key))
	selectedServer := int(algorithm.Sum32()) % len(config.Servers.Kvstores)
	log.Println("Key goes to: ", selectedServer)
	return selectedServer
}
