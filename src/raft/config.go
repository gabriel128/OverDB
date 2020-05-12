package raft

type Config struct {
	minElectionTimeMs int
	maxElectionTimeMs int
	heartBeatRateMs int
}

var config Config = Config{3000, 7000, 1000}
