package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/SUMUKHA-PK/Raft-Distributed-Consensus/servermanagement"
	"github.com/SUMUKHA-PK/Raft-Distributed-Consensus/types"
)

func main() {
	raftServers := make(map[string]types.RaftServer)
	configuration := initializeConfiguration("server.config.json")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go servermanagement.StartServers(configuration)
	go servermanagement.StartSignal(configuration, raftServers, 1, wg)
	wg.Wait()
}

func initializeConfiguration(filepath string) types.Configuration {
	file, err := os.Open(filepath)

	if err != nil {
		log.Panic("Error reading from %v: %v", filepath, err)
	}

	var configuration types.Configuration

	err = json.NewDecoder(file).Decode(&configuration)

	if err != nil {
		log.Panic("Error decoding %v: %v", filepath, err)
	}

	return configuration
}
