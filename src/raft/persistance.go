package raft

import "os"
import "io"
import "encoding/json"
import "bytes"
import "strconv"
import "log"

func marshal(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func unmarshal(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func saveFile(path string, obj interface{}) error {
	file, err := os.Create(path)
	defer file.Close()

	if err != nil {
		return err
	}

	reader, err := marshal(obj)

	if err != nil {
		return err
	}
	_, err = io.Copy(file, reader)
	return err
}

func loadFile(path string, obj interface{}) error {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return err
	}
	return unmarshal(file, obj)
}

func (rf *Raft) persist() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	path := "./" + strconv.Itoa(rf.me) + ".tmp"

	dao := RaftDao{CurrentTerm: rf.currentTerm, VotedFor: rf.votedFor, Log: rf.log}

	if err := saveFile(path, dao); err != nil {
		log.Fatalln(err)
	}
}

func (rf *Raft) readPersisted() {
	var dao RaftDao

	path := "./" + strconv.Itoa(rf.me) + ".tmp"

	err := loadFile(path, &dao)

	if err != nil && !os.IsNotExist(err) {
		log.Fatal("%+v", err)
	} else if err == nil {
		rf.mu.Lock()
		defer rf.mu.Unlock()

		rf.votedFor = dao.VotedFor
		rf.log = dao.Log
		rf.currentTerm = dao.CurrentTerm
	}
}
